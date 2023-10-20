package web

import (
	"net/http"

	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.StaticFS("/static", http.FS(views.EmbedFs))

	r.GET(IndexUrl, indexHandler)
	r.POST(UsersListUrl, usersListHandler)
	r.GET(TournamentTableUrl, tournamentBracketHandler)
}
