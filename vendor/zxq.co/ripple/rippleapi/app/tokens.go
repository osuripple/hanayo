package app

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"

	"zxq.co/ripple/rippleapi/common"
	"github.com/jmoiron/sqlx"
)

// GetTokenFull retrieves an user ID and their token privileges knowing their API token.
func GetTokenFull(token string, db *sqlx.DB) (common.Token, bool) {
	var (
		t             common.Token
		tokenPrivsRaw uint64
		userPrivsRaw  uint64
		priv8         bool
	)
	err := db.QueryRow(`SELECT 
	t.id, t.user, t.privileges, t.private, u.privileges
FROM tokens t
LEFT JOIN users u ON u.id = t.user
WHERE token = ? LIMIT 1`,
		fmt.Sprintf("%x", md5.Sum([]byte(token)))).
		Scan(
			&t.ID, &t.UserID, &tokenPrivsRaw, &priv8, &userPrivsRaw,
		)
	updateTokens <- t.ID
	if priv8 {
		// all privileges, they'll get removed by canOnly anyway.
		tokenPrivsRaw = (common.PrivilegeBeatmap << 1) - 1
	}
	t.UserPrivileges = common.UserPrivileges(userPrivsRaw)
	t.TokenPrivileges = common.Privileges(tokenPrivsRaw).CanOnly(t.UserPrivileges)
	switch {
	case err == sql.ErrNoRows:
		return common.Token{}, false
	case err != nil:
		panic(err)
	default:
		t.Value = token
		return t, true
	}
}

var updateTokens = make(chan int, 100)

func tokenUpdater(db *sqlx.DB) {
	for {
		// prepare a queue of tokens to update.
		tokensToUpdate := make([]int, 0, 50)
	AwaitLoop:
		for {
			// if we got ten, move on and update
			if len(tokensToUpdate) >= 50 {
				break
			}
			// if we ain't got any, add what we get straight from updateTokens
			if len(tokensToUpdate) == 0 {
				x := <-updateTokens
				tokensToUpdate = append(tokensToUpdate, x)
				continue
			}

			// otherwise, wait from updateTokens with a timeout of 10 seconds
			select {
			case x := <-updateTokens:
				tokensToUpdate = append(tokensToUpdate, x)
			case <-time.After(10 * time.Second):
				// wondering what this means?
				// https://golang.org/ref/spec#Break_statements
				break AwaitLoop
			}
		}

		q, a, _ := sqlx.In("UPDATE tokens SET last_updated = ? WHERE id IN (?)", time.Now().Unix(), tokensToUpdate)

		q = db.Rebind(q)
		_, err := db.Exec(q, a...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
