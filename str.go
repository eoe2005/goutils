package goutils

import "strings"

func StrGetBody(src, start, end string) (int, string) {
	si := strings.Index(src, start)
	ei := strings.Index(src[si+len(start):], end)
	if si < 0 || ei < 0 {
		return -1, ""
	}
	ei = si + len(start) + ei
	return ei + len(end), src[si+len(start) : ei]
}
