package human_size

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

var reg *regexp.Regexp

func init() {
	reg = regexp.MustCompile("^(\\d+) *(([kKMGT]?)i?B)?$")
}

func Parse(size string) (int, bool) {
	matches := reg.FindStringSubmatch(size)
	if len(matches) != 4 {
		return 0, false
	}
	cnt, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, false
	}
	switch matches[3] {
	case "":
		return cnt, true
	case "K", "k":
		return cnt * KB, true
	case "M":
		return cnt * MB, true
	case "G":
		return cnt * GB, true
	case "T":
		return cnt * TB, true
	}
	return 0, false
}

func Format(cnt int) string {
	if cnt < KB {
		return fmt.Sprintf("%d B", cnt)
	}
	if cnt < MB {
		return fmt.Sprintf("%.2f kB", float32(cnt)/float32(KB))
	}
	if cnt < GB {
		return fmt.Sprintf("%.2f MB", float32(cnt)/float32(MB))
	}
	if cnt < TB {
		return fmt.Sprintf("%.2f GB", float32(cnt)/float32(GB))
	}
	return fmt.Sprintf("%.2f TB", float32(cnt)/float32(TB))
}
