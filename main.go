package main

import (
	"fileStore/handler"
	"fileStore/handler/middleware"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", middleware.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/success", middleware.HTTPInterceptor(handler.UploadSuccessHandler))
	http.HandleFunc("/file/meta", middleware.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/query", middleware.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/download", middleware.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", middleware.HTTPInterceptor(handler.UpdateHandler))
	http.HandleFunc("/file/delete", middleware.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/fastupload", middleware.HTTPInterceptor(handler.TryFastUploadHandler))

	http.HandleFunc("/user/signup", handler.SignUpHander)
	http.HandleFunc("/user/signin", handler.SignInHander)
	http.HandleFunc("/user/info", middleware.HTTPInterceptor(handler.UserInfoHandler))

	http.HandleFunc("/file/mpupload/init", middleware.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", middleware.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", middleware.HTTPInterceptor(handler.CompleteUploadHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err: %s\n", err.Error())
		panic(err)
	}

	println("Running code after ListenAndServe (only happens when server shuts down)")
}
