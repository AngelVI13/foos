package routes

import (
	"fmt"
	"strings"

	"github.com/AngelVI13/foos/game"
)

const (
	IndexUrl                 = "/"
	UsersListUrl             = "/users_list"
	TournamentTableUrl       = "/tournament_table"
	TournamentTableUpdateUrl = "/tournament_table/:match_id/:team/update"
)

func RouteWithParam(route, name, value string) string {
	return strings.Replace(route, fmt.Sprintf(":%v", name), value, 1)
}

func MakeMatchUpdateUrl(match *game.Match, teamNum int) string {
	url := RouteWithParam(TournamentTableUpdateUrl, "match_id", match.Id)
	url = RouteWithParam(url, "team", fmt.Sprint(teamNum))
	return url
}
