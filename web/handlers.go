package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AngelVI13/foos/game"
	"github.com/AngelVI13/foos/log"
	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context, msg string) {
	c.HTML(http.StatusOK, "", views.Page(views.Error(msg)))
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "", views.Page(views.UsersForm()))
}

func usersListHandler(c *gin.Context) {
	playersRawInput := c.PostForm("playersListInput")
	playersSplit := strings.Split(playersRawInput, "\r\n")
	var players []string
	for _, p := range playersSplit {
		// remove X) preffix where X is a number
		// i.e. 1) Ugnius
		if strings.Count(p, ")") == 1 {
			index := strings.Index(p, ")")
			p = p[index+1:]
		}

		player := strings.TrimSpace(p)
		if player == "" {
			continue
		}
		players = append(players, player)
	}

	if !(len(players) == 8 || len(players) == 12) {
		errorHandler(c, fmt.Sprintf("Expected 8 or 12 playes but got %d", len(players)))
		return
	}

	teams, err := game.GenerateTeams(players)
	if err != nil {
		errorHandler(c, fmt.Sprintf("Failed to generate teams: %v", err))
		return
	}
	log.L.Info("", "teams", teams)
	state = NewGlobalState(players, teams)

	c.Redirect(http.StatusFound, TournamentTableUrl)
}

func tournamentBracketHandler(c *gin.Context) {
	log.L.Info("", "rounds", state.Rounds)
	c.HTML(http.StatusOK, "", views.Page(views.Rounds(state.Rounds)))
}
