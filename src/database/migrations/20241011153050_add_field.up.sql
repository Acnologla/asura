SET
    statement_timeout = 0;

ALTER TABLE
    Item
ADD
    COLUMN extra INT;

CREATE TABLE Transactions(
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userID BIGINT REFERENCES Users(ID) ON DELETE CASCADE,
    authorID BIGINT REFERENCES Users(ID) ON DELETE CASCADE,
    amount INT,
    createdAt BIGINT
);

CREATE INDEX transaction_user ON Transactions (userID);