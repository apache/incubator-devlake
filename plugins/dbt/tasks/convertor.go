package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
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

	dbType := "mysql"
	var dbUsername string
	var dbPassword string
	var dbServer string
	var dbPort string
	var dbDataBase string
	var dbSchema string

	dbUrl := taskCtx.GetConfig("DB_URL")
	flag := strings.Contains(dbUrl, "://")
	if flag {
		// other database
		u, err := url.Parse(dbUrl)
		if err != nil {
			return err
		}
		dbType = u.Scheme
		dbUsername = u.User.Username()
		dbPassword, _ = u.User.Password()
		dbServer, dbPort, _ = net.SplitHostPort(u.Host)
		dbDataBase = u.Path[1:]
		mapQuery, _ := url.ParseQuery(u.RawQuery)
		if value, ok := mapQuery["search_path"]; ok{
			dbSchema = value[0]
		}else{
			dbSchema = "public"
		}

	} else {
		// mysql database
		dbSlice := strings.FieldsFunc(dbUrl, func(r rune) bool { return strings.ContainsRune(":@()/?", r) })
		if len(dbSlice) < 6 {
			return fmt.Errorf("DB_URL data parsing error, please check the DB_URL value in .env file")
		}
		dbUsername = dbSlice[0]
		dbPassword = dbSlice[1]
		dbServer = dbSlice[3]
		dbPort = dbSlice[4]
		dbDataBase = dbSlice[5]
		dbSchema = dbDataBase
	}

	err := os.Chdir(projectPath)
	if err != nil {
		return err
	}
	config := viper.New()
	config.Set(projectName+".target", projectTarget)
	config.Set(projectName+".outputs."+projectTarget+".type", dbType)
	
	dbPortInt, _ := strconv.Atoi(dbPort)
	config.Set(projectName+".outputs."+projectTarget+".port", dbPortInt)
	config.Set(projectName+".outputs."+projectTarget+".password", dbPassword)
	config.Set(projectName+".outputs."+projectTarget+".schema", dbSchema)
	if flag {
		config.Set(projectName+".outputs."+projectTarget+".host", dbServer)
		config.Set(projectName+".outputs."+projectTarget+".user", dbUsername)
		config.Set(projectName+".outputs."+projectTarget+".dbname", dbDataBase)
	} else {
		config.Set(projectName+".outputs."+projectTarget+".server", dbServer)
		config.Set(projectName+".outputs."+projectTarget+".username", dbUsername)
		config.Set(projectName+".outputs."+projectTarget+".database", dbDataBase)	
	}
		
	err = config.WriteConfigAs("profiles.yml")
	if err != nil {
		return err
	}

	dbtExecParams := []string{"dbt", "run", "--profiles-dir", projectPath}
	if projectVars != nil {
		jsonProjectVars, err := json.Marshal(projectVars)
		if err != nil {
			return fmt.Errorf("parameters vars json marshal error")
		}
		dbtExecParams = append(dbtExecParams, "--vars")
		dbtExecParams = append(dbtExecParams, string(jsonProjectVars))
	}
	dbtExecParams = append(dbtExecParams, "--select")
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
