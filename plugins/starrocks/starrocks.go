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

package main

import (
	"github.com/apache/incubator-devlake/plugins/starrocks/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.StarRocks

func main() {
	cmd := &cobra.Command{Use: "StarRocks"}
	sourceType := cmd.Flags().StringP("source_type", "st", "", "Source type")
	sourceDsn := cmd.Flags().StringP("source_dsn", "sd", "", "Source dsn")
	_ = cmd.MarkFlagRequired("host")
	host := cmd.Flags().StringP("host", "h", "", "StarRocks host")
	_ = cmd.MarkFlagRequired("port")
	port := cmd.Flags().StringP("port", "p", "", "StarRocks port")
	_ = cmd.MarkFlagRequired("port")
	beHost := cmd.Flags().StringP("be_host", "BH", "", "StarRocks be host")
	bePort := cmd.Flags().StringP("be_port", "BP", "", "StarRocks be port")
	_ = cmd.MarkFlagRequired("user")
	user := cmd.Flags().StringP("user", "u", "", "StarRocks user")
	_ = cmd.MarkFlagRequired("password")
	password := cmd.Flags().StringP("password", "P", "", "StarRocks password")
	_ = cmd.MarkFlagRequired("database")
	database := cmd.Flags().StringP("database", "d", "", "StarRocks database")
	_ = cmd.MarkFlagRequired("table")
	tables := cmd.Flags().StringArrayP("table", "t", []string{}, "StarRocks table")
	_ = cmd.MarkFlagRequired("batch_size")
	batchSize := cmd.Flags().StringP("batch_size", "b", "", "StarRocks insert batch size")
	_ = cmd.MarkFlagRequired("batch_size")
	extra := cmd.Flags().StringP("extra", "e", "", "StarRocks create table sql extra")
	orderBy := cmd.Flags().StringP("order_by", "o", "", "Source tables order by, default is primary key")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"source_type": sourceType,
			"source_dsn":  sourceDsn,
			"host":        host,
			"port":        port,
			"user":        user,
			"password":    password,
			"database":    database,
			"be_host":     beHost,
			"be_port":     bePort,
			"tables":      tables,
			"batch_size":  batchSize,
			"extra":       extra,
			"order_by":    orderBy,
		})
	}
	runner.RunCmd(cmd)
}
