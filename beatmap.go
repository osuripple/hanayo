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
	"zxq.co/ripple/rippleapi/app/v1"
	"zxq.co/ripple/rippleapi/common"
)

// idk if this shud b exported, just getting it to work for now.
// TODO: figure ^ out

type userData struct {
	ID             int                  `json:"id"`
	Username       string               `json:"username"`
	UsernameAKA    string               `json:"username_aka"`
	RegisteredOn   common.UnixTimestamp `json:"registered_on"`
	Privileges     uint64               `json:"privileges"`
	LatestActivity common.UnixTimestamp `json:"latest_activity"`
	Country        string               `json:"country"`
}

type beatmapScore struct {
	v1.Score
	User userData `json:"user"`
}

type scoresResponse struct {
	common.ResponseBase
	Scores []beatmapScore `json:"scores"`
}

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    models.Beatmap
	Beatmapset models.Set
	Scores     []beatmapScore
}

type beatmapsList []models.Beatmap

func (s beatmapsList) Len() int {
	return len(s)
}

func (s beatmapsList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s beatmapsList) Less(i, j int) bool {
	if s[i].Mode != s[j].Mode {
		return s[i].Mode < s[j].Mode
	}
	return s[i].DifficultyRating < s[j].DifficultyRating
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
		sort.Sort(beatmapsList(bset.ChildrenBeatmaps))
	}

	data.Found = beatmapFound
	if !beatmapFound {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	scores, err := getScoresData(beatmap)
	if err != nil {
		data.Messages = append(data.Messages, errorMessage{T(c, "Could not retrieve scores for this map.")})
		c.Error(err)
	} else {
		data.Scores = scores
	}

	data.Beatmap = beatmap
	data.Beatmapset = bset
	data.TitleBar = T(c, "%s - %s", bset.Artist, bset.Title)
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

func getScoresData(beatmap models.Beatmap) (scores []beatmapScore, err error) {
	scoreResp := new(scoresResponse)

	resp, err := http.Get(config.API + "scores?b=" + strconv.Itoa(beatmap.ID) + "&p=1&l=50")
	if err != nil {
		return scores, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return scores, err
	}
	fmt.Printf("%#v\n", body)

	err = json.Unmarshal(body, &scoreResp)
	if err != nil {
		return scores, err
	}
	fmt.Printf("%#v\n", scores)

	return scoreResp.Scores, nil
}
