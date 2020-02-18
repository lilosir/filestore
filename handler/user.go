package handler

import (
	dbUser "fileStore/db"
	"fileStore/util"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	passwordSalt  = "+@#$%"
	validDuration = int64(60 * 60 * 24 * 365)
)

// SignUpHander handle user sign up
func SignUpHander(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// page, err := ioutil.ReadFile("./static/view/signup.html")
		// if err != nil {
		// 	io.WriteString(w, "internal server error: could not find sign up page")
		// 	return
		// }
		// w.Write(page)
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
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
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

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
		w.Write([]byte("FAILED"))
		return
	}

	//step 2: authentication (token or session/cookies)
	token := GenToken(username)
	ok := dbUser.UpdateToken(username, token)
	if !ok {
		w.Write([]byte("FAILED"))
		return
	}

	//setp 3: redirect to main page after login
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	// http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	resp := util.NewRespMsg(http.StatusTemporaryRedirect, "OK",
		struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	)
	w.Write(resp.JSONBytes())
}

// UserInfoHandler handle user info query
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//step 1: parse form
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	// token := r.Form.Get("token")

	//step 2: validate token
	// added a middleware to validate the token
	// ok := IsTokenValid(token)
	// if !ok {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	//step 3: query
	user, err := dbUser.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//step 4: respond with specific info
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// GenToken get a user token which is used for authentication
func GenToken(username string) string {
	// 40characters md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid check if token is valid
func IsTokenValid(token, username string) bool {
	if len(token) != 40 {
		return false
	}
	//validate if token is expired. The second parameter is base 16, according the GenToken is using base16
	tokenTime, _ := strconv.ParseInt(token[len(token)-8:], 16, 64)
	now := time.Now().Unix()
	if now-tokenTime > validDuration {
		return false
	}
	//query token from tbl_user_token table via username
	dbToken, err := dbUser.GetUserToken(username)
	if err != nil {
		fmt.Printf("fetch token error: %s", err.Error())
		return false
	}

	if dbToken == token {
		return true
	}
	//match two tokens
	return false
}
