package main

import (
	"fmt"

	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
)

type GameInfo struct {
	SummonerID int64
	Game       goriot.Game
	StatID     int64
}

// use list of summoners to download game data
// TODO: use a channel to ouput summoner IDs as they finish
func updateGames(summoners <-chan SummonerInfo, dbmap *gorp.DbMap) (gameCount int) {

	var out int = 0
	gameChan1 := make(chan GameInfo)
	gameChan2 := make(chan GameInfo)

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
			for _, game := range summonerGames {
				gameChan1 <- GameInfo{s.ID, game, 0}
				gameChan2 <- GameInfo{s.ID, game, 0}
			}
			globalWg.Done()
		}()

		// Update the shared game information
		updateGameInfo(gameChan1, dbmap)
		gameStats := updateSummonerStatistics(gameChan2, dbmap)

		updateSummonerGames(gameStats, dbmap)
	}

	return out
}

// Updates the common game information
func updateGameInfo(games <-chan GameInfo, db *gorp.DbMap) {
	// TODO: Use insert/update query instead
	var gameInfoQuery string = `
                            INSERT IGNORE INTO game_info
                                (id, mode, type, subType, mapId, date)
                            VALUES
                                (?, ?, ?, ?, ?, FROM_UNIXTIME(?))`

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
		}

		globalWg.Done()
	}()
}

// Updates a summoners specific game information
func updateSummonerGames(
	gameStats <-chan GameInfo,
	db *gorp.DbMap) {

	// save summoner_game
	// TODO: Use insert/update query instead
	var summonerGameQuery string = `
                        INSERT INTO summoner_games
                            (summonerId, gameId, championId, spellOne, spellTwo,
                             statId, won)
                        VALUES
                            (?,?,?,?,?,?,?)`
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
				gs.StatID,
				gs.Game.Statistics.Win)
			checkErr(sgErr, "Could not save summoner game info")
		}

		globalWg.Done()
	}()

}

func updateSummonerStatistics(
	gameInfo <-chan GameInfo,
	db *gorp.DbMap) <-chan GameInfo {
	// save stats
	var statsQuery string = `
                        INSERT INTO summoner_stats
                            (assists, barracksKilled, championsKilled,
                             combatPlayerScore, consumablesPurchased,
                             damageDealtPlayer, doubleKills, firstBlood,
                             gold, goldEarned, goldSpent, item0, item1,
                             item2, item3, item4, item5, item6, itemsPurchased,
                             killingSprees, largestCriticalStrike,
                             largestKillingSpree, largestMultiKill,
                             legendaryItemsCreated, level, magicDamageDealtPlayer,
                             magicDamageDealtToChampions, magicDamageTaken,
                             minionsDenied, minionsKilled, neutralMinionsKilled,
                             neutralMinionsKilledEnemyJungle,
                             neutralMinionsKilledYourJungle, nexusKilled, nodeCapture,
                             nodeCaptureAssist, nodeNeutralize, nodeNeutralizeAssist,
                             numDeaths, numItemsBought, objectivePlayerScore, pentaKills,
                             physicalDamageDealtPlayer, physicalDamageDealtToChampions,
                             physicalDamageTaken, quadraKills, sightWardsBought,
                             spellOneCast, spellTwoCast, spellThreeCast, spellFourCast,
                             summonerSpellOneCast, summonerSpellTwoCast,
                             superMonsterKilled, team, teamObjective, timePlayed,
                             totalDamageDealt, totalDamageDealtToChampions,
                             totalDamageTaken, totalHeal, totalPlayerScore,
                             totalScoreRank, totalTimeCrowdControlDealt,
                             totalUnitsHealed, tripleKills, trueDamageDealtPlayer,
                             trueDamageDealtToChampions, trueDamageTaken, turretsKilled,
                             unrealKills, victoryPointTotal, visionWardsBought,
                             wardKilled, wardPlaced)
                        VALUES
                            (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,
                             ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,
                             ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	out := make(chan GameInfo)
	globalWg.Add(1)
	go func() {

		for gi := range gameInfo {
			stats := gi.Game.Statistics
			statRes, statErr := db.Exec(
				statsQuery,
				stats.Assists,
				stats.BarracksKilled,
				stats.ChampionsKilled,
				stats.CombatPlayerScore,
				stats.ConsumablesPurchased,
				stats.DamageDealtPlayer,
				stats.DoubleKills,
				stats.FirstBlood,
				stats.Gold,
				stats.GoldEarned,
				stats.GoldSpent,
				stats.Item0,
				stats.Item1,
				stats.Item2,
				stats.Item3,
				stats.Item4,
				stats.Item5,
				stats.Item6,
				stats.ItemsPurchased,
				stats.KillingSprees,
				stats.LargestCriticalStrike,
				stats.LargestKillingSpree,
				stats.LargestMultiKill,
				stats.LegendaryItemsCreated,
				stats.Level,
				stats.MagicDamageDealtPlayer,
				stats.MagicDamageDealtToChampions,
				stats.MagicDamageTaken,
				stats.MinionsDenied,
				stats.MinionsKilled,
				stats.NeutralMinionsKilled,
				stats.NeutralMinionsKilledEnemyJungle,
				stats.NeutralMinionsKilledYourJungle,
				stats.NexusKilled,
				stats.NodeCapture,
				stats.NodeCaptureAssist,
				stats.NodeNeutralize,
				stats.NodeNeutralizeAssist,
				stats.NumDeaths,
				stats.NumItemsBought,
				stats.ObjectivePlayerScore,
				stats.PentaKills,
				stats.PhysicalDamageDealtPlayer,
				stats.PhysicalDamageDealtToChampions,
				stats.PhysicalDamageTaken,
				stats.QuadraKills,
				stats.SightWardsBought,
				stats.Spell1Cast,
				stats.Spell2Cast,
				stats.Spell3Cast,
				stats.Spell4Cast,
				stats.SummonSpell1Cast,
				stats.SummonSpell2Cast,
				stats.SuperMonsterKilled,
				stats.Team,
				stats.TeamObjective,
				stats.TimePlayed,
				stats.TotalDamageDealt,
				stats.TotalDamageDealtToChampions,
				stats.TotalDamageTaken,
				stats.TotalHeal,
				stats.TotalPlayerScore,
				stats.TotalScoreRank,
				stats.TotalTimeCrowdControlDealt,
				stats.TotalUnitsHealed,
				stats.TripleKills,
				stats.TrueDamageDealtPlayer,
				stats.TrueDamageDealtToChampions,
				stats.TrueDamageTaken,
				stats.TurretsKilled,
				stats.UnrealKills,
				stats.VictoryPointTotal,
				stats.VisionWardsBought,
				stats.WardKilled,
				stats.WardPlaced)
			checkErr(statErr, "Could not insert stats")

			statId, statIdErr := statRes.LastInsertId()
			checkErr(statIdErr, "Could not get last insterted id")

			// drop new versions of GameInfo into a channel
			out <- GameInfo{gi.SummonerID, gi.Game, statId}
		}
		close(out)

		globalWg.Done()
	}()

	return out
}
