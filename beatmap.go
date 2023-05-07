package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/osuripple/cheesegull/models"
)

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    models.Beatmap
	Beatmapset models.Set
	SetJSON    string
}

func beatmapInfo(c *gin.Context) {
	data := new(beatmapPageData)
	defer resp(c, 200, "beatmap.html", data)

	sid := c.Param("sid")
	b := ""
	// in this scenario, a sid looks like "<numbers>#[osu|taiko|fruits|mania]"
	// where the #<mode> is optional
	// we need to split the mode from the sid

	// if the sid is just numbers, it's a beatmapset id
	// if the sid is numbers followed by a # and a mode, there is a beatmap id

	for _, mode := range []string{"osu", "taiko", "fruits", "mania"} {
		if strings.HasSuffix(sid, "#"+mode) {
			sid = strings.TrimSuffix(sid, "#"+mode)
			b = c.Param("bid")
			break
		}
	}

	if b == "" {
		b = c.Param("bid")
	}


	log.Println("beatmapInfo", b, sid)

	if _, err := strconv.Atoi(b); err == nil {
		// try beatmap id first
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
	} else if sidNum, err := strconv.Atoi(sid); err == nil {
		// fallback to set id
		data.Beatmapset, err = getBeatmapSetData(models.Beatmap{ParentSetID: sidNum})
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

		data.Beatmap = data.Beatmapset.ChildrenBeatmaps[0]
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

	data.TitleBar = T(c, "%s - %s", data.Beatmapset.Artist, data.Beatmapset.Title)
	data.Scripts = append(data.Scripts, "/static/tablesort.js", "/static/beatmap.js")
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
