package main

type StarRocksConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	BePort   int `mapstructure:"be_port"`
	Tables   []string
}
