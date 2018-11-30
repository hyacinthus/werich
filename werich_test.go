package werich

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

type metaData struct {
	ID     string
	Date   string
	Title  string
	Author string
	Tags   []string
}

func TestRun(t *testing.T) {
	input, err := ioutil.ReadFile("test.md")
	if err != nil {
		t.Error(err)
	}
	md := NewMD(input)
	// meta
	var meta = new(metaData)
	assert.Nil(t, md.Meta(meta))
	assert.Equal(t, "", meta.ID)
	assert.Equal(t, "2018-11-20 14:14", meta.Date)
	assert.Equal(t, "WeRich Test", meta.Title)
	assert.Equal(t, "Muninn", meta.Author)
	assert.Equal(t, []string{"golang", "wechat"}, meta.Tags)
	// html
	div := md.HTML()
	// t.Log(string(div))
	_, err = html.Parse(bytes.NewReader(div))
	assert.Nil(t, err)
	// rich
	rich := md.Rich()
	// t.Log(string(rich))
	dst, err := ioutil.ReadFile("test.txt")
	assert.Equal(t, rich, dst, "not expected")
	var pretty bytes.Buffer
	err = json.Indent(&pretty, rich, "", "    ")
	assert.Nil(t, err)
	// t.Log(pretty.String())
}
