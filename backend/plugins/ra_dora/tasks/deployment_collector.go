package tasks

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	RAW_DEPLOYMENT = "ra_deployment"
)

type RaDoraTaskData struct {
	Deployments *models.CICDDeployment
}

// Task metadata
var CollectDeploymentsMeta = plugin.SubTaskMeta{
	Name: "collect_deployments",
	EntryPoint: func(taskCtx plugin.SubTaskContext) errors.Error {
		return errors.Convert(CollectDeployments(taskCtx))
	},
	EnabledByDefault: true,
	Description:      "Coleta deployments do ArgoCD via banco de dados",
}

// Coletor principal
func CollectDeployments(taskCtx plugin.SubTaskContext) error {
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

	// Inserir dados na tabela _raw
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT)
	for _, deployment := range deployments {
		err := InsertRawData(rawDataSubTaskArgs, deployment)
		if err != nil {
			return err
		}
	}

	log.Println("Coleta de deployments concluída com sucesso!")
	return nil
}

func CreateRawDataSubTaskArgs(subtaskCtx plugin.SubTaskContext, Table string) (*api.RawDataSubTaskArgs, *RaDoraTaskData) {
	data := subtaskCtx.GetData().(*RaDoraTaskData)
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:   subtaskCtx,
		Table: Table,
	}
	return rawDataSubTaskArgs, data
}

func InsertRawData(args *api.RawDataSubTaskArgs, deployment models.CICDDeployment) error {
	value, _ := strconv.ParseUint(deployment.ID, 10, 64)

	rawData := &api.RawData{
		ID: value,
		Data: func() []byte {
			data, err := json.Marshal(deployment)
			if err != nil {
				log.Fatalf("failed to marshal deployment: %v", err)
			}
			return data
		}(),
	}
	return args.Ctx.GetDal().Create(rawData)
}
