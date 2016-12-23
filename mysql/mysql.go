// Package mysql implements hanayo's domain types with MySQL and sqlx.
package mysql

import "github.com/jmoiron/sqlx"

// ServiceProvider is the struct implementing all the services specified
// in the hanayo package that can be used with MySQL.
//
// DB must be a valid DB connection through sqlx.
type ServiceProvider struct {
	DB *sqlx.DB
}
