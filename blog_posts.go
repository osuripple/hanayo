package main

import "strconv"

// BlogPost ...
type BlogPost struct {
	Title       string
	Slug        string
	HTML        string
	CreatorID   int
	CreatorName string
}

func getBlogPosts(count int) []BlogPost {
	rows, err := db.Query("SELECT p.title, p.slug, p.html, p.author, u.username FROM anchor_posts p LEFT JOIN users u ON u.id = p.author ORDER BY p.id DESC LIMIT " + strconv.Itoa(count))
	if err != nil {
		return nil
	}
	posts := make([]BlogPost, 0, count)
	for rows.Next() {
		var p BlogPost
		err := rows.Scan(&p.Title, &p.Slug, &p.HTML, &p.CreatorID, &p.CreatorName)
		if err == nil {
			posts = append(posts, p)
		}
	}
	return posts
}
