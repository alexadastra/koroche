package config

import "flag"

// TODO: figure out how to pass args more easy

type Args struct {
	ServiceName string
	ConfigPath  string
}

func ParseFlags() *Args {
	var a Args
	flag.StringVar(&a.ServiceName, "name", "", "defines service name")
	flag.StringVar(&a.ConfigPath, "config_path", "/etc/config/config.yaml",
		"defines the path where the service reads config from")
	flag.Parse()
	return &a
}
