package drone

import (
	_ "github.com/spf13/pflag"
)

/**
 */
type Options struct {
	Host  string `json:",omitempty" yaml:"host" description:"drone service host address"`
	Token string `json:",omitempty" yaml:"token" description:"drone service token address"`
}

/**
 */
func NewDevopsDroneOptions() *Options {
	return &Options{
		Host:  "",
		Token: "",
	}
}
