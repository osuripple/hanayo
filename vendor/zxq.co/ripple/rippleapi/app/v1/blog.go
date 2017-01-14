package v1

import (
	"time"

	"zxq.co/ripple/rippleapi/common"
)

type blogPost struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
	Author  userData  `json:"author"`
}

type blogPostsResponse struct {
	common.ResponseBase
	Posts []blogPost `json:"posts"`
}

// BlogPostsGET retrieves the latest blog posts on the Ripple blog.
func BlogPostsGET(md common.MethodData) common.CodeMessager {
	var and string
	var params []interface{}
	if md.Query("id") != "" {
		and = "b.id = ?"
		params = append(params, md.Query("id"))
	}
	rows, err := md.DB.Query(`
	SELECT 
		b.id, b.title, b.slug, b.created,
		
		u.id, u.username, s.username_aka, u.register_datetime,
		u.privileges, u.latest_activity, s.country
	FROM anchor_posts b
	INNER JOIN users u ON b.author = u.id
	INNER JOIN users_stats s ON b.author = s.id
	WHERE status = "published" `+and+`
	ORDER BY b.id DESC `+common.Paginate(md.Query("p"), md.Query("l"), 50), params...)
	if err != nil {
		md.Err(err)
		return Err500
	}

	var r blogPostsResponse
	for rows.Next() {
		var post blogPost
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Created,

			&post.Author.ID, &post.Author.Username, &post.Author.UsernameAKA, &post.Author.RegisteredOn,
			&post.Author.Privileges, &post.Author.LatestActivity, &post.Author.Country,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Posts = append(r.Posts, post)
	}
	r.Code = 200

	return r
}

type blogPostContent struct {
	common.ResponseBase
	Content string `json:"content"`
}

// BlogPostsContentGET retrieves the content of a specific blog post.
func BlogPostsContentGET(md common.MethodData) common.CodeMessager {
	field := "markdown"
	if md.HasQuery("html") {
		field = "html"
	}
	var (
		by  string
		val string
	)
	switch {
	case md.Query("slug") != "":
		by = "slug"
		val = md.Query("slug")
	case md.Query("id") != "":
		by = "id"
		val = md.Query("id")
	default:
		return ErrMissingField("id|slug")
	}
	var r blogPostContent
	err := md.DB.QueryRow("SELECT "+field+" FROM anchor_posts WHERE "+by+" = ? AND status = 'published'", val).Scan(&r.Content)
	if err != nil {
		return common.SimpleResponse(404, "no blog post found")
	}
	r.Code = 200
	return r
}
