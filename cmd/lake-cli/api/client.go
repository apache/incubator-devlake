package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/merico-dev/lake/logger"
	"github.com/robfig/cron/v3"
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

// ./lake-cli api task -m POST --body "[{\"plugin\":\"jira\", \"options\": {\"boardId\": 8}}]" --cron "@every 5s"
func (o *apiOptions) Run(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	fmt.Println(*o)

	sleep := make(chan bool)

	if strings.TrimSpace(o.Cron) != "" {
		c := cron.New()
		_, err := c.AddFunc(o.Cron, func() {
			err := DoRequest(fmt.Sprintf("%s/%s", o.Host, args[0]), o.Method, o.Body)
			if err != nil {
				logger.Error("failed to do request", err)
			}
		})
		if err != nil {
			return err
		}
		c.Start()
		<-sleep
		return nil
	}
	err := DoRequest(fmt.Sprintf("%s/%s", o.Host, args[0]), o.Method, o.Body)
	return err
}

func (o *apiOptions) Manuals() string {
	return ""
}

func DoRequest(url, method, body string) error {
	switch strings.ToUpper(method) {
	case "POST":
		return Post(url, body)
	}
	return nil
}

func Post(url, body string) error {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(responseBody))
	return nil
}
