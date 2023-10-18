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

	r.StaticFS("/static", http.FS(views.EmbedFs))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.Page(views.UsersForm()))
	})

	log.Info("Running")

	r.Run(":5555")
}
