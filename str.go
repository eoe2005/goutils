package goutils

import "strings"

func StrGetBody(src, start, end string) (int, string) {
	si := strings.Index(src, start)
	ei := strings.Index(src[si+len(start):], end)
	if si < 0 || ei < 0 || si >= ei {
		return -1, ""
	}
	return ei + len(end), src[si+len(start) : ei]
}
