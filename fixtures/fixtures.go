package fixtures

import (
	"fmt"
	"tpbt/models"
	"tpbt/services"
)

func InitializeDatabase(prv *services.Provider) error {
	tables := []models.DbModel{
		initializedDatabase{},
		group{},
		models.BTStreamer{},
		models.BTPlayer{},
		scoreboard{},
	}

	fmt.Println("\t- Generating the database")

	fmt.Println("\t\t- Dropping the old tables")
	for i := len(tables) - 1; i >= 0; i-- {
		prv.DB.Exec(`DROP TABLE IF EXISTS ` + tables[i].GetTableName())
	}

	fmt.Println("\t\t- Creating the new tables")
	for _, t := range tables {
		_, e := prv.DB.Exec(t.GetCreationScript())
		if e != nil {
			return e
		}
	}

	return InsertFixtures(prv)
}

func InsertFixtures(prv *services.Provider) error {
	groups := map[int]string{
		0: "Bannis",
		1: "Visiteur",
		2: "Streamer",
		3: "Streamer+",
		99: "Admin",
	}

	fmt.Println("\t\t- Inserting groups")
	for k, v := range groups {
		_, err := prv.DB.Exec("INSERT INTO GROUPS(GRP_ID, ROLE) VALUES ($1, $2)", k, v)
		if err != nil {
			return err
		}
	}

	return nil
}