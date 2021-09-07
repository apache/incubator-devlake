package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newApiCommand() *cobra.Command {
	options := &apiOptions{}
	cmd := &cobra.Command{
		Use:   "api {endpoint}",
		Short: "Create api request",
		Long:  options.Manuals(),
		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Validate(cmd, args); err != nil {
				log.Fatalln(err)
			}
			if err := options.Run(cmd, args); err != nil {
				log.Fatalln(err)
			}
		},
	}
	options.Install(cmd.Flags())
	return cmd
}

type apiOptions struct {
	Host   string
	Method string
	Body   string
	Cron   string
}

func (o *apiOptions) Install(flags *pflag.FlagSet) {
	flags.StringVarP(&o.Host, "host", "H", "http://127.0.0.1:8080", "lake server host")
	flags.StringVarP(&o.Method, "method", "m", "GET", "request method")
	flags.StringVar(&o.Body, "body", "", "request body")
	flags.StringVar(&o.Cron, "cron", "", "create cron job for api request")
}

func (o *apiOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("endpoint required, using like this: lake-cli api {endpoint}")
	}
	var supportedEndpoints = make(map[string]bool)
	supportedEndpoints["task"] = true
	var endpoint = args[0]
	if support, ok := supportedEndpoints[endpoint]; !support || !ok {
		return fmt.Errorf("unsupported endpoint %s", endpoint)
	}
	return nil
}

// ./lake-cli api task -m POST --body "{'plugin':'jira', 'options': {'boardId': 8}}"
func (o *apiOptions) Run(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	fmt.Println(*o)
	if strings.TrimSpace(o.Cron) != "" {
		// TODO: create cron job

		return nil
	}
	err := DoRequest(fmt.Sprintf("%s%s", o.Host, args[0]), o.Body)
	return err
}

func (o *apiOptions) Manuals() string {
	return ""
}

func DoRequest(url string, body string) error {
	// TODO: do request
	return nil
}
