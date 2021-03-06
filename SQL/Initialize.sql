DROP TABLE SCOREBOARD;
DROP TABLE VIEWER;
DROP TABLE TWITCH_USERS;
DROP TABLE GROUPS;

CREATE TABLE GROUPS (
    GRP_ID SERIAL PRIMARY KEY,
    ROLE VARCHAR(20)
);

CREATE TABLE TWITCH_USERS (
  TW_ID          TEXT PRIMARY KEY,
  TW_NAME        TEXT,
  APP_TOKEN      TEXT,
  ACCESS_TOKEN   TEXT,
  REFRESH_TOKEN  TEXT,
  BOT_IN_CHANNEL BOOLEAN NOT NULL DEFAULT false,
  GRP_ID         INTEGER CONSTRAINT fk_usr_grp REFERENCES GROUPS(GRP_ID)
);

CREATE TABLE VIEWER (
    TW_ID TEXT PRIMARY KEY,
    LAST_KNOWN_USERNAME TEXT
);

CREATE TABLE SCOREBOARD (
  VIEWER_ID TEXT,
  STREAMER_ID TEXT,
  SCORE INT,
  FOREIGN KEY (VIEWER_ID) REFERENCES VIEWER(TW_ID),
  FOREIGN KEY (STREAMER_ID) REFERENCES TWITCH_USERS(TW_ID),
  PRIMARY KEY (VIEWER_ID, STREAMER_ID)
);

INSERT INTO GROUPS (GRP_ID, ROLE) VALUES (0, 'Bannis');
INSERT INTO GROUPS (GRP_ID, ROLE) VALUES (1, 'Visiteur');
INSERT INTO GROUPS (GRP_ID, ROLE) VALUES (2, 'Streamer');
INSERT INTO GROUPS (GRP_ID, ROLE) VALUES (3, 'Streamer+');
INSERT INTO GROUPS (GRP_ID, ROLE) VALUES (99, 'Admin');

