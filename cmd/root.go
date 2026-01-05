package cmd

import (
	"github.com/maximal/cloudflare-speedtest/cmd/test"
	"github.com/maximal/cloudflare-speedtest/internal/exit"
)

func Execute() {
	if err := test.TestCmd.Execute(); err != nil {
		exit.Exit(exit.StatusGeneralError, err)
	}
	exit.Exit(exit.StatusOk)
}
