package convertor

import (
	"strconv"
	"strings"
)

func ToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func ToInt(s string) (int, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, err
	}
	return i, nil
}
