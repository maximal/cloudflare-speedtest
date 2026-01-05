package test

import (
	"time"

	"github.com/maximal/cloudflare-speedtest/internal/duration"
)

type Flags struct {
	BaseUrl         string
	Insecure        bool
	LatencyCount    uint64
	Format          string
	Version         bool
	NoProgress      bool
	SoftLimitString string
	SoftLimit       time.Duration
}

var flags = Flags{
	BaseUrl:         "https://" + SPEED_BASE_URL,
	Insecure:        false,
	LatencyCount:    20,
	Format:          "text",
	Version:         false,
	NoProgress:      false,
	SoftLimitString: "0",
	SoftLimit:       0,
}

func validateFlags() bool {
	switch flags.Format {
	case "text", "json", "jsonl", "influx", "tsv":
		break
	default:
		stderrRed("Error: format must be one of: text, json, jsonl, influx, tsv")
		return false
	}
	softLimit, err := duration.Parse(flags.SoftLimitString)
	if err != nil {
		stderrRed("Error parsing soft limit: %s", err.Error())
		return false
	}
	flags.SoftLimit = softLimit
	return true
}
