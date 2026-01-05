package duration

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	secondsRegex        = regexp.MustCompile(`^\s*(\d+)\s*[sS]?\s*$`)
	minutesSecondsRegex = regexp.MustCompile(`^\s*(\d+)\s*[mM]\s*(\d+)\s*[sS]\s*$`)
)

func Parse(string string) (time.Duration, error) {
	match := secondsRegex.FindStringSubmatch(string)
	if len(match) == 2 {
		sec, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(sec * int64(time.Second)), nil
	}

	match = minutesSecondsRegex.FindStringSubmatch(string)
	if len(match) == 3 {
		min, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		sec, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(min*int64(time.Minute) + sec*int64(time.Second)), nil
	}

	return 0, errors.New("string does not match any of duration patterns: `<n>s`, `<n>m<n>s`")
}
