package web

import (
	"net/http"

	"github.com/AngelVI13/foos/routes"
	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.StaticFS("/static", http.FS(views.EmbedFs))

	r.GET(routes.IndexUrl, indexHandler)
	r.POST(routes.UsersListUrl, usersListHandler)
	r.GET(routes.TournamentTableUrl, tournamentBracketHandler)
	r.GET(routes.TournamentTableUpdateUrl, tournamentBracketMatchUpdateHandler)
	r.GET(routes.TournamentTableShowUrl, tournamentBracketMatchShowHandler)
	r.POST(routes.TournamentTableShowUrl, tournamentBracketMatchShowHandler)
	r.GET(routes.TournamentTableEndRoundUrl, tournamentBracketEndRoundHandler)
	r.GET(routes.NewTournamentUrl, newTournamentHandler)
}
