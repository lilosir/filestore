package handler

import (
	"encoding/json"
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
			FilePath: "./tmp/" + fileHeader.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.FilePath)
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
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
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

	fileMeta := meta.GetFileMeta(fileHash)

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

	limitCount, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas, err := meta.GetLastMetas(limitCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// object ---> json
	data, err := json.Marshal(fileMetas)
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
	file, err := os.Open(fileMeta.FilePath)
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
	err := os.Remove(fileMeta.FilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete file index
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
