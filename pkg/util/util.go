package util

import (
	"strings"
	"time"
)

func GetPortFromHost(host string) string {
	ss := strings.Split(host, ":")
	if len(ss) > 1 {
		return ss[1]
	}

	return "80"
}

func ConvertSecondsToDuration(s int) *time.Duration {
	if s == 0 {
		s = 5
	}

	d := time.Duration(s) * time.Second
	return &d
}
