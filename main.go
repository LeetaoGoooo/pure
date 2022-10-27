package main

import (
	"log"
	"net/http"
	"os"
	"pure/core"
	"strconv"
	"strings"
	"text/template"

	"github.com/shurcooL/githubv4"
)

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

var api = core.NewApi(os.Getenv("GITHUB_USER_NAME"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_ACCESS_TOKEN"))

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

func main() {

	http.HandleFunc("/", FetchPosts)
	http.HandleFunc("/article/", FetchPost)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./templates/js/"))))
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
