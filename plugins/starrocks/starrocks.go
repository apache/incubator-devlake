package main

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

type StarRocks string

func (s StarRocks) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		LoadDataTaskMeta,
	}
}

func (s StarRocks) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op StarRocksConfig
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (s StarRocks) Description() string {
	return "Sync data from database to StarRocks"
}

func (s StarRocks) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/starrocks"
}

var PluginEntry StarRocks

func main() {
	cmd := &cobra.Command{Use: "StarRocks"}
	_ = cmd.MarkFlagRequired("host")
	host := cmd.Flags().StringP("host", "h", "", "StarRocks host")
	_ = cmd.MarkFlagRequired("port")
	port := cmd.Flags().StringP("port", "p", "", "StarRocks port")
	_ = cmd.MarkFlagRequired("port")
	bePort := cmd.Flags().StringP("be_port", "BP", "", "StarRocks be port")
	_ = cmd.MarkFlagRequired("user")
	user := cmd.Flags().StringP("user", "u", "", "StarRocks user")
	_ = cmd.MarkFlagRequired("password")
	password := cmd.Flags().StringP("password", "P", "", "StarRocks password")
	_ = cmd.MarkFlagRequired("database")
	database := cmd.Flags().StringP("database", "d", "", "StarRocks database")
	_ = cmd.MarkFlagRequired("table")
	tables := cmd.Flags().StringArrayP("table", "t", []string{}, "StarRocks table")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"host":     host,
			"port":     port,
			"user":     user,
			"password": password,
			"database": database,
			"be_port":  bePort,
			"tables":   tables,
		})
	}
	runner.RunCmd(cmd)
}
