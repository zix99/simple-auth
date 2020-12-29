// +build testmain

package main

import (
	"os"
	"simple-auth/pkg/config"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(t *testing.T) {
	logrus.Info("Starting main as test...")
	args := slicePostDelim(os.Args[1:])
	simpleAuthServer(config.Load(args...))
}

func slicePostDelim(args []string) []string {
	for i, v := range args {
		if v == "--" {
			return args[i+1:]
		}
	}
	return args
}
