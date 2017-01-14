package v1

import (
	"database/sql"

	"zxq.co/ripple/rippleapi/common"
)

type docFile struct {
	ID      int    `json:"id"`
	DocName string `json:"doc_name"`
	Public  bool   `json:"public"`
	IsRule  bool   `json:"is_rule"`
}

type docResponse struct {
	common.ResponseBase
	Files []docFile `json:"files"`
}

// DocGET retrieves a list of documentation files.
func DocGET(md common.MethodData) common.CodeMessager {
	var wc string
	if md.User.TokenPrivileges&common.PrivilegeBlog == 0 || md.Query("public") == "1" {
		wc = "WHERE public = '1'"
	}
	rows, err := md.DB.Query("SELECT id, doc_name, public, is_rule FROM docs " + wc)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r docResponse
	for rows.Next() {
		var f docFile
		err := rows.Scan(&f.ID, &f.DocName, &f.Public, &f.IsRule)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Files = append(r.Files, f)
	}
	r.Code = 200
	return r
}

type docContentResponse struct {
	common.ResponseBase
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DocContentGET retrieves the raw markdown file of a doc file
func DocContentGET(md common.MethodData) common.CodeMessager {
	docID := common.Int(md.Query("id"))
	if docID == 0 {
		return common.SimpleResponse(404, "Documentation file not found!")
	}
	var wc string
	if md.User.TokenPrivileges&common.PrivilegeBlog == 0 || md.Query("public") == "1" {
		wc = "AND public = '1'"
	}
	var r docContentResponse
	err := md.DB.QueryRow("SELECT doc_name, doc_contents FROM docs WHERE id = ? "+wc+" LIMIT 1", docID).Scan(&r.Title, &r.Content)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "Documentation file not found!"
	case err != nil:
		md.Err(err)
		return Err500
	default:
		r.Code = 200
	}
	return r
}

// DocRulesGET gets the rules.
func DocRulesGET(md common.MethodData) common.CodeMessager {
	var r docContentResponse
	r.Title = "Rules"
	err := md.DB.QueryRow("SELECT doc_contents FROM docs WHERE is_rule = '1' LIMIT 1").Scan(&r.Content)
	const ruleFree = "# This Ripple instance is rule-free! Yay!"
	switch {
	case err == sql.ErrNoRows:
		r.Content = ruleFree
	case err != nil:
		md.Err(err)
		return Err500
	case len(r.Content) == 0:
		r.Content = ruleFree
	}
	r.Code = 200
	return r
}
