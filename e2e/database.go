/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	mysqlGorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDb() (*sql.DB, error) {
	v := LoadConfigFile()
	dbUrl := v.GetString("DB_URL")
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "mysql" {
		dbUrl = fmt.Sprintf(("%s@tcp(%s)%s?%s"), u.User.String(), u.Host, u.Path, u.RawQuery)
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, err
	}
	fmt.Println("Connected!")
	return db, nil
}

func InitializeGormDb() (*gorm.DB, error) {
	connectionString := "merico:merico@tcp(localhost:3306)/lake"
	db, err := gorm.Open(mysqlGorm.Open(connectionString))
	if err != nil {
		return nil, err
	}
	return db, nil
}
