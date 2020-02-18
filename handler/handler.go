package handler

import (
	"encoding/json"
	"fileStore/db"
	"fileStore/meta"
	"fileStore/util"
	"strconv"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler handle file uploading
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// return upload HTML view
		file, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			// panic("could not find file")
			io.WriteString(w, "internal server error: could not find file")
			return
		}
		io.WriteString(w, string(file))
	} else if r.Method == "POST" {
		// receive file I/O stream, and store to local
		// 1. get the file via post method form
		// 2. create a FileMeta for the current file
		// 3. create a new file
		// 4. update the FileMeta info: fileSize and fileSha1
		// 5. update the global fileMetas, which is used for CRUD

		// get the form file via post method
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get file, err: %s\n", err.Error())
			return
		}
		defer file.Close()

		// create filemeta which is used for get/add/update file
		fileMeta := meta.FileMeta{
			FileName: fileHeader.Filename,
			FileAddr: "./tmp/" + fileHeader.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.FileAddr)
		if err != nil {
			fmt.Printf("Failed to create file, err: %v\n", err)
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save file, err: %s\n", err.Error())
			return
		}
		_, err = newFile.Seek(0, 0)
		if err != nil {
			fmt.Printf("Failed to get file info, err: %s\n", err.Error())
			return
		}
		fileMeta.FileSha1 = util.FileSha1(newFile)

		// update the global fileMetas
		// meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		//TODO: update user file table
		r.ParseForm()
		username := r.Form.Get("username")
		ok := db.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if !ok {
			w.Write([]byte("Upload Failed."))
		}

		http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	}
}

// UploadSuccessHandler return success message
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "uploaded successfully!")
}

// GetFileMetaHandler get file meta via querying
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["filehash"][0]

	// fileMeta := meta.GetFileMeta(fileHash)
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// any interface ---> json
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileQueryHandler query latest file metas by limitation
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	limitCount, _ := strconv.Atoi(r.Form.Get("limit"))
	// fileMetas, err := meta.GetLastMetas(limitCount)

	userFiles, err := db.QueryUserFileMetas(username, limitCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// object ---> json
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler download file via querying filehash
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileMeta := meta.GetFileMeta(r.Form.Get("filehash"))
	file, err := os.Open(fileMeta.FileAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// if the file is small, use following ioutil.ReadAll
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// UpdateHandler just update file name
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta.FileName = newFileName
	meta.UpdateFileMeta(fileMeta)

	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler handle delete a file
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileSha1)

	//delete file from disk
	err := os.Remove(fileMeta.FileAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete file index
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}

// TryFastUploadHandler logic
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 1. parse form parameters
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2. check if there is same filehash in the database
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 3. return false if there is no record
	if fileMeta.FileSha1 == "" {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "fast upload failed, use regular upload please",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 4. otherwise upload file to database and return true
	ok := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if ok {
		resp := util.RespMsg{
			Code: 1,
			Msg:  "fast upload successfully",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: 1,
		Msg:  "upload failed, try it again.",
	}
	w.Write(resp.JSONBytes())
	return
}
