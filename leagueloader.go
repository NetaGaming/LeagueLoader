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
	"github.com/yvasiyarov/gorelic"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"os"
	"time"
)

/* Config elements */
type Configuration struct {
	ApiKey   string      `json:"apiKey"`
	NewRelic string      `json:"newRelicKey"`
	DbConfig MysqlConfig `json:"mysqlConfig"`
}

type MysqlConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

/* Database tables */
// so far, unused
type Summoner struct {
	Id       int64
	Name     string
	RealName string `db:"real_name"`
	TeamId   int    `db:"neta_team"`
	Level    int
}

// Streamlines checking for errors
func checkErr(e error, message string) {
	if e != nil {
		log.Fatalln(message, e)
	}
}

func main() {

	var config Configuration = openAndReadConfig("config.json")
	var dbConfig MysqlConfig = config.DbConfig

	// Goriot setup
	goriot.SetAPIKey(config.ApiKey)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	goriot.SetLongRateLimit(500, 10*time.Minute)

	// New Relic setup
	agent := gorelic.NewAgent()
	agent.NewrelicLicense = config.NewRelic
	agent.NewrelicName = "League Loader"
	agent.Run()

	// Make connection
	dbmap := initDb(dbConfig.Database, dbConfig.Username, dbConfig.Password)
	defer dbmap.Db.Close()

	// get list of available summoners
	var summoners []int64 = getSummoners(dbmap)

	// update summoner information
	updateSummoners(summoners, dbmap)
	fmt.Println("Summoners updated")

	// update game information
	updateGames(summoners, dbmap)
	fmt.Println("Games updated")

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

	dbmap.AddTableWithName(Summoner{}, "summoners").SetKeys(false, "Id")

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
