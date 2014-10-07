package main

import (
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
)

type PlayedGame struct {
	SummonerId int64 `db:"summonerId"`
	GameId     int64 `db:"gameId"`
}

// use list of summoners to download game data
// TODO: use a channel as input for summoner IDs
// TODO: use a channel to ouput summoner IDs as they finish
func updateGames(summoners []int64, dbmap *gorp.DbMap) (gameCount int) {

	// get stored game ids for summoners
	var summonerGameIdQuery string = `
        SELECT gameId
        FROM summoner_games
        INNER JOIN game_info gi ON
            gi.id = gameId
        WHERE summonerId = %d
        ORDER by gi.date desc`
	var gameIdQuery string = "SELECT id FROM game_info"
	var gameIds []int64 = make([]int64, 1)
	_, err := dbmap.Select(&gameIds, gameIdQuery)
	checkErr(err, "Could not get game ids from database")

	// get most recent games for each summoner
	var playedGames []int64
	var savedGames = make([]int64, 1)
	for _, summonerId := range summoners {

		playedGames = make([]int64, 1)
		summonerGames, riotErr := goriot.RecentGameBySummoner(goriot.NA, summonerId)
		if riotErr != nil {
			panic(riotErr)
		}

		_, err := dbmap.Select(
			&playedGames,
			fmt.Sprintf(summonerGameIdQuery, summonerId))
		checkErr(err, "Could not find games for summoner")

		// save game if we don't already have it
		for _, game := range summonerGames {

			//TODO: get the last played game ID and compare to
			//      downloaded games

			updateGameInfo(game, dbmap)

			var stats goriot.GameStat = game.Statistics
			statId := updateSummonerStatistics(summonerId, stats, dbmap)

			updateSummonerGames(
				summonerId, game, statId, stats.Win, dbmap)
		}

	}

	return len(savedGames)
}

func updateGameInfo(game goriot.Game, db *gorp.DbMap) {
	// save game to db
	// TODO: Use insert/update query instead
	var gameInfoQuery string = `
                            INSERT INTO game_info
                                (id, mode, type, subType, mapId, date)
                            VALUES
                                (?, ?, ?, ?, ?, FROM_UNIXTIME(?))`
	_, infoErr := db.Exec(
		gameInfoQuery,
		game.GameID,
		game.GameMode,
		game.GameType,
		game.SubType,
		game.MapID,
		game.CreateDate/1000)
	if infoErr != nil {
		panic(infoErr)
	}
}

func updateSummonerGames(
	summonerId int64,
	game goriot.Game,
	statId int64,
	wonGame bool,
	db *gorp.DbMap) {
	// save summoner_game
	// TODO: Use insert/update query instead
	var summonerGameQuery string = `
                        INSERT INTO summoner_games
                            (summonerId, gameId, championId, spellOne, spellTwo,
                             statId, won)
                        VALUES
                            (?,?,?,?,?,?,?)`
	_, sgErr := db.Exec(
		summonerGameQuery,
		summonerId,
		game.GameID,
		game.ChampionID,
		game.Spell1,
		game.Spell2,
		statId,
		wonGame)
	checkErr(sgErr, "Could not save summoner game info")
}

func updateSummonerStatistics(
	summonerId int64,
	stats goriot.GameStat,
	db *gorp.DbMap) (statId int64) {
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

	return statId
}
