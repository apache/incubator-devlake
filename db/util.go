package db

import (
	"fmt"

	"github.com/merico-dev/lake/config"
)

func GetConnectionString(dbParams string, includeDriver bool) string {
	user := config.V.GetString("DB_USER")
	pass := config.V.GetString("DB_PASSWORD")
	host := config.V.GetString("DB_HOST")
	port := config.V.GetString("DB_PORT")
	name := config.V.GetString("DB_DATABASE")
	driver := config.V.GetString("DB_DRIVER")

	params := fmt.Sprintf("%v&%v", dbParams, config.V.GetString("DB_PARAMS"))
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", user, pass, host, port, name, params)

	if includeDriver {
		return fmt.Sprintf("%v://%v", driver, connectionString)
	} else {
		return connectionString
	}
}
