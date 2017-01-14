// Package ocl allows you to do calculation of osu! levels in Go.
package ocl

import "math"

// GetLevel calculates what's the level of a score. It will stop at the
// level 10,000, after which it will give up. If you want to calculate the
// level without brakes of any kind, use GetLevelWithMax(score, -1).
func GetLevel(score int64) int {
	return GetLevelWithMax(score, 10000)
}

// GetLevelWithMax calculates what's the level of a score, having a maximum
// level. Set brake to a negative number to free yourself from any brakes.
func GetLevelWithMax(score int64, brake int) int {
	i := 1
	for {
		if brake > 0 && i >= brake {
			return i
		}
		lScore := GetRequiredScoreForLevel(i)
		if score < lScore {
			return i - 1
		}
		i++
	}
}

// GetRequiredScoreForLevel retrieves the score required to reach a certain
// level.
func GetRequiredScoreForLevel(level int) int64 {
	if level <= 100 {
		if level > 1 {
			return int64(math.Floor(float64(5000)/3*(4*math.Pow(float64(level), 3)-3*math.Pow(float64(level), 2)-float64(level)) + math.Floor(1.25*math.Pow(1.8, float64(level)-60))))
		}
		return 1
	}
	return 26931190829 + 100000000000*int64(level-100)
}

// GetLevelPrecise gets a precise level, meaning that decimal digits are
// included. There isn't any maximum level.
func GetLevelPrecise(score int64) float64 {
	baseLevel := GetLevelWithMax(score, -1)
	baseLevelScore := GetRequiredScoreForLevel(baseLevel)
	scoreProgress := score - baseLevelScore
	scoreLevelDifference := GetRequiredScoreForLevel(baseLevel+1) - baseLevelScore
	res := float64(scoreProgress)/float64(scoreLevelDifference) + float64(baseLevel)
	if math.IsInf(res, 0) || math.IsNaN(res) {
		return 0
	}
	return res
}
