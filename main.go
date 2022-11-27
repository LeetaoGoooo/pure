package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"pure/core"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/githubv4"
	"gopkg.in/yaml.v2"
)

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Category struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
}

type PageQuery struct {
	Next         string `form:"next,omitempty"`
	Pre          string `form:"pre,omitempty"`
	CategoryId   string `uri:"category_id"`
	CategoryName string `uri:"category_name"`
}

type PostQuery struct {
	Id    uint64 `uri:"id" binding:"required"`
	Title string `uri:"title" binding:"required"`
}

type PureConfig struct {
	UserName    string     `yaml:"username"`
	Repo        string     `yaml:"repo"`
	RepoId      string     `yaml:"repoId"`
	AccessToken string     `yaml:"accessToken"`
	Categories  []Category `yaml:"categories"`
}

var config PureConfig
var api core.BlogApi

func init() {
	configFile, err := os.ReadFile("pure.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		panic(err)
	}

	api = core.NewApi(config.UserName, config.Repo, config.AccessToken)

}

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
	"unescapeHtml": func(bodyHtml githubv4.HTML) template.HTML {
		return template.HTML(string(bodyHtml))
	},
}

func FetchPosts(c *gin.Context) {
	var pageQuery PageQuery
	if err := c.ShouldBind(&pageQuery); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindUri(&pageQuery); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	categoryId := config.Categories[0].Id

	if len(pageQuery.CategoryId) > 0 {

		categoryId = pageQuery.CategoryId
	}

	discussions, err := api.FetchPosts(pageQuery.Pre, pageQuery.Next, categoryId)

	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", map[string]any{
		"Posts":   discussions,
		"Navbars": config.Categories,
	})
}

func FetchPost(c *gin.Context) {
	var postQuery PostQuery
	if err := c.ShouldBindUri(&postQuery); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	discussion, err := api.FetchPost(postQuery.Id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "post.html", map[string]any{
		"Post":    discussion,
		"Navbars": config.Categories,
		"Repo":    fmt.Sprintf("%s/%s", config.UserName, config.Repo),
		"RepoId":  config.RepoId,
	})
}

func main() {

	r := gin.Default()
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/css", "templates/css")
	r.GET("/", FetchPosts)
	r.GET("/category/:category_id/:category_name", FetchPosts)
	r.GET("/post/:id/:title", FetchPost)
	r.Run()
}
