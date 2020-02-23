package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}

	return nil
}

func main() {
	username := "admin"
	token := "58878710bef1ea9213046ae5c45e835f5e520f30"
	filehash := "376b96274e9f1ff2c4ca63b336d45ace2244191c"

	// 1. 请求初始化分块上传接口
	// resp, err := http.PostForm(
	// 	"http://localhost:8080/file/mpupload/init",
	// 	url.Values{
	// 		"username": {username},
	// 		"token":    {token},
	// 		"filehash": {filehash},
	// 		"filesize": {"677499765"},
	// 	})

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }

	// var dat map[string]interface{}
	// json.Unmarshal(body, &dat)

	// // 2. 得到uploadID以及服务端指定的分块大小chunkSize
	// uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	// chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	// ChunkCount := jsonit.Get(body, "data").Get("ChunkCount").ToInt()
	// fmt.Printf("chunk count: %d", ChunkCount)
	// fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// // 3. 请求分块上传接口
	// filename := "/Users/osir/Desktop/Dockerdmg.zip"
	// tURL := "http://localhost:8080/file/mpupload/uppart?" +
	// 	"username=admin&token=" + token + "&uploadid=" + uploadID
	// multipartUpload(filename, tURL, chunkSize)

	// 4. 请求分块完成接口
	resp, err := http.PostForm(
		"http://localhost:8080/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"677499765"},
			"filename": {"Dockerdmg.zip"},
			"uploadid": {"admin323032302d30322d32332030303a35343a35342e313436323239202d3035303020455354206d3d2b32342e303835303737363935"},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}
