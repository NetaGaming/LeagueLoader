/*
   LeagueLoader is an application that loads data from the Riot
   API into a database for specific users. The data can then be used
   to track their progression for any number of things: tournaments,
   personal statistics, or achievements
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"os"
	"time"
)

var dtFormat string = "2006-01-02 15:04:05"

/* Config elements */
type Configuration struct {
	ApiKey   string      `json:"apiKey"`
	DbConfig MysqlConfig `json:"mysqlConfig"`
}

type MysqlConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

// Streamlines checking for errors
func checkErr(e error, message string) {
	if e != nil {
		log.Fatalln(message, e)
	}
}

func main() {

	// set loader start time
	var startTime string = time.Now().Format(dtFormat)

	var config Configuration = openAndReadConfig("config.json")
	var dbConfig MysqlConfig = config.DbConfig

	// Make connection
	dbmap := initDb(dbConfig.Database, dbConfig.Username, dbConfig.Password)
	defer dbmap.Db.Close()

	// Goriot setup
	// TODO: move limits to configuration
	goriot.SetAPIKey(config.ApiKey)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	goriot.SetLongRateLimit(500, 10*time.Minute)

	// get list of available summoners
	var summoners []int64 = getSummoners(dbmap)

	// update summoner information
	updateSummoners(summoners, dbmap)
	fmt.Println("Summoners updated")

	// update game information
	var updatedGameCount int = updateGames(summoners, dbmap) - 1
	fmt.Println("Games updated: ", updatedGameCount)

	// end loader time and save
	var endTime string = time.Now().Format(dtFormat)
	saveLoadReport(startTime, endTime, updatedGameCount, dbmap)

	return
}

/***
 * Opens configuration files and
 * implements associated structs
 */
func openAndReadConfig(configFileName string) (config Configuration) {

	// load config file
	configFile, err := os.Open(configFileName)
	checkErr(err, "Unable to open config file")

	// parse config file
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	checkErr(err, "Unable to decode json")

	return config
}

/***
 * Sets up the "ORM" by connecting to db and mapping
 * structs to tables
 */
func initDb(database string, username string, password string) *gorp.DbMap {
	db, err := sql.Open("mymysql", database+"/"+username+"/"+password)
	checkErr(err, "Connection failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "utf-8"}}

	return dbmap
}

/***
 * Checks a slice of int64 for a given
 * value
 */
func existsInSlice(search int64, values []int64) (exists bool) {

	for _, value := range values {
		if value == search {
			return true
		}
	}

	return false
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
