package api

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	supportedEndpoints["pipeline"] = true
	var endpoint = args[0]
	if support, ok := supportedEndpoints[endpoint]; !support || !ok {
		return fmt.Errorf("unsupported endpoint %s", endpoint)
	}
	return nil
}

func (o *apiOptions) Run(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(o.Body) != "" {
		o.Body = readBodyFromFile(o.Body)
	}
	fmt.Printf("%+v\n", o)
	err := DoRequest(fmt.Sprintf("%s/%s", o.Host, args[0]), o.Method, o.Body)
	if err != nil {
		return err
	}
	if strings.TrimSpace(o.Cron) != "" {
		// create cron job
		sleep := make(chan bool)
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
	}
	return nil
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
	log.Printf("POST to %v\n", url)
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

func readBodyFromFile(input string) string {
	// TODO: check if input is path like

	// read file if exist, otherwise return input string
	if _, err := os.Stat(input); err != nil {
		log.Println(err)
		return input
	}
	f, err := os.Open(input)
	defer close(f)
	if err != nil {
		log.Fatalln(err)
		return input
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
		return input
	}
	return string(content)
}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
