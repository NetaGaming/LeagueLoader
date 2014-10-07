package main

import (
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
	"strings"
)

/**
 * Get list of available summoners
 */
func getSummoners(dbmap *gorp.DbMap) <-chan int64 {

	out := make(chan int64)

	// select summoners
	var summoners []int64
	_, err := dbmap.Select(
		&summoners,
		"select id from summoners")
	checkErr(err, "Selecting summoner ids failed")

	go func() {
		for _, n := range summoners {
			out <- n
		}
		close(out)
	}()

	return out
}

// Updates summoner name and level
func updateSummoners(summoners <-chan int64, dbmap *gorp.DbMap) <-chan int64 {

	out := make(chan int64)

	// we'll keep all the summoners details in here
	var SummonersGoRiot map[int64]goriot.Summoner

	// Get summoner data from Riot in batches of
	// forty and then combine them
	go func() {
		var summonerGroup []int64
		for s := range summoners {
			summonerGroup = append(summonerGroup, s)

			if len(summonerGroup) == 40 {
				// get riot data
				riotData, err := goriot.SummonerByID(goriot.NA, summonerGroup...)
				checkErr(err, "Could not load summoners from Riot")

				// add to larger structure
				for k, v := range riotData {
					SummonersGoRiot[k] = v
				}

				// reset slice
				summonerGroup = nil
			}
			close(out)
		}
	}()

	// Build a slice of select queries that will be UNIONd
	// together to help reduce DB calls
	var selectQueries []string
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
        SET s.level = r.level, s.name = r.name;`,
		strings.Join(selectQueries, " UNION "))

	// run query and check for errors
	_, err := dbmap.Exec(updateQuery)
	checkErr(err, "Summoner update failed")

	return out
}
