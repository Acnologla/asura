SET
    statement_timeout = 0;

CREATE TABLE Guilds(
    ID BigInt PRIMARY KEY,
    disableLootbox BOOLEAN,
    lootBoxChannel BigInt
)