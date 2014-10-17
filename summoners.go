package main

import (
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
	"strings"
)

// Get list of available summoners
func getSummoners(dbmap *gorp.DbMap) (
	<-chan SummonerInfo, <-chan SummonerInfo) {

	summonerChan1 := make(chan SummonerInfo)
	summonerChan2 := make(chan SummonerInfo)

	// select summoners
	var summoners []SummonerInfo
	_, err := dbmap.Select(
		&summoners,
		"select id from summoners")
	checkErr(err, "Selecting summoner ids failed")

	globalWg.Add(1)
	go func() {
		for _, n := range summoners {
			summonerChan1 <- n
			summonerChan2 <- n
		}
		close(summonerChan1)
		close(summonerChan2)
		globalWg.Done()
	}()

	return summonerChan1, summonerChan2
}

// Updates summoner name and level
func updateSummoners(summoners <-chan SummonerInfo, dbmap *gorp.DbMap) {

	// we'll keep all the summoners details in here
	var SummonersGoRiot map[int64]goriot.Summoner = make(map[int64]goriot.Summoner)

	// Get summoner data from Riot in batches of
	// forty and then combine them
	go func() {
		globalWg.Add(1)
		var selectQueries []string
		for s := range summoners {

			// get riot data
			riotData, err := goriot.SummonerByID(goriot.NA, s.ID)
			checkErr(err, "Could not load summoners from Riot")

			// add to larger structure
			for k, v := range riotData {
				SummonersGoRiot[k] = v
			}

		}

		// Build a slice of select queries that will be UNIONd
		// together to help reduce DB calls
		for _, summoner := range SummonersGoRiot {
			selectQueries = append(
				selectQueries,
				fmt.Sprintf(
					"SELECT %d id, %d level, '%s' name",
					summoner.ID,
					summoner.SummonerLevel,
					summoner.Name))
		}

		// build final update query
		var updateQuery string = fmt.Sprintf(
			`UPDATE summoners s INNER JOIN (
				%s
			) r USING(id)
			SET s.level = r.level, s.name = r.name, s.last_update = UTC_TIMESTAMP();`,
			strings.Join(selectQueries, " UNION "))

		// run query and check for errors
		_, err := dbmap.Exec(updateQuery)
		checkErr(err, "Summoner update failed")
		globalWg.Done()
	}()

}
