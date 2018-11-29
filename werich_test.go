package werich

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestRun(t *testing.T) {
	input, err := ioutil.ReadFile("test.md")
	if err != nil {
		t.Error(err)
	}
	md := NewMD(input)
	t.Log(string(md.HTML()))
	rich := md.Rich()
	t.Log(string(rich))
	var pretty bytes.Buffer
	err = json.Indent(&pretty, rich, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(pretty.String())
}
