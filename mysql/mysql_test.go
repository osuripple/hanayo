package mysql

import (
	"os"

	"github.com/jmoiron/sqlx"
)

var _sp *ServiceProvider

// sp retrieves the test service provider.
func sp() *ServiceProvider {
	if _sp != nil {
		return _sp
	}
	e := "root@/ripple"
	if os.Getenv("TEST_DSN") != "" {
		e = os.Getenv("TEST_DSN")
	}
	db, err := sqlx.Open("mysql", e)
	if err != nil {
		panic(err)
	}
	db.MapperFunc(MapperFunc)
	_sp = &ServiceProvider{db}
	return _sp
}
