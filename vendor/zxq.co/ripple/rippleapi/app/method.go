package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"strings"

	"zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
)

// Method wraps an API method to a HandlerFunc.
func Method(f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		initialCaretaker(c, f, privilegesNeeded...)
	}
}

func initialCaretaker(c *gin.Context, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
	rateLimiter()

	var doggoTags []string

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	token := ""
	switch {
	case c.Request.Header.Get("X-Ripple-Token") != "":
		token = c.Request.Header.Get("X-Ripple-Token")
	case c.Query("token") != "":
		token = c.Query("token")
	case c.Query("k") != "":
		token = c.Query("k")
	default:
		token, _ = c.Cookie("rt")
	}
	c.Set("token", fmt.Sprintf("%x", md5.Sum([]byte(token))))

	md := common.MethodData{
		DB:          db,
		RequestData: data,
		C:           c,
		Doggo:       doggo,
		R:           red,
	}
	if token != "" {
		tokenReal, exists := GetTokenFull(token, db)
		if exists {
			md.User = tokenReal
			doggoTags = append(doggoTags, "authorised")
		}
	}

	var ip string
	if requestIP, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err != nil {
		panic(err)
	} else {
		// if requestIP is not 127.0.0.1, means no reverse proxy is being used => direct request.
		if requestIP != "127.0.0.1" {
			ip = requestIP
		}
	}

	// means we're using reverse-proxy, so X-Real-IP
	if ip == "" {
		ip = c.ClientIP()
	}

	// requests from hanayo should not be rate limited.
	if !(c.Request.Header.Get("H-Key") == cf.HanayoKey && c.Request.UserAgent() == "hanayo") {
		perUserRequestLimiter(md.ID(), c.ClientIP())
	} else {
		doggoTags = append(doggoTags, "hanayo")
	}

	doggo.Incr("requests.v1", doggoTags, 1)

	missingPrivileges := 0
	for _, privilege := range privilegesNeeded {
		if uint64(md.User.TokenPrivileges)&uint64(privilege) == 0 {
			missingPrivileges |= privilege
		}
	}
	if missingPrivileges != 0 {
		c.IndentedJSON(401, common.SimpleResponse(401, "You don't have the privilege(s): "+common.Privileges(missingPrivileges).String()+"."))
		return
	}

	resp := f(md)
	if md.HasQuery("pls200") {
		c.Writer.WriteHeader(200)
	} else {
		c.Writer.WriteHeader(resp.GetCode())
	}

	if md.HasQuery("callback") {
		c.Header("Content-Type", "application/javascript; charset=utf-8")
	} else {
		c.Header("Content-Type", "application/json; charset=utf-8")
	}

	mkjson(c, resp)
}

// Very restrictive, but this way it shouldn't completely fuck up.
var callbackJSONP = regexp.MustCompile(`^[a-zA-Z_\$][a-zA-Z0-9_\$]*$`)

// mkjson auto indents json, and wraps json into a jsonp callback if specified by the request.
// then writes to the gin.Context the data.
func mkjson(c *gin.Context, data interface{}) {
	exported, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		c.Error(err)
		exported = []byte(`{ "code": 500, "message": "something has gone really really really really really really wrong." }`)
	}
	cb := c.Query("callback")
	willcb := cb != "" &&
		len(cb) < 100 &&
		callbackJSONP.MatchString(cb)
	if willcb {
		c.Writer.Write([]byte("/**/ typeof " + cb + " === 'function' && " + cb + "("))
	}
	c.Writer.Write(exported)
	if willcb {
		c.Writer.Write([]byte(");"))
	}
}
