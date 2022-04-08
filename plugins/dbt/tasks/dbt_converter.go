package tasks

import (
	"bufio"
	"fmt"
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
	projectVars := data.Options.ProjectVars

	dbUrl := taskCtx.GetConfig("DB_URL")
	dbSlice := strings.FieldsFunc(dbUrl, func(r rune) bool { return strings.ContainsRune(":@()/?", r) })
	if len(dbSlice) < 6 {
		return fmt.Errorf("DB_URL parsing error")
	}

	dbUsername := dbSlice[0]
	dbPassword := dbSlice[1]
	dbServer := dbSlice[3]
	dbPort, _ := strconv.Atoi(dbSlice[4])
	dbSchema := dbSlice[5]

	err := os.Chdir(projectPath)
	if err != nil {
		return err
	}
	config := viper.New()
	config.Set(projectName+".target", projectTarget)
	config.Set(projectName+".outputs."+projectTarget+".type", "mysql")
	config.Set(projectName+".outputs."+projectTarget+".server", dbServer)
	config.Set(projectName+".outputs."+projectTarget+".port", dbPort)
	config.Set(projectName+".outputs."+projectTarget+".schema", dbSchema)
	config.Set(projectName+".outputs."+projectTarget+".database", dbSchema)
	config.Set(projectName+".outputs."+projectTarget+".username", dbUsername)
	config.Set(projectName+".outputs."+projectTarget+".password", dbPassword)
	err = config.WriteConfigAs("profiles.yml")
	if err != nil {
		return err
	}

	dbtExecParams := []string{"dbt", "run", "--profiles-dir", projectPath, "--select"}
	if projectVars != "" {
		dbtExecParams = []string{"dbt", "run", "--profiles-dir", projectPath, "--vars", projectVars, "--select"}
	}
	dbtExecParams = append(dbtExecParams, models...)
	cmd := exec.Command(dbtExecParams[0], (dbtExecParams[1:])...)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "of") && strings.Contains(line, "OK") {
			taskCtx.IncProgress(1)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var DbtConverterMeta = core.SubTaskMeta{
	Name:             "DbtConverter",
	EntryPoint:       DbtConverter,
	EnabledByDefault: true,
	Description:      "Convert data by dbt",
}
