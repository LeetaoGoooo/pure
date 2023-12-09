package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"pure/core"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
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
	UserName string `yaml:"username"`
	Repo     string `yaml:"repo"`
	RepoId   string `yaml:"repoId"`
	Website  struct {
		Host  string `yaml:"host"`
		Bio   string `yaml:"bio"`
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	} `yaml:"website"`
	AccessToken string     `yaml:"accessToken"`
	Categories  []Category `yaml:"categories"`
	About       uint64     `yaml:"about"`
}

var config PureConfig
var api core.BlogApi
var memoryCache persist.CacheStore

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
	memoryCache = persist.NewMemoryStore(1 * time.Minute)
}

var funcMap = template.FuncMap{
	"formatDate": func(unformated githubv4.DateTime) string {
		return unformated.Time.Format("2006-01-02")
	},
	"previewContent": func(fullContent githubv4.String) string {
		if len([]rune(fullContent)) >= 100 {
			return string([]rune(fullContent)[:100])
		}
		return string(fullContent)
	},
	"unescapeHtml": func(bodyHtml githubv4.HTML) template.HTML {
		return template.HTML(string(bodyHtml))
	},
	"isExisted": func(categoryId githubv4.String) bool {
		for _, category := range config.Categories {
			if category.Id == string(categoryId) {
				return true
			}
		}
		return false
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

	categoryId := ""

	if len(pageQuery.CategoryId) > 0 {
		categoryId = pageQuery.CategoryId
	}

	// 获取所有文章
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
		"About":   config.About,
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
		"About":   config.About,
		"Repo":    fmt.Sprintf("%s/%s", config.UserName, config.Repo),
		"RepoId":  config.RepoId,
	})
}

func AboutPage(c *gin.Context) {
	discussion, err := api.FetchPost(config.About)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "post.html", map[string]any{
		"Post":    discussion,
		"Navbars": config.Categories,
		"About":   config.About,
		"Repo":    fmt.Sprintf("%s/%s", config.UserName, config.Repo),
		"RepoId":  config.RepoId,
	})
}

func GenerateFeed(c *gin.Context) {
	feed := &feeds.Feed{
		Title:       config.Website.Name,
		Link:        &feeds.Link{Href: config.Website.Host},
		Description: config.Website.Bio,
		Author:      &feeds.Author{Name: config.Website.Name, Email: config.Website.Email},
		Created:     time.Now(),
	}

	discussions, err := api.FetchPosts("", "", config.Categories[0].Id)

	if err != nil {
		c.XML(http.StatusOK, gin.H{
			"Message": fmt.Sprintf("Something seems error, please contact %s(%s)", config.Website.Name, config.Website.Email),
		})
		return
	}

	for _, disdiscussion := range discussions.Nodes {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       string(disdiscussion.Title),
			Description: string([]rune(disdiscussion.Body)[:200]),
			Author:      &feeds.Author{Name: config.Website.Name, Email: config.Website.Email},
			Created:     disdiscussion.CreatedAt.Time,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d/%s", config.Website.Host, disdiscussion.Number, disdiscussion.Title)},
		})
	}

	feed.WriteAtom(c.Writer)
}

func main() {

	r := gin.Default()
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/css", "templates/css")
	r.StaticFile("/favicon.ico", "templates/favicon.ico")
	r.GET("/", cache.CacheByRequestURI(memoryCache, 30*time.Second), FetchPosts)
	r.GET("/category/:category_id/:category_name", cache.CacheByRequestURI(memoryCache, 30*time.Second), FetchPosts)
	r.GET("/post/:id/:title", cache.CacheByRequestURI(memoryCache, 1*time.Hour), FetchPost)
	r.GET("/atom.xml", cache.CacheByRequestURI(memoryCache, 24*time.Hour), GenerateFeed)
	r.GET("/404", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "error.html", nil)
	})
	if config.About > 0 {
		r.GET("/about", cache.CacheByRequestURI(memoryCache, 1*time.Hour), AboutPage)
	}
	r.Run()
}
