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
	TournamentTableShowUrl   = "/tournament_table/:match_id/:team/show"
)

func RouteWithParam(route, name, value string) string {
	return strings.Replace(route, fmt.Sprintf(":%v", name), value, 1)
}

func MakeMatchUpdateUrl(match *game.Match, teamNum int) string {
	url := RouteWithParam(TournamentTableUpdateUrl, "match_id", match.Id)
	url = RouteWithParam(url, "team", fmt.Sprint(teamNum))
	return url
}

func MakeMatchShowUrl(match *game.Match, teamNum int) string {
	url := RouteWithParam(TournamentTableShowUrl, "match_id", match.Id)
	url = RouteWithParam(url, "team", fmt.Sprint(teamNum))
	return url
}
