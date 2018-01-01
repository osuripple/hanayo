package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/cheesegull/models"
)

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    models.Beatmap
	Beatmapset models.Set
	SetJSON    string
}

func beatmapInfo(c *gin.Context) {
	var (
		beatmap      models.Beatmap
		bset         models.Set
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
			return
		}
		bset, err = getBeatmapSetData(beatmap)
		if err != nil {
			c.Error(err)
			return
		}
		beatmapFound = true
		sort.Slice(bset.ChildrenBeatmaps, func(i, j int) bool {
			if bset.ChildrenBeatmaps[i].Mode != bset.ChildrenBeatmaps[j].Mode {
				return bset.ChildrenBeatmaps[i].Mode < bset.ChildrenBeatmaps[j].Mode
			}
			return bset.ChildrenBeatmaps[i].DifficultyRating < bset.ChildrenBeatmaps[j].DifficultyRating
		})
	}

	data.Found = beatmapFound
	if !beatmapFound {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	data.Beatmap = beatmap
	data.Beatmapset = bset

	setJson, err := json.Marshal(bset)
	if err == nil {
		data.SetJSON = fmt.Sprintf("%s", setJson)
	} else {
		data.SetJSON = "[]"
	}

	data.TitleBar = T(c, "%s - %s", bset.Artist, bset.Title)
	data.Scripts = append(data.Scripts, "/static/beatmap.js")
}

func getBeatmapData(b string) (beatmap models.Beatmap, err error) {
	resp, err := http.Get(config.CheesegullAPI + "/b/" + b)
	if err != nil {
		return beatmap, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return beatmap, err
	}

	err = json.Unmarshal(body, &beatmap)
	if err != nil {
		return beatmap, err
	}

	return beatmap, nil
}

func getBeatmapSetData(beatmap models.Beatmap) (bset models.Set, err error) {
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

	return bset, nil
}
