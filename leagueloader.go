/*
   LeagueLoader is an application that loads data from the Riot
   API into a database for specific users. The data can then be used
   to track their progression for any number of things: tournaments,
   personal statistics, or achievements
*/

package main

import (
	"database/sql"
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"runtime"
	"sync"
	"time"
)

var dtFormat string = "2006-01-02 15:04:05"
var globalWg sync.WaitGroup

// We'll fill this from the database and pass it
// around where ever we're updating Summoner info
type SummonerInfo struct {
	ID int64 `db:"id"`
}

// Streamlines checking for errors
func checkErr(e error, message string) {
	if e != nil {
		log.Fatalln(message, e)
	}
}

func main() {

	// set loader start time
	//var startTime string = time.Now().Format(dtFormat)

	runtime.GOMAXPROCS(runtime.NumCPU())

	var config Configuration = openAndReadConfig("config.json")
	var dbConfig MysqlConfig = config.DbConfig

	// Make connection
	dbmap := initDb(dbConfig.Database, dbConfig.Username, dbConfig.Password)
	defer dbmap.Db.Close()

	// Goriot setup
	// TODO: move limits to configuration
	goriot.SetAPIKey(config.ApiKey)
	goriot.SetSmallRateLimit(3000, 10*time.Second)
	goriot.SetLongRateLimit(180000, 10*time.Minute)

	// get channel that streams summoner ids
	//summoners, gameSummoners := getSummoners(dbmap)
	summoners1, summoners2 := getSummoners(dbmap)

	// update summoner information
	updateSummoners(summoners1, dbmap)

	// update game information
	updateGames(summoners2, dbmap)

	// end loader time and save
	//var endTime string = time.Now().Format(dtFormat)
	//saveLoadReport(startTime, endTime, updatedGameCount, dbmap)

	globalWg.Wait()

	return
}

// Creates a connection to a MySQL database
func initDb(database string, username string, password string) *gorp.DbMap {
	db, err := sql.Open("mymysql", database+"/"+username+"/"+password)
	checkErr(err, "Connection failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "utf-8"}}

	return dbmap
}

// Saves runttime report to db
func saveLoadReport(StartTime string, EndTime string, Records int, dbmap *gorp.DbMap) {
	var reportQuery string = `INSERT INTO runtimes
            (startTime, endTime, records)
         VALUES
            (?,?,?)`
	_, err := dbmap.Exec(reportQuery,
		StartTime,
		EndTime,
		Records)

	checkErr(err, "Could not report to database")

	fmt.Println("Reported")
}
