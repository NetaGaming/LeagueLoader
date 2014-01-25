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
func getSummoners(dbmap *gorp.DbMap) (summoners []int64) {

	// select summoners
	_, err := dbmap.Select(
		&summoners,
		"select id from summoners")
	checkErr(err, "Selecting summoner ids failed")

	return summoners
}

/***
 * Updates summoner name and level
 */
func updateSummoners(summoners []int64, dbmap *gorp.DbMap) {

	// get summoner data
	riotData, err := goriot.SummonerByID(goriot.NA, summoners...)
	checkErr(err, "Could not load summoners from Riot")

	// Build a slice of select queries that will be UNIONd
	// together to help reduce DB calls
	var selectQueries []string
	for _, summoner := range riotData {
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
	_, err = dbmap.Exec(updateQuery)
	checkErr(err, "Summoner update failed")

	return
}
