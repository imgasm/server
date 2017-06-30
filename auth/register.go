package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/render"
	"github.com/gorilla/sessions"
	"github.com/imgasm/server/session"
	"github.com/imgasm/server/user"
	"io/ioutil"
	// "log"
	"net/http"
	"regexp"
)

type RegistrationForm struct {
	Username string
	Password string
}

// curl -i -X POST -d "{\"username\": \"awesomeusername\", \"password\": \"secretpassword\"}" http://localhost:8080/auth/register
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)

	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 15 {
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
	var form RegistrationForm
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

	_, err = user.GetUserByNameKey(ctx, form.Username)
	if err == nil {
		registerAttempt(sess)
		sess.Save(r, w)
		http.Error(w, "Username already taken", 409)
		return
	}

	_, err = user.PutUserByNameKey(ctx, form.Username, []byte(form.Password))
	if err != nil {
		http.Error(w, "Failed to register user, please try again", 422)
		return
	}

	session.ClearAllSessionValues(sess)
	sess.Values["username"] = form.Username
	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to store cookies.", 500)
		return
	}

	render.JSON(w, r, form)
}

func validateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("Username required.")
	}
	if len(username) > 20 {
		return errors.New("Username can at most be 20 characters.")
	}
	if regMatch, _ := regexp.MatchString("^[A-Za-z0-9-]+$", username); !regMatch {
		return errors.New("Username can only contain letters, numbers, and hyphenation.")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) == 0 {
		return errors.New("Password required.")
	}
	if len(password) > 128 {
		return errors.New("Password can at most be 128 characters.")
	}
	return nil
}

func registerAttempt(sess *sessions.Session) {
	if sess.Values["register_attempt"] == nil {
		sess.Values["register_attempt"] = 1
	} else {
		sess.Values["register_attempt"] = sess.Values["register_attempt"].(int) + 1
	}
}
