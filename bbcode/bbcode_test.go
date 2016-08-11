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
	}
	for _, x := range w {
		bbc, err := bbcode.Parse(x)
		if err != nil {
			t.Error("err:", err)
		}
		t.Logf("%v", bbc)
	}
}
