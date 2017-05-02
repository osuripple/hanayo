package doc

import "io/ioutil"

const referenceLanguage = "en"

var docFiles []Document

// File represents the single documentation file of a determined language.
type File struct {
	IsUpdated      bool
	Title          string
	referencesFile string
}

// Data retrieves data from file's actual file on disk.
func (f File) Data() (string, error) {
	data, err := ioutil.ReadFile(f.referencesFile)
	return string(data), err
}

// Document represents a documentation file, providing its old ID, its slug,
// and all its variations in the various languages.
type Document struct {
	Slug      string
	OldID     int
	Languages map[string]File
}

// File retrieves a Document's File based on the passed language, and returns
// the values for the referenceLanguage (en) if in the passed language they are
// not available
func (d Document) File(lang string) File {
	if vals, ok := d.Languages[lang]; ok {
		return vals
	}
	return d.Languages[referenceLanguage]
}

// LanguageDoc has the only purpose to be returned by GetDocs.
type LanguageDoc struct {
	Title string
	Slug  string
}

// GetDocs retrieves a list of documents in a certain language, with titles and
// slugs.
func GetDocs(lang string) []LanguageDoc {
	var docs []LanguageDoc

	for _, file := range docFiles {
		docs = append(docs, LanguageDoc{
			Slug:  file.Slug,
			Title: file.File(lang).Title,
		})
	}

	return docs
}

// SlugFromOldID gets a doc file's slug from its old ID
func SlugFromOldID(i int) string {
	for _, d := range docFiles {
		if d.OldID == i {
			return d.Slug
		}
	}

	return ""
}

// GetFile retrieves a file, given a slug and a language.
func GetFile(slug, language string) File {
	for _, f := range docFiles {
		if f.Slug != slug {
			continue
		}
		if val, ok := f.Languages[language]; ok {
			return val
		}
		return f.Languages[referenceLanguage]
	}
	return File{}
}
