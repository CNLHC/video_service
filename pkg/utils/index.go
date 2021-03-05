package utils

import (
	"strconv"
	"strings"
)

type FFMpegStats map[string]string

func (FFMpegStats) Parse(r string) (res FFMpegStats) {
	pairs := strings.Split(r, "\n")
	res = make(FFMpegStats)
	for _, pair := range pairs {
		t := strings.Trim(pair, "\n \t")
		if len(t) > 0 {
			item := strings.Split(t, "=")
			if len(item) == 2 {
				res[item[0]] = item[1]
			}
		}
	}
	return res
}

func (c FFMpegStats) GetOutputMs() (out int) {
	if s, ok := c["out_time_ms"]; ok {
		if i, err := strconv.Atoi(s); err != nil {
			return -1
		} else {
			return i
		}
	}
	return -1
}
