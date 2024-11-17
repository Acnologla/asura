SET
    statement_timeout = 0;

create Table Tower (
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users (ID) ON DELETE CASCADE NOT NULL,
    floor int,
    lose bool
);

CREATE INDEX tower_user ON Tower (userID);