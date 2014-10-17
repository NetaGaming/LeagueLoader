package main

import (
	"fmt"

	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
)

type GameInfo struct {
	SummonerID int64
	Game       goriot.Game
}

// use list of summoners to download game data
func updateGames(summoners <-chan SummonerInfo, dbmap *gorp.DbMap) (gameCount int) {

	var out int = 0

	for s := range summoners {

		summonerGames, riotErr := goriot.RecentGameBySummoner(goriot.NA, s.ID)
		checkErr(
			riotErr,
			fmt.Sprintf(
				"Unable to get summoner's (%d) recent games",
				s.ID))

		// dump game info into channel
		globalWg.Add(1)
		go func() {
			gameChan1 := make(chan GameInfo)

			for _, game := range summonerGames {
				gameChan1 <- GameInfo{s.ID, game}
			}

			close(gameChan1)

			// Update the shared game information
			// TODO: these two methods aren't working here
			gameChan2 := updateGameInfo(gameChan1, dbmap)
			updateSummonerGames(gameChan2, dbmap)

			globalWg.Done()
		}()

	}

	return out
}

// Updates the common game information
func updateGameInfo(games <-chan GameInfo, db *gorp.DbMap) <-chan GameInfo {
	// TODO: Use insert/update query instead
	var gameInfoQuery string = `
                            INSERT IGNORE INTO game_info
                                (id, mode, type, subType, mapId, date)
                            VALUES
                                (?, ?, ?, ?, ?, FROM_UNIXTIME(?))`

	out := make(chan GameInfo)
	globalWg.Add(1)
	go func() {

		for gi := range games {
			_, infoErr := db.Exec(
				gameInfoQuery,
				gi.Game.GameID,
				gi.Game.GameMode,
				gi.Game.GameType,
				gi.Game.SubType,
				gi.Game.MapID,
				gi.Game.CreateDate/1000)
			checkErr(infoErr, "Unable to insert new game info")

			out <- gi
		}
		close(out)

		globalWg.Done()
	}()

	return out
}

// Updates a summoners specific game information
func updateSummonerGames(
	gameStats <-chan GameInfo,
	db *gorp.DbMap) {

	// save summoner_game
	// TODO: Use insert/update query instead
	var summonerGameQuery string = `
                        INSERT IGNORE INTO summoner_games
                            (summonerId,
							 gameId,
							 championId,
							 spellOne,
							 spellTwo,
							 minionsKilled,
							 numDeaths,
							 assists,
							 championsKilled,
                             won)
                        VALUES
                            (?,?,?,?,?,?,?,?,?,?)`
	globalWg.Add(1)
	go func() {

		for gs := range gameStats {
			_, sgErr := db.Exec(
				summonerGameQuery,
				gs.SummonerID,
				gs.Game.GameID,
				gs.Game.ChampionID,
				gs.Game.Spell1,
				gs.Game.Spell2,
				gs.Game.Statistics.MinionsKilled,
				gs.Game.Statistics.NumDeaths,
				gs.Game.Statistics.Assists,
				gs.Game.Statistics.ChampionsKilled,
				gs.Game.Statistics.Win)
			checkErr(sgErr, "Could not save summoner game info")
		}

		globalWg.Done()
	}()
}
