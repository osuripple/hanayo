package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/osuripple/cheesegull/models"
)

type beatmapPageData struct {
	baseTemplateData

	Found        bool
	Beatmap      Beatmap
	Beatmapset   Set
	SetJSON      string
	MapJSON      string
	MapID        string
}

type Beatmap struct {
	ID               int `json:"BeatmapID"`
	ParentSetID      int
	DiffName         string
	FileMD5          string
	Mode             int
	BPM              float64
	AR               float32
	OD               float32
	CS               float32
	HP               float32
	TotalLength      int
	HitLength        int
	Playcount        int
	Passcount        int
	MaxCombo         int
	DifficultyRating float64
	Status           int
}

type Set struct {
	ID               int `json:"SetID"`
	ChildrenBeatmaps []Beatmap
	RankedStatus     int
	ApprovedDate     time.Time
	LastUpdate       time.Time
	LastChecked      time.Time
	Artist           string
	Title            string
	Creator          string
	Source           string
	Tags             string
	HasVideo         bool
	Genre            int
	Language         int
	Favourites       int
}

func beatmapInfo(c *gin.Context) {
	data := new(beatmapPageData)
	defer resp(c, 200, "beatmap.html", data)

	b := c.Param("bid")
	if _, err := strconv.Atoi(b); err != nil {
		c.Error(err)
	} else {
		data.MapID = c.Param("bid")
		if err != nil {
            c.Error(err)
			return
		}
		data.Beatmap, err = getBeatmapData(b)
		if err != nil {
			c.Error(err)
			return
		}
		data.Beatmapset, err = getBeatmapSetData(data.Beatmap)
		if err != nil {
			c.Error(err)
			return
		}
		sort.Slice(data.Beatmapset.ChildrenBeatmaps, func(i, j int) bool {
			if data.Beatmapset.ChildrenBeatmaps[i].Mode != data.Beatmapset.ChildrenBeatmaps[j].Mode {
				return data.Beatmapset.ChildrenBeatmaps[i].Mode < data.Beatmapset.ChildrenBeatmaps[j].Mode
			}
			return data.Beatmapset.ChildrenBeatmaps[i].DifficultyRating < data.Beatmapset.ChildrenBeatmaps[j].DifficultyRating
		})
	}

	if data.Beatmapset.ID == 0 {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	data.KyutGrill = fmt.Sprintf("https://assets.ppy.sh/beatmaps/%d/covers/cover.jpg?%d", data.Beatmapset.ID, data.Beatmapset.LastUpdate.Unix())
	data.KyutGrillAbsolute = true

	setJSON, err := json.Marshal(data.Beatmapset)
	if err == nil {
		data.SetJSON = string(setJSON)
	} else {
		data.SetJSON = "[]"
	}

	mapJSON, err := json.Marshal(data.Beatmap)
	if err == nil {
		data.MapJSON = string(mapJSON)
	} else {
		data.MapJSON = "[]"
	}

	data.TitleBar = T(c, "%s - %s", data.Beatmapset.Artist, data.Beatmapset.Title)
	data.Scripts = append(data.Scripts, "/static/tablesort.js", "/static/beatmap.js")
}

func getBeatmapData(b string) (beatmap Beatmap, err error) {
	// Get beatmap data from Cheesegull API
	resp, err := http.Get(config.CheesegullAPI + "/b/" + b)
	if err != nil {
		return beatmap, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return beatmap, err
	}

	// Unmarshal beatmap data
	err = json.Unmarshal(body, &beatmap)
	if err != nil {
		return beatmap, err
	}

	// Get beatmap data from Kawata API to get approved value
	APIResp, err := http.Get(config.BaseURL + "/api/v1/get_beatmaps?b=" + b)
	if err != nil {
		return beatmap, err
	}
	defer APIResp.Body.Close()
	APIBody, err := ioutil.ReadAll(APIResp.Body)
	if err != nil {
		return beatmap, err
	}

	// Unmarshal API response to get approved value
	var apiBeatmap []struct {
		Approved string `json:"approved"`
	}
	err = json.Unmarshal(APIBody, &apiBeatmap)
	if err != nil {
		return beatmap, err
	}

	// Set approved value in beatmap object
	if len(apiBeatmap) > 0 {
		approved, err := strconv.Atoi(apiBeatmap[0].Approved)
		if err == nil {
			beatmap.Status = approved
		}
	}

	return beatmap, nil
}

func getBeatmapSetData(beatmap Beatmap) (bset Set, err error) {
	resp, err := http.Get(config.CheesegullAPI + "/s/" + strconv.Itoa(beatmap.ParentSetID))
	if err != nil {
		return bset, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bset, err
	}

	err = json.Unmarshal(body, &bset)
	if err != nil {
		return bset, err
	}

	for i := range bset.ChildrenBeatmaps {
		childBeatmap := &bset.ChildrenBeatmaps[i]
		// Get beatmap data from Kawata API to get approved value
		APIResp, err := http.Get(config.BaseURL + "/api/v1/get_beatmaps?b=" + strconv.Itoa(childBeatmap.ID))
		if err != nil {
			return bset, err
		}
		defer APIResp.Body.Close()
		APIBody, err := ioutil.ReadAll(APIResp.Body)
		if err != nil {
			return bset, err
		}

		// Unmarshal API response to get approved value
		var apiBeatmap []struct {
			Approved string `json:"approved"`
		}
		err = json.Unmarshal(APIBody, &apiBeatmap)
		if err != nil {
			return bset, err
		}

		// Set approved value in child beatmap object
		if len(apiBeatmap) > 0 {
			approved, err := strconv.Atoi(apiBeatmap[0].Approved)
			if err == nil {
				childBeatmap.Status = approved
			}
		}
	}

	return bset, nil
}
