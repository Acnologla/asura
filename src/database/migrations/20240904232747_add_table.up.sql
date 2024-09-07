SET
    statement_timeout = 0;

create Table Trials (
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    win int,
    rooster int
);

CREATE INDEX trial_user ON Trials (userID);