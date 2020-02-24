package handler

import (
	rPool "fileStore/cache/redis"
	dbLayer "fileStore/db"
	"fileStore/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// MultiPartUploadInfo denote info
type MultiPartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

//InitialMultipartUploadHandler initializes multiple upload
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse form
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 2. get redis pool connection
	rConn := rPool.CreateRedisPool().Get()
	defer rConn.Close()

	// 3. generate multi-part upload intial info
	chunkSize := 10 * 1024 * 1024 //10 MB
	uploadInfo := MultiPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now()),
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / float64(chunkSize))),
	}

	//4. wirte the initial info to redis cache
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunkcount", uploadInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filehash", uploadInfo.FileHash)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filesize", uploadInfo.FileSize)

	//5. respond client
	w.Write(util.NewRespMsg(0, "OK", uploadInfo).JSONBytes())
}

// UploadPartHandler upload multi parts of a file
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse form
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 2. get redis pool connection
	rConn := rPool.CreateRedisPool().Get()
	defer rConn.Close()

	// 3. create a file which is used for multi store info
	fPath := "./tmp/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	file, err := os.Create(fPath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024) // 1 MB
	for {
		n, err := r.Body.Read(buf)
		file.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. uodate redis cache status
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. respond client
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler inform user upload finished
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse form
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uploadID := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// 2. get redis pool connection
	rConn := rPool.CreateRedisPool().Get()
	defer rConn.Close()

	// 3. Via uploadID check if redis has finished uploading all parts
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed, redis query error", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0

	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	fmt.Printf("total: %d, chunkcount: %d\n", totalCount, chunkCount)
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "Invalid request", nil).JSONBytes())
		return
	}

	// 4. TODO: merge multi parts

	// 5. upload tbl_file and tbl_user_file
	fsize, _ := strconv.Atoi(filesize)
	dbLayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dbLayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// 6. respond client
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
