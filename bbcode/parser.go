package bbcode

// Parse parses raw BBCode into the type BBCode.
func Parse(raw string) (BBCode, error) {
	const (
		writing = iota
		namingTag
		givingAttribute
		closingTag
	)

	var (
		status  int
		current = &Tag{
			Name:    "^",
			Content: make(BBCode, 0),
		}
		closingTagName string
	)

	invalidClosingTag := func(ch rune) {
		status = writing
		c := *current
		setToParent(&current).removeLast()
		str := "[" + current.Name
		if current.Attribute != "" {
			str += "=" + current.Attribute
		}
		str += "]"
		current.addString(str)
		current.Content = append(current.Content, c.Content...)
		str = "[/" + closingTagName + string(ch)
		current.addString(str)
		closingTagName = ""
	}

	for _, ch := range raw {
		switch status {
		case writing:
			// if we intend to start writing a new tag,
			// change the status and replace current
			if ch == '[' {
				status = namingTag
				_new := &Tag{
					Name:    "",
					Content: make(BBCode, 0),
					parent:  current,
				}
				current.Content = append(current.Content, _new)
				current = _new
				continue
			}
			current.addString(string(ch))
		case namingTag:
			// if we have not even started writing the tag and it starts with /,
			// it means we want to close a tag. change the status, change current tag
			// to be the parent, which is the one that should be then closed.
			if ch == '/' && current.Name == "" {
				status = closingTag
				setToParent(&current).removeLast()
				continue
			}
			// if we can switch to attribute, do it
			if ch == '=' {
				if current.Name == "" {
					// if it's [= it's not valid. Just output it and go back to writing mode in the parent tag.
					status = writing
					setToParent(&current).removeLast()
					current.addString("[=")
					continue
				}
				status = givingAttribute
			}
			if ch == ']' && current.Name != "" {
				status = writing
				continue
			}
			// if the letter is not valid, assume it's text of the upper parent.
			if !isValidNameLetter(ch) {
				status = writing
				strToAdd := "[" + current.Name + string(ch)
				setToParent(&current).removeLast()
				current.addString(strToAdd)
				continue
			}
			// passed all tests, it's a valid character for the name.
			current.Name += string(ch)
		case givingAttribute:
			if ch == ']' {
				status = writing
				continue
			}
			current.Attribute += string(ch)
		case closingTag:
			// if we're closing the closing tag, check it's the same as the opening tag.
			if ch == ']' && closingTagName != "" {
				if closingTagName != current.Name {
					invalidClosingTag(ch)
					continue
				}
				status = writing
				setToParent(&current)
				closingTagName = ""
				continue
			}
			// If the closing tag is not valid, let the parent inherit all
			// of the current's content, as well as the raw tags.
			if !isValidNameLetter(ch) {
				invalidClosingTag(ch)
				continue
			}
			closingTagName += string(ch)
		}
	}

	// unparse tags not closed
	for current.Name != "^" {
		c := *current
		setToParent(&current).removeLast()
		str := "[" + c.Name
		if c.Attribute != "" {
			str += "=" + c.Attribute
		}
		str += "]"
		current.addString(str)
		current.Content = append(current.Content, c.Content...)
	}

	return current.Content, nil
}

func (t *Tag) addString(str string) {
	// make sure last element of current.Content exists and is *String,
	// so that we can append our new character to it.
	if len(t.Content) == 0 {
		t.Content = append(t.Content, new(String))
	} else if _, ok := t.Content[len(t.Content)-1].(*String); !ok {
		t.Content = append(t.Content, new(String))
	}
	s, ok := t.Content[len(t.Content)-1].(*String)
	if !ok {
		panic("last element of t.Content is not a string, despite having made sure it would have been")
	}
	*s += String(str)
}

// yes. we need a double pointer.
// https://play.golang.org/p/3N8n-D4FXZ
// this almost feels like reading C
func setToParent(t **Tag) *Tag {
	*t = (**t).parent
	return *t
}
func (t *Tag) removeLast() {
	t.Content[len(t.Content)-1] = nil
	t.Content = t.Content[:len(t.Content)-1]
}

func isValidNameLetter(b rune) bool {
	return false ||
		(b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z')
}
