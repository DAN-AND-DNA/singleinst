package cmd

import (
	flag "github.com/spf13/pflag"
)

var (
	ConfigPath string
)

func init() {
	flag.StringVarP(&ConfigPath, "config", "c", ".", "配置目录位置")
	flag.Parse()
}
