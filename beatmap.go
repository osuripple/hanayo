package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type beatmapData struct {
	AR float32
	CS float32
	HP float32
	OD float32

	BeatmapID   int
	ParentSetID int

	BPM              float32
	DiffName         string
	DifficultyRating float32
	FileMD5          string
	HitLength        int
	MaxCombo         int
	Mode             int
	Passcount        int
	Playcount        int
	TotalLength      int
}

type beatmapSetData struct {
	ApprovedDate     string // todo an actual date
	Artist           string
	ChildrenBeatmaps []*beatmapData
	Creator          string
	Favourites       int
	Genre            int
	HasVideo         bool
	Language         int
	LastChecked      string
	LastUpdate       string
	RankedStatus     int
	SetID            int
	Source           string
	Tags             string
	Title            string
}

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    *beatmapData
	Beatmapset *beatmapSetData
}

func beatmapInfo(c *gin.Context) {
	var (
		beatmap      *beatmapData
		bset         *beatmapSetData
		beatmapFound bool
	)

	data := new(beatmapPageData)
	defer resp(c, 200, "beatmap.html", data)

	b := c.Param("bid")
	if _, err := strconv.Atoi(b); err != nil {
		c.Error(err)
	} else {
		beatmap, err = getBeatmapData(b)
		if err != nil {
			c.Error(err)
		} else {
			bset, err = getBeatmapSetData(beatmap)
			if err != nil {
				c.Error(err)
			} else {
				fmt.Printf("set: %#v\n", bset)
				beatmapFound = true
			}
		}
	}

	data.Found = beatmapFound
	if !beatmapFound {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	data.Beatmap = beatmap
	data.Beatmapset = bset
	data.TitleBar = T(c, "%s - %s", bset.Artist, bset.Title)
}

func getBeatmapData(b string) (beatmap *beatmapData, err error) {
	obj := new(beatmapData)

	resp, err := http.Get(config.CheesegullAPI + "/b/" + b)
	if err != nil {
		return obj, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func getBeatmapSetData(beatmap *beatmapData) (bset *beatmapSetData, err error) {
	obj := new(beatmapSetData)

	resp, err := http.Get(config.CheesegullAPI + "/s/" + strconv.Itoa(beatmap.ParentSetID))
	if err != nil {
		return obj, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}
