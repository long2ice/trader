package utils

import (
	"reflect"
	"strconv"
	"time"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func GetTypeName(v interface{}) string {
	valueOf := reflect.ValueOf(v)
	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	} else {
		return valueOf.Type().Name()
	}
}
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 0, 64)
}
func TsToTime(ts float64) time.Time {
	return time.Unix(int64(ts/1000), 0)
}
