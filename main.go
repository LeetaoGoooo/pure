package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"pure/core"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/shurcooL/githubv4"
)

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

var api = core.NewApi(os.Getenv("GITHUB_USER_NAME"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_ACCESS_TOKEN"))
var storage core.Storage = *core.NewStorage()

var funcMap = template.FuncMap{
	"formatDate": func(unformated githubv4.DateTime) string {
		return unformated.Time.Format("2006-01-02")
	},
	"previewContent": func(fullContent githubv4.String) string {
		if len(fullContent) >= 100 {
			return string(fullContent)[0:100]
		}
		return string(fullContent)
	},
}

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer r.Body.Close()
	next := r.URL.Query().Get("next")
	pre := r.URL.Query().Get("pre")

	discussions, err := api.FetchPosts(pre, next)

	if err != nil {
		redictTo404(w, Response[core.Discussions]{
			Code:    400,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	indexTemplate, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/base/navbar.html", "templates/base/footer.html", "templates/index.html")
	if err != nil {
		redictTo404(w, Response[core.Discussions]{
			Code:    400,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	indexTemplate.Execute(w, discussions)
}

func redictTo404(w http.ResponseWriter, r Response[core.Discussions]) {
	errTemplate := template.Must(template.ParseFiles("templates/error.html"))
	errTemplate.Execute(w, r)
}

func FetchPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer r.Body.Close()
	idAndTitle := strings.TrimPrefix(r.URL.Path, "/article/")
	idAndTitleArr := strings.Split(idAndTitle, "/")
	number, err := strconv.Atoi(idAndTitleArr[0])
	if err != nil {
		redictTo404(w, Response[core.Discussions]{
			Code:    400,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	discussion, err := api.FetchPost(number)
	if err != nil {
		redictTo404(w, Response[core.Discussions]{
			Code:    400,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	postTemplate, err := template.New("post.html").Funcs(funcMap).ParseFiles("templates/base/navbar.html", "templates/base/footer.html", "templates/post.html")
	if err != nil {
		redictTo404(w, Response[core.Discussions]{
			Code:    400,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	postTemplate.Execute(w, discussion)
}

// Cached 缓存页面
func Cached(duration string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		content := storage.Get(r.RequestURI)
		if content != nil {
			w.Write(content)
		} else {
			c := httptest.NewRecorder()
			handler(c, r)

			for k, v := range c.Header() {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.Bytes()

			if d, err := time.ParseDuration(duration); err == nil {
				storage.Set(r.RequestURI, content, d)
			}

			w.Write(content)
		}

	})
}

func main() {

	http.Handle("/", Cached("10m", FetchPosts))
	http.Handle("/article/", Cached("1h", FetchPost))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./templates/js/"))))
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
