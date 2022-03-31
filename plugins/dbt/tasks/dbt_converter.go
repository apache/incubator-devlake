package tasks

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
)

func DbtConverter(taskCtx core.SubTaskContext) error {

	taskCtx.SetProgress(0, -1)
	data := taskCtx.GetData().(*DbtTaskData)
	models := data.Options.SelectedModels
	projectPath := data.Options.ProjectPath
	projectName := data.Options.ProjectName
	projectTarget := data.Options.ProjectTarget

	dbUrl := taskCtx.GetConfig("DB_URL")
	dbSlice := strings.FieldsFunc(dbUrl, func(r rune) bool { return strings.ContainsRune(":@()/?", r) })
	dbUsername := dbSlice[0]
	dbPassword := dbSlice[1]
	dbServer := dbSlice[3]
	dbPort, _ := strconv.Atoi(dbSlice[4])
	dbSchema := dbSlice[5]

	os.Chdir(projectPath)
	config := viper.New()
	config.Set(projectName+".target", projectTarget)
	config.Set(projectName+".outputs."+projectTarget+".type", "mysql")
	config.Set(projectName+".outputs."+projectTarget+".server", dbServer)
	config.Set(projectName+".outputs."+projectTarget+".port", dbPort)
	config.Set(projectName+".outputs."+projectTarget+".schema", dbSchema)
	config.Set(projectName+".outputs."+projectTarget+".database", dbSchema)
	config.Set(projectName+".outputs."+projectTarget+".username", dbUsername)
	config.Set(projectName+".outputs."+projectTarget+".password", dbPassword)
	config.WriteConfigAs("profiles.yml")

	dbtExecParams := []string{"dbt", "run", "--profiles-dir", projectPath, "--select"}
	dbtExecParams = append(dbtExecParams, models...)
	cmd := exec.Command(dbtExecParams[0], (dbtExecParams[1:])...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {

		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		substr1 := "of"
		substr2 := "OK"
		if strings.Contains(line, substr1) && strings.Contains(line, substr2) {
			taskCtx.IncProgress(1)
		}
	}

	return nil
}

var DbtConverterMeta = core.SubTaskMeta{
	Name:             "DbtConverter",
	EntryPoint:       DbtConverter,
	EnabledByDefault: true,
	Description:      "Convert data by dbt",
}
