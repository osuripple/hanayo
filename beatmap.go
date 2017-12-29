package main

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type beatmapData struct {
	baseTemplateData
	SongName string
}

func beatmapInfo(c *gin.Context) {
	var songName string

	b := c.Param("bid")
	if _, err := strconv.Atoi(b); err == nil {
		err := db.QueryRow("SELECT song_name FROM beatmaps WHERE beatmap_id = ? LIMIT 1", b).Scan(&songName)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	}

	data := new(beatmapData)
	data.SongName = songName
	defer resp(c, 200, "beatmap.html", data)
}
