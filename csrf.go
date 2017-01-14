package main

// functions to check csrf

import (
	"strconv"
	"time"

	"zxq.co/x/rs"
	"github.com/thehowl/cieca"
)

var cStore = new(cieca.DataStore)

func csrfGenerate(u int) string {
	var s string
	for {
		s = rs.String(10)
		_, e := cStore.GetWithExist(s)
		if !e {
			break
		}
	}
	cStore.SetWithExpiration(strconv.Itoa(u)+s, nil, time.Minute*15)
	return s
}
func csrfExist(u int, token string) bool {
	_, e := cStore.GetWithExist(strconv.Itoa(u) + token)
	if e {
		cStore.Delete(strconv.Itoa(u) + token)
	}
	return e
}
