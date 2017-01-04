package mysql

import (
	"testing"
)

func TestUser(t *testing.T) {
	_, err := sp().User(1000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserByEmail(t *testing.T) {
	_, err := sp().UserByEmail("fo@kab.ot")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCountry(t *testing.T) {
	_, err := sp().GetCountry(1000)
	if err != nil {
		t.Fatal(err)
	}
}
