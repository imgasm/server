package auth

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/render"
	"github.com/gorilla/sessions"
	"github.com/imgasm/server/session"
	"github.com/imgasm/server/user"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	// "log"
	"net/http"
)

type LoginForm struct {
	Username string
	Password string
}

// curl -i -X POST -d "{\"username\": \"awesomeusername\", \"password\": \"secretpassword\"}" http://localhost:8080/auth/login
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)

	if sess.Values["login_attempt"] != nil && sess.Values["login_attempt"].(int) >= 15 {
		http.Error(w, "In order to protect ourselves against spam, we only allow a limited amount of erroneous requests. Please try again later.", 429)
		return
	}

	if sess.Values["username"] != nil {
		http.Error(w, "Already logged in. If you wish to change user, please logout first.", 409)
		return
	}

	// read json request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request. Please ensure the data is sent correctly (see official API documentation) and try again.", 422)
		// todo: use 422 or 500? improve error message
	}
	var form LoginForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		http.Error(w, "Failed to parse JSON-encoded data. Please ensure the data is sent correctly (see official API documentation) and try again.", 422)
		// todo: use 422 or 500? improve error message
	}

	err = validateUsername(form.Username)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	err = validatePassword(form.Password)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	userStruct, err := user.GetUserByNameKey(ctx, form.Username)
	if err != nil {
		loginAttempt(sess)
		sess.Save(r, w)
		http.Error(w, "Username does not exist.", 404)
		return
	}

	err = bcrypt.CompareHashAndPassword(userStruct.Password, []byte(form.Password))
	if err != nil {
		loginAttempt(sess)
		sess.Save(r, w)
		http.Error(w, "Password does not match.", 422)
		return
	}

	session.ClearAllSessionValues(sess)
	sess.Values["username"] = form.Username
	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to store cookies.", 500)
		return
	}

	render.JSON(w, r, userStruct)
}

func loginAttempt(sess *sessions.Session) {
	if sess.Values["login_attempt"] == nil {
		sess.Values["login_attempt"] = 1
	} else {
		sess.Values["login_attempt"] = sess.Values["login_attempt"].(int) + 1
	}
}
