CREATE TABLE IF NOT EXISTS game_info (
    id bigint not null,
    mode varchar(15) not null,
    type varchar(20) not null,
    subType varchar(25) not null,
    mapId int not null,
    date datetime not null,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS league_games (
    gameid bigint not null,
    name varchar(60) null,
    confirmed bit(1) null,
    PRIMARY KEY (gameid),
    FOREIGN KEY (gameid) REFERENCES game_info(id)
);

CREATE TABLE IF NOT EXISTS neta_teams (
    id int(11) auto_increment not null,
    name varchar(60) not null,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS summoners (
    id bigint(20) not null,
    name varchar(40) not null,
    neta_team int(11) null,
    `level` int(11) null,
    last_update datetime null,
    PRIMARY KEY (id),
    FOREIGN KEY (neta_team) REFERENCES neta_teams(id)
);

CREATE TABLE IF NOT EXISTS summoner_stats (
    id int(11) AUTO_INCREMENT not null,
    assists int,
    barracksKilled int,
    championsKilled int,
    combatPlayerScore int,
    consumablesPurchased int,
    damageDealtPlayer int,
    doubleKills int,
    firstBlood int,
    gold int,
    goldEarned int,
    goldSpent int,
    item0 int,
    item1 int,
    item2 int,
    item3 int,
    item4 int,
    item5 int,
    item6 int,
    itemsPurchased int,
    killingSprees int,
    largestCriticalStrike int,
    largestKillingSpree int,
    largestMultiKill int,
    legendaryItemsCreated int,
    level int,
    magicDamageDealtPlayer int,
    magicDamageDealtToChampions int,
    magicDamageTaken int,
    minionsDenied int,
    minionsKilled int,
    neutralMinionsKilled int,
    neutralMinionsKilledEnemyJungle int,
    neutralMinionsKilledYourJungle int,
    nexusKilled bit(1),
    nodeCapture int,
    nodeCaptureAssist int,
    nodeNeutralize int,
    nodeNeutralizeAssist int,
    numDeaths int,
    numItemsBought int,
    objectivePlayerScore int,
    pentaKills int,
    physicalDamageDealtPlayer int,
    physicalDamageDealtToChampions int,
    physicalDamageTaken int,
    quadraKills int,
    sightWardsBought int,
    spellOneCast int,
    spellTwoCast int,
    spellThreeCast int,
    spellFourCast int,
    summonerSpellOneCast int,
    summonerSpellTwoCast int,
    superMonsterKilled int,
    team int,
    teamObjective int,
    timePlayed int,
    totalDamageDealt int,
    totalDamageDealtToChampions int,
    totalDamageTaken int,
    totalHeal int,
    totalPlayerScore int,
    totalScoreRank int,
    totalTimeCrowdControlDealt int,
    totalUnitsHealed int,
    tripleKills int,
    trueDamageDealtPlayer int,
    trueDamageDealtToChampions int,
    trueDamageTaken int,
    turretsKilled int,
    unrealKills int,
    victoryPointTotal int,
    visionWardsBought int,
    wardKilled int,
    wardPlaced int,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS summoner_games (
    id int(11) AUTO_INCREMENT not null,
    summonerId bigint not null,
    gameId bigint not null,
    championId int not null,
    spellOne int not null,
    spellTwo int not null,
    statId int(11) not null,
    won bit(1),
    PRIMARY KEY (id),
    FOREIGN KEY (statId) REFERENCES summoner_stats(id),
    FOREIGN KEY (summonerId) REFERENCES summoners(id),
    FOREIGN KEY (gameId) REFERENCES game_info(id)
);

CREATE TABLE IF NOT EXISTS runtimes (
    id int(11) AUTO_INCREMENT not null,
    startTime datetime not null,
    endTime datetime,
    records int(11),
    PRIMARY KEY (id)
);

-- Pre-insert summoners
INSERT INTO `summoners`
    (name, id, `level`)
VALUES
     ('Darkmist16',77804,30)
    ,('SuicideSnowman',21305835,30)
    ,('Misaga',21465652,30)
    ,('Delath',134961,30)
    ,('AzayakaAkari',19772280,30)
    ,('riskman64',24199871,30)
    ,('Khoza',,)
    ,('echoblaze',,)
    ,('TheJuggler',,)
    ,('m1tsu',,)
    ,('019Ky',,)
    ,('Mystenance',,)
    ,('titan alibaba',,)
    ,('striderfox',,)
    ,('Faucetin',,)
    ,('Takeiteasyonme',,)
    ,('HybridEleven',,)
    ,('taaaakun',,)
    ,('Psychotic Idiot',,)
    ,('W4yl4nder',,)
    ,('Ngsanity',,)
    ,('desunman',,)
    ,('fadedlightx',,)
    ,('akiraK',,)
    ,('stopisme',,)
    ,('ginourmous',,)
    ,('karenkun',,)
    ,('kalun85',,)
    ,('juebag',,)
    ,('1337bagger',,);
