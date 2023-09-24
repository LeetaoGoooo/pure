package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"pure/constants"
	"pure/core"
	"pure/enties"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/shurcooL/githubv4"
)

var config enties.PureConfig = constants.BlogConfig

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

var api core.BlogApi

func init() {
	api = core.NewApi(config.UserName, config.Repo, config.AccessToken)
}

func FetchPosts(c *gin.Context) {
	var pageQuery enties.PageQuery
	if err := c.ShouldBind(&pageQuery); err != nil {
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
			"Message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindUri(&pageQuery); err != nil {
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
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
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
			"Message": err.Error(),
		})
		return
	}

	tmpl, _ := template.New("posts").Funcs(funcMap).Parse(constants.IndexPage)

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, map[string]any{
		"Posts":   discussions,
		"Navbars": config.Categories,
		"About":   config.About,
	}); err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(tpl.Bytes())
}

func FetchPost(c *gin.Context) {

	var postQuery enties.PostQuery
	if err := c.ShouldBindUri(&postQuery); err != nil {
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
			"Message": err.Error(),
		})
		return
	}

	discussion, err := api.FetchPost(postQuery.Id)
	if err != nil {
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
			"Message": err.Error(),
		})
		return
	}

	tmpl, _ := template.New("post").Funcs(funcMap).Parse(constants.PostPage)
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, map[string]any{
		"Post":    discussion,
		"Navbars": config.Categories,
		"About":   config.About,
		"Repo":    fmt.Sprintf("%s/%s", config.UserName, config.Repo),
		"RepoId":  config.RepoId,
	}); err != nil {
		panic(err)
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(tpl.Bytes())
}

func AboutPage(c *gin.Context) {

	discussion, err := api.FetchPost(config.About)
	if err != nil {
		c.HTML(http.StatusBadRequest, constants.ErrorPage, map[string]any{
			"Message": err.Error(),
		})
		return
	}

	tmpl, _ := template.New("post").Funcs(funcMap).Parse(constants.PostPage)
	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, map[string]any{
		"Post":    discussion,
		"Navbars": config.Categories,
		"About":   config.About,
		"Repo":    fmt.Sprintf("%s/%s", config.UserName, config.Repo),
		"RepoId":  config.RepoId,
	}); err != nil {
		panic(err)
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(tpl.Bytes())
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
