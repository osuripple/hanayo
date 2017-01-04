package mysql

import (
	"testing"
)

func TestMapperFunc(t *testing.T) {
	tt := []struct{ want, give string }{
		{"id", "ID"},
		{"api", "API"},
		{"user_id", "UserID"},
		{"registered_on", "RegisteredOn"},
	}
	for _, x := range tt {
		v := MapperFunc(x.give)
		if v != x.want {
			t.Errorf("want %q, got %q", x.want, v)
		}
	}
}
