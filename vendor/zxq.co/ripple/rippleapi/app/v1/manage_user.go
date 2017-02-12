package v1

import (
	"time"

	"zxq.co/ripple/rippleapi/common"
)

type setAllowedData struct {
	UserID  int `json:"user_id"`
	Allowed int `json:"allowed"`
}

// UserManageSetAllowedPOST allows to set the allowed status of an user.
func UserManageSetAllowedPOST(md common.MethodData) common.CodeMessager {
	data := setAllowedData{}
	if err := md.Unmarshal(&data); err != nil {
		return ErrBadJSON
	}
	if data.Allowed < 0 || data.Allowed > 2 {
		return common.SimpleResponse(400, "Allowed status must be between 0 and 2")
	}
	var banDatetime int64
	var privsSet string
	if data.Allowed == 0 {
		banDatetime = time.Now().Unix()
		privsSet = "privileges = (privileges & ~3)"
	} else {
		banDatetime = 0
		privsSet = "privileges = (privileges | 3)"
	}
	_, err := md.DB.Exec("UPDATE users SET "+privsSet+", ban_datetime = ? WHERE id = ?", banDatetime, data.UserID)
	if err != nil {
		md.Err(err)
		return Err500
	}
	go fixPrivileges(data.UserID, md.DB)
	query := `
SELECT users.id, users.username, register_datetime, privileges,
	latest_activity, users_stats.username_aka,
	users_stats.country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.id=?
LIMIT 1`
	return userPutsSingle(md, md.DB.QueryRowx(query, data.UserID))
}
