// Package bbcode implements bbcode parsing and translation into HTML.
//
// In this parser, a BBCode tag is defined in EBNF as following:
//
//  BBCode          = { Letter | Tag } .
//  Tag             = "[" Name [ "=" Attribute ] "]" BBCode "[/" Name "]" . // With both the first and the last Name being the same
//  Attribute       = Letter { Letter } .
//  Letter          = /* Everything except [ or ] */ .
//  NameLetter      = "0" … "9" | "A" … "Z" | "a" … "z" .
//  Name            = NameLetter { NameLetter } .
//
// TODO: Update EBNF thing.
package bbcode

import "fmt"

// BBCode is a slice of elements being either a String or a Tag.
type BBCode []BBCodeElement

func (b BBCode) lastOfChain(toAppendString bool) BBCodeElement {
	if len(b) == 0 {
		return nil
	}
	el := b[len(b)-1]
	switch el := el.(type) {
	case *String:
		return el
	case *Tag:
		loc := el.Content.lastOfChain(toAppendString)
		if loc == nil {
			if toAppendString {
				s := new(String)
				el.Content = append(el.Content, s)
				return s
			}
			t := new(Tag)
			el.Content = append(el.Content, t)
			return t
		}
		if _, ok := loc.(*String); ok && !toAppendString {
			t := new(Tag)
			el.Content = append(el.Content, t)
			return t
		}
		return loc
	}
	return nil
}

// BBCodeElement is one of String or Tag.
type BBCodeElement interface {
	BBCodeElement()
}

// String is just a string, but it satisfies the BBCodeElement interface.
type String string

// BBCodeElement is implemented just to satisfy BBCodeElement.
func (s *String) BBCodeElement() {}

func (s String) String() string {
	return string(s)
}

// Tag is a tag in BBCode.
type Tag struct {
	Name      string
	Attribute string
	Content   BBCode
	parent    *Tag
}

func (t Tag) String() string {
	return fmt.Sprintf("{tag %s, attr %s, val %v}", t.Name, t.Attribute, t.Content)
}

// BBCodeElement is implemented just to satisfy BBCodeElement.
func (t *Tag) BBCodeElement() {}
