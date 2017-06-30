package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/imgasm/server/config"
	"github.com/spf13/viper"
	"net/http"
)

var (
	sessionStore *sessions.CookieStore
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(config.ViperConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s \n", err))
	}
	sessionStore = sessions.NewCookieStore(
		[]byte(viper.GetString("session.keys.authentication")),
		[]byte(viper.GetString("session.keys.encryption")))
	env := viper.GetString("environment")
	switch env {
	case "dev":
		sessionStore.Options = &sessions.Options{
			Domain: viper.GetString("session.local.domain"),
			Path:   viper.GetString("session.local.path"),
			MaxAge: viper.GetInt("session.local.max_age"),
		}
	case "prod":
		sessionStore.Options = &sessions.Options{
			Domain: viper.GetString("session.production.domain"),
			Path:   viper.GetString("session.production.path"),
			MaxAge: viper.GetInt("session.production.max_age"),
		}
	default:
		panic("environment not set")
	}
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := sessionStore.Get(r, "imgasm-session")
	return session
}

func ClearAllSessionValues(sess *sessions.Session) {
	for key := range sess.Values {
		delete(sess.Values, key)
	}
}

func ClearSessionByKey(key interface{}, sess *sessions.Session) {
	delete(sess.Values, key)
}
