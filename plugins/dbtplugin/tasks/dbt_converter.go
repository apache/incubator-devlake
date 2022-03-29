package tasks

import (
	"os/exec"
	"os"
	"github.com/merico-dev/lake/plugins/core"
	
)

func DbtConverter(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*DbtTaskData)
	models := data.Options.SelectedModels
	current_dir, _ := os.Getwd()
	os.Chdir(current_dir+"/plugins/dbtplugin")
	cmd := exec.Command("bash", "-c", "dbt run --select " + models)
	if err := cmd.Run(); err != nil{
		return err
	}
	os.Chdir(current_dir)
	return nil
}

var DbtConverterMeta = core.SubTaskMeta{
	Name:             "DbtConverter",
	EntryPoint:       DbtConverter,
	EnabledByDefault: true,
	Description:      "Convert data by dbt",
}
