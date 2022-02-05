SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Clan(
    name VARCHAR(25) PRIMARY KEY,
    xp INT,
    createdAt BigInt,
    background VARCHAR(300),
    money INT,
    membersUpgrade INT,
    banksUpgrade INT,
    missionsUpgrade INT,
    mission BigInt,
    missionProgress INT
);

CREATE TABLE Users(
    ID BigInt PRIMARY KEY,
    xp INT,
    upgrades INT [],
    win int,
    lose int,
    money INT,
    clan VARCHAR(26) REFERENCES Clan(name),
    dungeon INT,
    dungeonReset INT,
    tradeMission BigInt,
    lastMission BigInt,
    vip BigInt,
    vipBackground varchar(200),
    trainLimit INT,
    asuraCoin INT,
    arenaActive BOOLEAN,
    arenaWin INT,
    arenaLose INT,
    arenaLastFight BIGINT,
    rank INT,
    tradeItem BigInt,
    daily BigInt,
    dailyStrikes INT,
    pity INT
);

CREATE TABLE Mission(
    userId BIGINT REFERENCES Users(ID) ON DELETE CASCADE PRIMARY KEY,
    type INT, 
    level INT,
    progress INT,
    adv int
);

CREATE TABLE Rooster(
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    name VARCHAR(26),
    resets INT,
	equip BOOL,
    xp INT,
    type INT,
    equipped INT []
);



CREATE TABLE ClanMember(
    ID BigInt PRIMARY KEY,
    clan VARCHAR(26) REFERENCES Clan(name) ON DELETE CASCADE NOT NULL,
    role INT,
    xp INT
);

CREATE TABLE Item(
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    quantity INT,
    itemID INT,
	equip BOOL,
    type INT
);

