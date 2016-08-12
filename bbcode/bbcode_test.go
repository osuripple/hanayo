package bbcode_test

import (
	"testing"

	"git.zxq.co/ripple/hanayo/bbcode"
)

func TestParse(t *testing.T) {
	w := []string{
		"[b]Test[/b]",
		"[b]Test",
		"Xdddd",
		"x[b=meme]Test[/b]",
		"[/^]",
		"[/meme]",
		"[!]meme[/!]",
		// Few posts I made on the official osu! server to test out.
		`... [url=https://github.com/Imvoo/GOsu]...[/url] ... [heading][url=https://github.com/thehowl/go-osuapi]GitHub[/url] | [url=https://godoc.org/gopkg.in/thehowl/go-osuapi.v1]Documentation[/url] | [url=https://github.com/thehowl/whosu]Sample application[/url][/heading]`,
		`[notice][heading]x[/heading]...[size=150][i]...[/i][/size][heading]...[/heading][size=150][centre][b][url=http://olc.howl.moe/]Website[/url][/b][/size] [size=150] | [b][url=https://github.com/TheHowl/OsuLevelCalculator]GitHub[/url][/b] | [b][url=https://github.com/TheHowl/OsuLevelCalculator/releases/tag/v1.2]Release[/url][/b] | [b][url=https://github.com/TheHowl/OsuLevelCalculator/releases/download/v1.1/OsuLevelCalculator_v1.2.zip]Release (zip file)[/url][/b][/centre][/size][/notice]`,
		`[quote="Howl"][quote="T3R4BYT3"]...[/quote]...[/quote]`,
	}
	for _, x := range w {
		bbc, err := bbcode.Parse(x)
		if err != nil {
			t.Error("err:", err)
		}
		t.Logf("%v", bbc)
	}
}
