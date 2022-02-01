package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

// type Mr struct {
// 	iid int
// }

func TestGitLabMrs(t *testing.T) {
	// var mrs []Mr
	// var mr Mr
	db, err := InitializeGormDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT * FROM gitlab_merge_requests"
	tx := db.Debug().Exec(sqlCommand)
	fmt.Println("KEVIN >>> tx.RowsAffected", tx.RowsAffected)
	fmt.Println("KEVIN >>> tx.Error", tx.Error)
	assert.Equal(t, tx.RowsAffected > 0, true)
	// assert.Nil(t, err)
	// defer rows.Close()

	// colTypes, _ := rows.ColumnTypes()
	// for _, col := range colTypes {

	// 	fmt.Println("KEVIN >>> col Name: ", col.Name())
	// 	fmt.Println("KEVIN >>> col Type: ", col.DatabaseTypeName())
	// }
	// for rows.Next() {
	// 	var mr Mr
	// 	if err := rows.Scan(colTypes); err != nil {
	// 		panic(err)
	// 	}
	// 	mrs = append(mrs, mr)
	// }
	// assert.Equal(t, len(mrs) > 0, true)
}
