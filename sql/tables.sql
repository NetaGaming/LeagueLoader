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

CREATE TABLE IF NOT EXISTS summoner_games (
    summonerId bigint not null,
    gameId bigint not null,
    championId int not null,
    spellOne int not null,
    spellTwo int not null,
    minionsKilled int,
    numDeaths int,
    assists int,
    championsKilled int,
    won bit(1),
    PRIMARY KEY (summonerId, gameId),
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
    ,('Khoza',72009,30)
    ,('echoblaze',21519850,30)
    ,('TheJuggler',82712,30)
    ,('m1tsu',31460782,30)
    ,('019Ky',30082388,30)
    ,('Mystenance',21731401,30)
    ,('titan alibaba',52354738,30)
    ,('striderfox',24461461,30)
    ,('Faucetin',50981870,30)
    ,('Takeiteasyonme',24200266,30)
    ,('HybridEleven',43238261,30)
    ,('taaaakun',24469600,30)
    ,('Psychotic Idiot',41764137,30)
    ,('W4yl4nder',58882619,30)
    ,('Ngsanity',24220006,30)
    ,('desunman',37051067,30)
    ,('fadedlightx',43538849,30)
    ,('akiraK',19928736,30)
    ,('stopisme',20428822,30)
    ,('ginourmous',24221257,30)
    ,('karenkun',37058280,30)
    ,('kalun85',24404121,30)
    ,('juebag',21867311,30)
    ,('1337bagger',50063257,30);
