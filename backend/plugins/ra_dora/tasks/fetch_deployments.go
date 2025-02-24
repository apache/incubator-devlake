package tasks

import (
	"log"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var FetchDeploymentsMeta = plugin.SubTaskMeta{
	Name: "Fetch Deployments",
}

func FetchDeployments() error {
	cfg := config.GetConfig()
	host := cfg.GetString("POSTGRESQL_HOST")
	user := cfg.GetString("POSTGRESQL_USER")
	password := cfg.GetString("POSTGRESQL_PASSWORD")
	dbname := cfg.GetString("POSTGRESQL_DBNAME")

	// Configuração do banco de dados PostgreSQL
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Buscar dados da tabela cicd_deployments
	var deployments []models.CICDDeployment
	result := db.Find(&deployments)
	if result.Error != nil {
		return result.Error
	}

	// Logar as informações
	for _, deployment := range deployments {
		log.Printf("Deployment: %+v\n", deployment)
	}

	return nil
}
