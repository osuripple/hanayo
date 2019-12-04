// Package getrank allows retrieval of the rank of a score on
// an osu! beatmap (rank as in one of SSH, SS, SH, S, A, B, C, D)
package getrank

import (
	"gopkg.in/thehowl/go-osuapi.v1"
)

const silver = osuapi.ModFlashlight | osuapi.ModHidden

// GetRank retrieves a rank of a score with the passed arguments.
func GetRank(gameMode osuapi.Mode, mods osuapi.Mods, acc float64, c300, c100, c50, cmiss int) string {
	total := c300 + c100 + c50 + cmiss

	switch gameMode {
	case osuapi.ModeOsu, osuapi.ModeTaiko:
		var (
			c300f   = float64(c300)
			totalf  = float64(total)
			perc300 = c300f / totalf
		)
		switch {
		case acc == 100:
			return s(true, mods)
		case perc300 > 0.90 && float64(c50)/totalf < 0.1 && cmiss == 0:
			return s(false, mods)
		case (perc300 > 0.80 && cmiss == 0) || (perc300 > 0.90):
			return "a"
		case (perc300 > 0.70 && cmiss == 0) || (perc300 > 0.80):
			return "b"
		case perc300 > 0.60:
			return "c"
		}
		return "d"

	case osuapi.ModeCatchTheBeat:
		if acc == 100 {
			if mods&silver > 0 {
				return "sshd"
			}
			return "ss"
		}

		if acc >= 98.01 && acc <= 99.99 {
			if mods&silver > 0 {
				return "shd"
			}
			return "s"
		}

		if acc >= 94.01 && acc <= 98.00 {
			return "a"
		}

		if acc >= 90.01 && acc <= 94.00 {
			return "b"
		}

		if acc >= 85.01 && acc <= 90.00 {
			return "c"
		}

		return "d"

	case osuapi.ModeOsuMania:
		switch {
		case acc == 100:
			if mods&silver > 0 {
				return "sshd"
			}
			return "ss"
		case acc > 95:
			if mods&silver > 0 {
				return "shd"
			}
			return "s"
		case acc > 90:
			return "a"
		case acc > 80:
			return "b"
		case acc > 70:
			return "c"
		}
		return "d"
	}
	return "a"
}

func s(s2 bool, h osuapi.Mods) string {
	a := "s"
	if s2 {
		a += "s"
	}
	if h&silver > 0 {
		a += "h"
	}
	return a
}
