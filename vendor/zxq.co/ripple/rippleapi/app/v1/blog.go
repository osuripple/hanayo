package v1

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"zxq.co/ripple/rippleapi/common"
)

// This basically proxies requests from Medium's API and is used on Ripple's
// home page to display the latest blog posts.

type mediumResp struct {
	Success bool `json:"success"`
	Payload struct {
		Posts []mediumPost `json:"posts"`
	} `json:"payload"`
}

type mediumPost struct {
	ID          string             `json:"id"`
	Creator     mediumUser         `json:"creator"`
	Title       string             `json:"title"`
	CreatedAt   int64              `json:"createdAt"`
	UpdatedAt   int64              `json:"updatedAt"`
	Virtuals    mediumPostVirtuals `json:"virtuals"`
	ImportedURL string             `json:"importedUrl"`
	UniqueSlug  string             `json:"uniqueSlug"`
}

type mediumUser struct {
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type mediumPostVirtuals struct {
	CreatedAtRelative string  `json:"createdAtRelative"`
	UpdatedAtRelative string  `json:"updatedAtRelative"`
	Snippet           string  `json:"snippet"`
	WordCount         int     `json:"wordCount"`
	ReadingTime       float64 `json:"readingTime"`
}

// there's gotta be a better way

type blogPost struct {
	ID          string    `json:"id"`
	Creator     blogUser  `json:"creator"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImportedURL string    `json:"imported_url"`
	UniqueSlug  string    `json:"unique_slug"`

	CreatedAtRelative string  `json:"created_at_relative"`
	UpdatedAtRelative string  `json:"updated_at_relative"`
	Snippet           string  `json:"snippet"`
	WordCount         int     `json:"word_count"`
	ReadingTime       float64 `json:"reading_time"`
}

type blogUser struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type blogPostsResponse struct {
	common.ResponseBase
	Posts []blogPost `json:"posts"`
}

// consts for the medium API
const (
	mediumAPIResponsePrefix = `])}while(1);</x>`
	mediumAPIAllPosts       = `https://blog.ripple.moe/all?format=json`
)

func init() {
	gob.Register([]blogPost{})
}

// BlogPostsGET retrieves the latest blog posts on the Ripple blog.
func BlogPostsGET(md common.MethodData) common.CodeMessager {
	// check if posts are cached in redis
	res := md.R.Get("api:blog_posts").Val()
	if res != "" {
		// decode values
		posts := make([]blogPost, 0, 20)
		err := gob.NewDecoder(strings.NewReader(res)).Decode(&posts)
		if err != nil {
			md.Err(err)
			return Err500
		}

		// create response and return
		var r blogPostsResponse
		r.Code = 200
		r.Posts = blogLimit(posts, md.Query("l"))
		return r
	}

	// get data from medium api
	resp, err := http.Get(mediumAPIAllPosts)
	if err != nil {
		md.Err(err)
		return Err500
	}

	// read body and trim the prefix
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		md.Err(err)
		return Err500
	}
	all = bytes.TrimPrefix(all, []byte(mediumAPIResponsePrefix))

	// unmarshal into response struct
	var mResp mediumResp
	err = json.Unmarshal(all, &mResp)
	if err != nil {
		md.Err(err)
		return Err500
	}

	if !mResp.Success {
		md.Err(errors.New("medium api call is not successful"))
		return Err500
	}

	// create posts slice and fill it up with converted posts from the medium
	// API
	posts := make([]blogPost, len(mResp.Payload.Posts))
	for idx, mp := range mResp.Payload.Posts {
		var p blogPost

		// convert structs
		p.ID = mp.ID
		p.Title = mp.Title
		p.CreatedAt = time.Unix(0, mp.CreatedAt*1000000)
		p.UpdatedAt = time.Unix(0, mp.UpdatedAt*1000000)
		p.ImportedURL = mp.ImportedURL
		p.UniqueSlug = mp.UniqueSlug

		p.Creator.UserID = mp.Creator.UserID
		p.Creator.Name = mp.Creator.Name
		p.Creator.Username = mp.Creator.Username

		p.CreatedAtRelative = mp.Virtuals.CreatedAtRelative
		p.UpdatedAtRelative = mp.Virtuals.UpdatedAtRelative
		p.Snippet = mp.Virtuals.Snippet
		p.WordCount = mp.Virtuals.WordCount
		p.ReadingTime = mp.Virtuals.ReadingTime

		posts[idx] = p
	}

	// save in redis
	bb := new(bytes.Buffer)
	err = gob.NewEncoder(bb).Encode(posts)
	if err != nil {
		md.Err(err)
		return Err500
	}
	md.R.Set("api:blog_posts", bb.Bytes(), time.Minute*5)

	var r blogPostsResponse
	r.Code = 200
	r.Posts = blogLimit(posts, md.Query("l"))
	return r
}

func blogLimit(posts []blogPost, s string) []blogPost {
	i := common.Int(s)
	if i >= len(posts) || i < 1 {
		return posts
	}
	return posts[:i]
}
