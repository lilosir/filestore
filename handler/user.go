package handler

import (
	dbUser "fileStore/db"
	"fileStore/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	passwordSalt = "+@#$%"
)

// SignUpHander handle user sign up
func SignUpHander(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			io.WriteString(w, "internal server error: could not find sign up page")
			return
		}
		w.Write(page)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Can not parse form, err: %v", err)
			return
		}

		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		if len(userName) < 3 || len(userPassword) < 5 {
			w.Write([]byte("Invalide parameter"))
			return
		}
		//encrypted password
		encPassword := util.Sha1([]byte(userPassword + passwordSalt))

		ok := dbUser.UserSignUp(userName, encPassword)
		if ok {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	default:
		fmt.Fprintf(w, "Request method not supported")
	}
}
