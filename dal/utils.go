package dal

import "tpbt/services"

func ShouldInitializeDatabase(prv *services.Provider) bool {
	_, err := prv.DB.Query("SELECT TRUE FROM INITIALIZED_DATABASE")
	if err != nil {
		return true
	}

	return false
}
