package db

import (
	"fmt"

	"github.com/merico-dev/lake/config"
)

func GetConnectionString(dbParams string) string {
	// For now, we only allow override of the params suffix

	user := config.V.GetString("DB_USER")
	pass := config.V.GetString("DB_PASS")
	host := config.V.GetString("DB_HOST")
	port := config.V.GetString("DB_PORT")
	name := config.V.GetString("DB_NAME")
	params := fmt.Sprintf("%v&%v", dbParams, config.V.GetString("DB_PARAMS"))
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", user, pass, host, port, name, params)
}
