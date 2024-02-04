SET
    statement_timeout = 0;

ALTER TABLE
    Rooster
ADD
    Evolved BOOL;

ALTER TABLE
    Users
ADD
    rankedLastFight BIGINT