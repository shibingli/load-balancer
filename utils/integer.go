package utils

import (
	"strconv"
	"strings"
)

// SearchInts similar to sort.SearchInts(), but faster.
func SearchInts(a []int, x int) (i int) {
	j := len(a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// MustInt converted to int, default is 0
func MustInt(v interface{}) int {
	switch i := v.(type) {
	default:
		d, ok := i.(int)
		if ok {
			return d
		}
		return 0
	case string:
		v, err := strconv.Atoi(strings.TrimSpace(i))
		if err == nil {
			return v
		}
		return 0
	case bool:
		if i {
			return 1
		}
		return 0
	case nil:
		return 0
	case int:
		return i
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	case uint:
		return int(i)
	case uint8:
		return int(i)
	case uint16:
		return int(i)
	case uint32:
		return int(i)
	case uint64:
		return int(i)
	case float32:
		return int(i)
	case float64:
		return int(i)
	}
}
