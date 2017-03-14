package locale

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/leonelquinteros/gotext"
)

var languageMap = make(map[string]*gotext.Po, 20)

func loadLanguages() {
	files, err := ioutil.ReadDir("./data/locales")
	if err != nil {
		fmt.Println("loadLanguages", err)
		return
	}
	for _, file := range files {
		if file.Name() == "templates.pot" || file.Name() == "." || file.Name() == ".." {
			continue
		}

		po := new(gotext.Po)
		po.ParseFile("./data/locales/" + file.Name())

		langName := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".po"), "templates-")
		languageMap[langName] = po
	}
}

func init() {
	loadLanguages()
}

// Get retrieves a string from a language
func Get(langs []string, str string, vars ...interface{}) string {
	for _, lang := range langs {
		l := languageMap[lang]
		if l != nil {
			return l.Get(str, vars...)
		}
	}

	return fmt.Sprintf(str, vars...)
}
