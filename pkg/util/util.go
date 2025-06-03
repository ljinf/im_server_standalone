package util

import "strconv"

func StringToNum(str string) int64 {
	if len(str) < 1 {
		return 0
	}

	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}
