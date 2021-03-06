package werich

import (
	"bytes"
	"io"

	bf "github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

// Delimiter
var yamldelim = []byte("---")

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

// Unix 将正文的换行符统一成unix样式
// blackfriday 现在有个bug，无法完美支持 Windows 换行符
func (md *MD) Unix() {
	md.body = bytes.Replace(md.body, []byte("\r\n"), []byte("\n"), -1)
	md.body = bytes.Replace(md.body, []byte("\r"), []byte("\n"), -1)
}

// Reader 给需要的函数提供全文 Reader
func (md *MD) Reader() io.Reader {
	var buffer bytes.Buffer
	if md.HasMeta() {
		buffer.Write(yamldelim)
		buffer.Write([]byte("\n"))
		buffer.Write(md.meta)
		buffer.Write([]byte("\n"))
		buffer.Write(yamldelim)
		buffer.Write([]byte("\n"))
	}
	buffer.Write(md.body)
	return &buffer
}

// Meta unmarshal meta to struct or map v
func (md *MD) Meta(v interface{}) error {
	return yaml.Unmarshal(md.meta, v)
}

// HasMeta 是否存在 meta data
func (md *MD) HasMeta() bool {
	return md.meta != nil && len(md.meta) != 0
}

// SetMeta 重新生成 meta 信息
func (md *MD) SetMeta(v interface{}) error {
	meta, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	md.meta = meta
	return nil
}

// HTML convert md to html
func (md *MD) HTML() []byte {
	return bf.Run(md.body, bf.WithExtensions(bf.CommonExtensions))
}

// Rich render markdown to weapp rich-text json struct
func (md *MD) Rich() []byte {
	renderer := &Renderer{}
	return bf.Run(md.body, bf.WithExtensions(bf.CommonExtensions), bf.WithRenderer(renderer))
}
