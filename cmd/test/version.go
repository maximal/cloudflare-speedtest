package test

import (
	"regexp"
	"runtime/debug"

	"github.com/maximal/cloudflare-speedtest/internal/exit"
)

var versionRegex = regexp.MustCompile(`v\d+\.\d+\.\d+\-\d{14}\-([a-f0-9]{8})[a-f0-9]*(\+dirty)?`)

func showVersion() exit.Status {
	var commit string
	if info, ok := debug.ReadBuildInfo(); ok {
		match := versionRegex.FindStringSubmatch(info.Main.Version)
		if len(match) == 3 {
			commit = match[1] + match[2]
		} else {
			commit = "dev"
		}
	} else {
		commit = "dev"
	}
	print(PROGRAM_TITLE + " v" + VERSION + "+" + commit)
	return exit.StatusOk
}
