package main

import (
	"github.com/AngelVI13/foos/views"
	"github.com/AngelVI13/foos/web"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.HTMLRender = &views.TemplRender{}

	web.SetupRoutes(r)
	r.Run(":5555")
}
