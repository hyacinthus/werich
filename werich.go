package werich

import (
	"bytes"

	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

// Delimiter
var yamldelim = []byte("---")

// bf instance
var bfhtml, bfrich *blackfriday.Markdown

func init() {
	bfhtml = blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	bfrich = blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
}

// MD Markdown entity
type MD struct {
	// meta data front matter
	meta []byte
	// source without front matter
	body []byte
	// may be more config
}

// NewMD make a MD
func NewMD(src []byte) *MD {
	var md = new(MD)
	if bytes.HasPrefix(src, yamldelim) {
		parts := bytes.SplitN(src, yamldelim, 3)
		md.meta = parts[1]
		md.body = parts[2]
	} else {
		md.body = src
	}
	return md
}

// Meta unmarshal meta to struct or map v
func (md *MD) Meta(v interface{}) error {
	return yaml.Unmarshal(md.meta, v)
}

// HTML convert md to html
func (md *MD) HTML() []byte {
	return bfhtml.Run(md.body)
}

// Rich render markdown to weapp rich-text json struct
