package utils

import "strings"

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
