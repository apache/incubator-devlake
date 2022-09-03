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

package cmd

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	goerrors "github.com/go-errors/errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile         string
	modifyExistCode bool

	rootCmd = &cobra.Command{
		Use:   `generator [command]`,
		Short: "Apache DevLake Cli Tool -- Code Generator",
	}
)

// Execute executes the root command.
func Execute() errors.Error {
	return errors.Default.WrapRaw(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is PROJECT/.env)")
	rootCmd.PersistentFlags().BoolVar(&modifyExistCode, "modifyExistCode", true, "allow generator modify exist code")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("..")
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !goerrors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && goerrors.As(err, notFound):
		// The config file is optional, we shouldn't exit when the config is not found
	default:
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
