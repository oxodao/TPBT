package fixtures

type group struct{}

func (g group) GetTableName() string {
	return "GROUPS"
}

func (g group) GetCreationScript() string {
	return ` CREATE TABLE GROUPS (
		GRP_ID SERIAL PRIMARY KEY,
		ROLE VARCHAR(20)
	);`
}

type scoreboard struct {}

func (sc scoreboard) GetTableName() string {
	return "SCOREBOARD"
}

func (sc scoreboard) GetCreationScript() string {
	return ` CREATE TABLE SCOREBOARD (
		  VIEWER_ID TEXT,
		  STREAMER_ID TEXT,
		  SCORE INT,
		  FOREIGN KEY (VIEWER_ID) REFERENCES VIEWER(TW_ID),
		  FOREIGN KEY (STREAMER_ID) REFERENCES TWITCH_USERS(TW_ID),
		  PRIMARY KEY (VIEWER_ID, STREAMER_ID)
	);`
}




type initializedDatabase struct {}

func (id initializedDatabase) GetTableName() string {
	return "INITIALIZED_DATABASE"
}

func (id initializedDatabase) GetCreationScript() string {
	return `CREATE TABLE INITIALIZED_DATABASE (ID INTEGER);`
}