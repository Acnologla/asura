SET statement_timeout = 0;
CREATE INDEX rooster_user ON Rooster (userID);
CREATE INDEX item_user ON Item (userID);
CREATE INDEX mission_user on Mission(userID);
CREATE INDEX clanMember_clan ON ClanMember (clan);