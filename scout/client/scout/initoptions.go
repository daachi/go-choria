// generated code; DO NOT EDIT

package scoutclient

import (
	"github.com/sirupsen/logrus"
)

type initOptions struct {
	cfgFile  string
	logger   *logrus.Entry
	ns       NodeSource
	progress bool
}

// InitializationOption is an optional setting used to initialize the client
type InitializationOption func(opts *initOptions)

// ConfigFile sets the config file to use, when not set will use the user default
func ConfigFile(f string) InitializationOption {
	return func(o *initOptions) {
		o.cfgFile = f
	}
}

// Logger sets the logger to use else one is made via the Choria framework
func Logger(l *logrus.Entry) InitializationOption {
	return func(o *initOptions) {
		o.logger = l
	}
}

// Discovery sets the NodeSource to use when finding nodes to manage
func Discovery(ns NodeSource) InitializationOption {
	return func(o *initOptions) {
		o.ns = ns
	}
}

// Progress enables displaying a progress bar
func Progress() InitializationOption {
	return func(o *initOptions) {
		o.progress = true
	}
}
