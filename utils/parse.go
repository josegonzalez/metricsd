package utils

import "strconv"

func ParseInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}
