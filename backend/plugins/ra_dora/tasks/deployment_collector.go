package tasks

import (
	"encoding/json"
	"log"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Task metadata
var CollectDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "collect_deployments",
	EntryPoint:       CollectDeployment,
	EnabledByDefault: true,
	Description:      "Coleta deployments do banco PostgreSQL do ArgoCD",
}

// Coletor principal
func CollectDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	cfg := config.GetConfig()
	host := cfg.GetString("POSTGRESQL_HOST")
	user := cfg.GetString("POSTGRESQL_USER")
	password := cfg.GetString("POSTGRESQL_PASSWORD")
	dbname := cfg.GetString("POSTGRESQL_DBNAME")

	// Configuração do banco de dados PostgreSQL
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.Default.Wrap(err, "Erro ao conectar ao banco PostgreSQL")
	}

	// Buscar dados da tabela cicd_deployments
	var deployments []models.DatabaseDeployments
	result := db.Find(&deployments)
	if result.Error != nil {
		return errors.Default.Wrap(err, "Erro ao buscar deployments")
	}

	// Obtendo o DAL (Data Access Layer) do DevLakes
	devlakeDb := taskCtx.GetDal()

	// Salvando os deployments no banco _raw do DevLake
	for _, deployment := range deployments {
		data, err := json.Marshal(deployment)
		if err != nil {
			return errors.Default.Wrap(err, "Erro ao serializar dados")
		}

		rawDeployment := models.RawDeployments{RawData: string(data)}

		err = devlakeDb.Create(&rawDeployment)
		if err != nil {
			return errors.Default.Wrap(err, "Erro ao salvar deployment no banco _raw do DevLake")
		}
	}

	log.Println("Coleta de deployments concluída com sucesso!")
	return nil
}
