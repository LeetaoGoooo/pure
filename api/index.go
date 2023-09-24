package api

import (
	"bytes"
	"html/template"
	"net/http"
	"pure/constants"
	"pure/handler"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

var memoryCache persist.CacheStore

func Handler(writer http.ResponseWriter, request *http.Request) {

	memoryCache = persist.NewMemoryStore(10 * time.Minute)

	r := gin.New()
	r.GET("/posts", cache.CacheByRequestURI(memoryCache, 30*time.Second), handler.FetchPosts)
	r.GET("/category/:category_id/:category_name", cache.CacheByRequestURI(memoryCache, 30*time.Second), handler.FetchPosts)
	r.GET("/post/:id/:title", cache.CacheByRequestURI(memoryCache, 1*time.Hour), handler.FetchPost)
	r.GET("/atom.xml", cache.CacheByRequestURI(memoryCache, 24*time.Hour), handler.GenerateFeed)
	r.GET("/404", func(ctx *gin.Context) {
		tmpl, _ := template.New("404").Parse(constants.ErrorPage)
		var tpl bytes.Buffer
		if err := tmpl.Execute(&tpl, nil); err != nil {
			panic(err)
		}
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(tpl.Bytes())

	})
	r.GET("/about", cache.CacheByRequestURI(memoryCache, 1*time.Hour), handler.AboutPage)
	r.ServeHTTP(writer, request)
}
