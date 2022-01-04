package db

import (
	"fmt"
	"log"
	"net/url"

	"github.com/merico-dev/lake/config"
)

func GetConnectionString(dbParams map[string]string, includeDriver bool) string {
	u, err := url.Parse(config.V.GetString("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	for k, v := range dbParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	fmt.Println("JON >>> Connection String", u)

	if includeDriver {
		return fmt.Sprintf("mysql://%v", u.String())
	} else {
		return u.String()
	}
}
