package stadfangaskra

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

//
func getSingleQueryValueOrEmpty(v url.Values, param string) (string, error) {

	if qval, ok := v[param]; ok {
		if len(qval) > 1 {
			return "", errors.New(fmt.Sprintf("Only accepts a single query '%s' parameter, got %v", param, qval))
		}
		return qval[0], nil
	}

	return "", nil
}

func getQueryParamsAsInt(v url.Values, param string) (values []int) {

	if value, ok := v[param]; ok {
		for _, i := range value {
			v, err := strconv.Atoi(i)
			if err == nil {
				values = append(values, v)
			}
		}
	}
	return
}
