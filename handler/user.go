package handler

import (
	dbUser "fileStore/db"
	"fileStore/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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

// SignInHander handle user sign in
func SignInHander(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + passwordSalt))

	//step 1: check username and password
	pwdChecked := dbUser.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("Failed"))
		return
	}

	//step 2: authentication (token or session/cookites)
	token := GenToken(username)
	ok := dbUser.UpdateToken(username, token)
	if !ok {
		w.Write([]byte("Failed"))
		return
	}

	//setp 3: redirect to main page after login
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

// GenToken get a user token which is used for authentication
func GenToken(username string) string {
	// 40characters md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
