package app

import (
	"strconv"
	"time"

	"zxq.co/ripple/rippleapi/limit"
)

const reqsPerSecond = 5000
const sleepTime = time.Second / reqsPerSecond

var limiter = make(chan struct{}, reqsPerSecond)

func setUpLimiter() {
	for i := 0; i < reqsPerSecond; i++ {
		limiter <- struct{}{}
	}
	go func() {
		for {
			limiter <- struct{}{}
			time.Sleep(sleepTime)
		}
	}()
}

func rateLimiter() {
	<-limiter
}
func perUserRequestLimiter(uid int, ip string) {
	if uid == 0 {
		limit.Request("ip:"+ip, 60)
	} else {
		limit.Request("user:"+strconv.Itoa(uid), 2000)
	}
}
