package main

import (
	"net/http"
	"os"

	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := gin.Default()
	r.HTMLRender = &views.TemplRender{}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.Page(global.Count))
	})

	r.POST("/", func(c *gin.Context) {
		v := c.PostForm("global")
		if v == "global" {
			global.Count++
		}
		log.Info("", "global", v)
		c.HTML(http.StatusOK, "", views.Page(global.Count))
	})
	log.Info("Running")

	r.Run(":5555")
}
