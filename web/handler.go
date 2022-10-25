package web

import (
	"errors"
	"fmt"
	"github.com/haileyyoon23/christmas-wish-box/content"
	"github.com/haileyyoon23/christmas-wish-box/db"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type templateStruct map[string]interface{}

var (
	ErrEmptyGift = errors.New("blank is not gift")
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	talkList, err := db.GetTalks(db.DB)

	if err != nil {
		panic(err)
	}

	errMsg := r.URL.Query().Get("errMsg")

	content.ExecuteTemplate(w, "home", templateStruct{
		"talkList": talkList,
		"errMsg":   errMsg,
	})
}

func ListPage(w http.ResponseWriter, r *http.Request) {
	talkList, err := db.GetTalks(db.DB)

	if err != nil {
		panic(err)
	}

	errMsg := r.URL.Query().Get("errMsg")

	content.ExecuteTemplate(w, "list", templateStruct{
		"talkList": talkList,
		"errMsg":   errMsg,
	})
}

func TalkAppendHandler(w http.ResponseWriter, r *http.Request) {
	talk := r.URL.Query().Get("talk")
	talkTemp := talk

	var path string
	var err error

	if strings.Trim(talkTemp, " ") == "" {
		err = ErrEmptyGift
	}
	if err != nil {
		goto ERRORS
	}

	err = db.AddTalk(db.DB, talk)
	if err != nil {
		goto ERRORS
	}

ERRORS:
	if err != nil {
		path = "?errMsg=" + err.Error()
	}

	http.Redirect(w, r, "/list"+path, http.StatusSeeOther)
}

func TalkLikeHandler(w http.ResponseWriter, r *http.Request) {
	talk := r.URL.Query().Get("talk")

	err := db.UpdateLike(talk)
	if err != nil {
		panic(err)
	}
}

func GiftDislikeHandler(w http.ResponseWriter, r *http.Request) {
	gift := r.URL.Query().Get("present")

	err := db.UpdateDislike(gift)
	if err != nil {
		panic(err)
	}
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				renderInternalError(w, e)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func renderInternalError(w http.ResponseWriter, err interface{}) {
	_, file, line, _ := runtime.Caller(3)

	b := []byte(fmt.Sprintf("ERROR\n\n%s:%d\n\n%s\n", path.Base(file), line, err))

	h := w.Header()
	h.Set("Content-Type", "text/plain;charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(b)
}
