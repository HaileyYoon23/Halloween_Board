package content

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"runtime"
	"strconv"
)

var (
	//go:embed *.gohtml
	fs  embed.FS
	tpl = template.Must(template.ParseFS(fs, "*.gohtml"))
)

func ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	statusCode := http.StatusOK
	b := new(bytes.Buffer)

	err := tpl.ExecuteTemplate(b, name+".gohtml", data)
	if err != nil {
		panic(err)
	}

	h := w.Header()
	h.Set("Content-Type", "text/html;charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(b.Len()))
	w.WriteHeader(statusCode)
	_, _ = b.WriteTo(w)
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
