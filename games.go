package main

import (
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/coopernurse/gorp"
)

// use list of summoners to download game data
func updateGames(summoners []int64, dbmap *gorp.DbMap) (gameCount int) {

	// get stored game ids for summoners
	var gameIdQuery string = `SELECT id FROM  game_info`
	var gameIds []int64 = make([]int64, 1)
	_, err := dbmap.Select(&gameIds, gameIdQuery)
	checkErr(err, "Could not get game ids from database")

	// get most recent games for each summoner
	var savedGames = make([]int64, 1)
	for _, summonerId := range summoners {
		summonerGames, riotErr := goriot.RecentGameBySummoner(goriot.NA, summonerId)
		if riotErr != nil {
			panic(riotErr)
		} else {

			// save game if we don't already have it
			for _, game := range summonerGames {

				if existsInSlice(game.GameID, gameIds) == false {

					if existsInSlice(game.GameID, savedGames) == false {

						fmt.Println(game.CreateDate)
						// save game to db
						var gameInfoQuery string = `
                            INSERT INTO game_info
                                (id, mode, type, subType, mapId, date)
                            VALUES
                                (?, ?, ?, ?, ?, ?)`
						_, infoErr := dbmap.Exec(
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

						// save game id to list to skip loading into
						// game_info, in the event it shows again
						savedGames = append(savedGames, game.GameID)
					}
					var stats goriot.GameStat = game.Statistics
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
					statRes, statErr := dbmap.Exec(
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

					// save summoner_game
					var summonerGameQuery string = `
                        INSERT INTO summoner_games
                            (summonerId, gameId, championId, spellOne, spellTwo,
                             statId, won)
                        VALUES
                            (?,?,?,?,?,?,?)`
					_, sgErr := dbmap.Exec(
						summonerGameQuery,
						summonerId,
						game.GameID,
						game.ChampionID,
						game.Spell1,
						game.Spell2,
						statId,
						stats.Win)
					checkErr(sgErr, "Could not save summoner game info")
				}
			}

		}
	}

	return len(savedGames)
}
