package meta

import (
	"fmt"
	"sort"
)

// FileMeta denotes file source data struct
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	FilePath string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta add/update file meta
func UpdateFileMeta(filemeta FileMeta) {
	fileMetas[filemeta.FileSha1] = filemeta
}

// GetFileMeta returns a file meta via file sha1
func GetFileMeta(filesha1 string) FileMeta {
	return fileMetas[filesha1]
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
