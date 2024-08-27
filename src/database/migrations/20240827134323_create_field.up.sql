SET
    statement_timeout = 0;

ALTER TABLE
    Users
ADD
    lastRank int;

ALTER TABLE
    Users
ADD
    rankedLastFight BIGINT