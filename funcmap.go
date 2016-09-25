package main

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
	"github.com/dustin/go-humanize"
)

// funcMap contains useful functions for the various templates.
var funcMap = template.FuncMap{
	// html disables HTML escaping on the values it is given.
	"html": func(value interface{}) template.HTML {
		return template.HTML(fmt.Sprint(value))
	},
	// avatars is a function returning the configuration constant AvatarURL
	// TODO: Replace with config function returning something from config.
	"avatars": func() string {
		return config.AvatarURL
	},
	// navbarItem is a function to generate an item in the navbar.
	// The reason why this exists is that I wanted to have the currently
	// selected element in the navbar having the "active" class.
	"navbarItem": func(currentPath, name, path string) template.HTML {
		var act string
		if path == currentPath {
			act = "active "
		}
		return template.HTML(fmt.Sprintf(`<a class="%sitem" href="%s">%s</a>`, act, path, name))
	},
	// curryear returns the current year.
	"curryear": func() string {
		return strconv.Itoa(time.Now().Year())
	},
	// hasAdmin returns, based on the user's privileges, whether they should be
	// able to see the RAP button (aka AdminPrivilegeAccessRAP).
	"hasAdmin": func(privs common.UserPrivileges) bool {
		return privs&common.AdminPrivilegeAccessRAP > 0
	},
	// isRAP returns whether the current page is in RAP.
	"isRAP": func(p string) bool {
		parts := strings.Split(p, "/")
		return len(parts) > 1 && parts[1] == "admin"
	},
	// favMode is just a helper function for user profiles. Basically checks
	// whether two floats are equal, and if they are it will return "active ",
	// so that the element in the mode menu of a user profile can be marked as
	// the current active element.
	"favMode": func(favMode, current float64) string {
		if favMode == current {
			return "active "
		}
		return ""
	},
	// slice generates a []interface{} with the elements it is given.
	// useful to iterate over some elements, like this:
	//  {{ range slice 1 2 3 }}{{ . }}{{ end }}
	"slice": func(els ...interface{}) []interface{} {
		return els
	},
	// int converts a float/int to an int.
	"int": func(f interface{}) int {
		if f == nil {
			return 0
		}
		switch f := f.(type) {
		case int:
			return f
		case float64:
			return int(f)
		case float32:
			return int(f)
		}
		return 0
	},
	// float converts an int to a float.
	"float": func(i int) float64 {
		return float64(i)
	},
	// atoi converts a string to an int and then a float64.
	// If s is not an actual int, it returns nil.
	"atoi": func(s string) interface{} {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		return float64(i)
	},
	// parseUserpage compiles BBCode to HTML.
	"parseUserpage": func(s string) template.HTML {
		return template.HTML(compileBBCode(s))
	},
	// time converts a RFC3339 timestamp to the HTML element <time>.
	"time": func(s string) template.HTML {
		t, _ := time.Parse(time.RFC3339, s)
		return template.HTML(fmt.Sprintf(`<time class="timeago" datetime="%s">%v</time>`, s, t))
	},
	// band is a bitwise AND.
	"band": func(i1 int, i ...int) int {
		for _, el := range i {
			i1 &= el
		}
		return i1
	},
	// countryReadable converts a country's ISO name to its full name.
	"countryReadable": countryReadable,
	"country": func(s string) template.HTML {
		c := countryReadable(s)
		if c == "" {
			return ""
		}
		return template.HTML(fmt.Sprintf(`<i class="%s flag smallpadd"></i> %s`, strings.ToLower(s), c))
	},
	// humanize pretty-prints a float, e.g.
	//     humanize(1000) == "1,000"
	"humanize": func(f float64) string {
		return humanize.Commaf(f)
	},
	// levelPercent basically does this:
	//     levelPercent(56.23215) == "23"
	"levelPercent": func(l float64) string {
		_, f := math.Modf(l)
		f *= 100
		return fmt.Sprintf("%.0f", f)
	},
	// level removes the decimal part from a float.
	"level": func(l float64) string {
		i, _ := math.Modf(l)
		return fmt.Sprintf("%.0f", i)
	},
	// trimPrefix returns s without the provided leading prefix string.
	// If s doesn't start with prefix, s is returned unchanged.
	"trimPrefix": strings.TrimPrefix,
	// log fmt.Printf's something
	"log": fmt.Printf,
	// has returns whether priv1 has all 1 bits of priv2, aka priv1 & priv2 == priv2
	"has": func(priv1, priv2 float64) bool {
		return uint64(priv1)&uint64(priv2) == uint64(priv2)
	},
	// _range is like python range's.
	// If it is given 1 argument, it returns a []int containing numbers from 0
	// to x.
	// If it is given 2 arguments, it returns a []int containing numers from x
	// to y if x < y, from y to x if y < x.
	"_range": func(x int, y ...int) ([]int, error) {
		switch len(y) {
		case 0:
			r := make([]int, x)
			for i := range r {
				r[i] = i
			}
			return r, nil
		case 1:
			nums, up := pos(y[0] - x)
			r := make([]int, nums)
			for i := range r {
				if up {
					r[i] = i + x + 1
				} else {
					r[i] = i + y[0]
				}
			}
			if !up {
				// reverse r
				sort.Sort(sort.Reverse(sort.IntSlice(r)))
			}
			return r, nil
		}
		return nil, errors.New("y must be at maximum 1 parameter")
	},
}

func pos(x int) (int, bool) {
	if x > 0 {
		return x, true
	}
	return x * -1, false
}
