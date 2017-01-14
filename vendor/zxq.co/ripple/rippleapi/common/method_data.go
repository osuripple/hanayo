package common

import (
	"encoding/json"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gopkg.in/redis.v5"
)

// MethodData is a struct containing the data passed over to an API method.
type MethodData struct {
	User        Token
	DB          *sqlx.DB
	RequestData RequestData
	C           *gin.Context
	Doggo       *statsd.Client
	R           *redis.Client
}

// Err logs an error into gin.
func (md MethodData) Err(err error) {
	md.C.Error(err)
}

// ID retrieves the Token's owner user ID.
func (md MethodData) ID() int {
	return md.User.UserID
}

// Query is shorthand for md.C.Query.
func (md MethodData) Query(q string) string {
	return md.C.Query(q)
}

// HasQuery returns true if the parameter is encountered in the querystring.
// It returns true even if the parameter is "" (the case of ?param&etc=etc)
func (md MethodData) HasQuery(q string) bool {
	_, has := md.C.GetQuery(q)
	return has
}

// RequestData is the body of a request. It is wrapped into this type
// to implement the Unmarshal function, which is just a shorthand to
// json.Unmarshal.
type RequestData []byte

// Unmarshal json-decodes Requestdata into a value. Basically a
// shorthand to json.Unmarshal.
func (r RequestData) Unmarshal(into interface{}) error {
	return json.Unmarshal([]byte(r), into)
}
