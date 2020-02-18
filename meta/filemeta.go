package meta

import (
	"fmt"
	"sort"

	mydb "fileStore/db"
)

// FileMeta denotes file source data struct
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	FileAddr string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta add to/update file meta, used for local storage
func UpdateFileMeta(filemeta FileMeta) {
	fileMetas[filemeta.FileSha1] = filemeta
}

// UpdateFileMetaDB add to/update database, prepare and execute
func UpdateFileMetaDB(filemeta FileMeta) bool {
	return mydb.OnFileUploadFinished(filemeta.FileSha1, filemeta.FileName, filemeta.FileSize, filemeta.FileAddr)
}

// GetFileMeta returns a file meta via file sha1
func GetFileMeta(filesha1 string) FileMeta {
	return fileMetas[filesha1]
}

// GetFileMetaDB returns a file meta via file sha1 from mysql
func GetFileMetaDB(filesha1 string) (FileMeta, error) {
	fileMeta := FileMeta{}
	tfile, err := mydb.GetFileMeta(filesha1)
	if err != nil || tfile == nil {
		return fileMeta, err
	}

	fileMeta.FileSha1 = tfile.FileHash
	fileMeta.FileName = tfile.FileName.String
	fileMeta.FileSize = tfile.FileSize.Int64
	fileMeta.FileAddr = tfile.FileAddr.String

	return fileMeta, nil
}

// GetLastMetas returns last amount of the global file metas
func GetLastMetas(count int) ([]FileMeta, error) {
	metaSlice := []FileMeta{}
	for _, v := range fileMetas {
		metaSlice = append(metaSlice, v)
	}

	if count > len(metaSlice) {
		return nil, fmt.Errorf("exceed limitation: %d", len(metaSlice))
	}

	sort.Sort(ByUploadTime(metaSlice))
	return metaSlice[:count], nil
}

// RemoveFileMeta delete a file
func RemoveFileMeta(filesha1 string) {
	//delete maybe thread is sychornized, probably need to add mutex
	delete(fileMetas, filesha1)
}
