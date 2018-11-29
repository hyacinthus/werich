// Package werich is a
// Blackfriday Markdown Processor WeApp rich-text Randerer
// Available at http://github.com/hyacinthus/werich
//
// Copyright © 2018 Muninn <hyacinthus@gmail.com>.
// Distributed MIT License.
// See README.md for details.
//
package werich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	bf "github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
)

// Renderer is a type that implements the Renderer interface for rich-text output.
type Renderer struct {
	// If you define this as "pre_", the tag class will seems like "pre_md_xxx",
	// don't forget "_". The default is empty, clsss will be "md_xxx".
	CSSPrefix string
	// Heading level offset, if it is 2, <h2> will change to <h4>
	HeadingOffset int
}

func isValidLink(link []byte) bool {
	// current directory : begin with "./"
	if bytes.HasPrefix(link, []byte("http")) {
		return true
	}
	return false
}

func (r *Renderer) tag(w io.Writer, name string) {
	_, err := fmt.Fprintf(w, `{"name":"%s"}`, name)
	if err != nil {
		logrus.WithError(err).Error("write tag failed")
	}
}

func (r *Renderer) text(w io.Writer, text []byte) {
	escaped, err := json.Marshal(string(text))
	if err != nil {
		logrus.WithError(err).Error("json marshal text failed")
	}
	_, err = fmt.Fprintf(w, `{"type":"text","text":"%s"}`, escaped[1:len(escaped)-1])
	if err != nil {
		logrus.WithError(err).Error("write text failed")
	}
}

func (r *Renderer) start(w io.Writer, tag string) {
	_, err := fmt.Fprintf(w, `{"name":"%s","attrs":{"class":"%smd_%s"},"children":[`, tag, r.CSSPrefix, tag)
	if err != nil {
		logrus.WithError(err).Error("write tag start failed")
	}
}

// the front part of a tag, if class is empty, default be prefix_md_tagname
func (r *Renderer) startWithClass(w io.Writer, tag, class string) {
	if class == "" {
		class = "md_" + tag
	}
	_, err := fmt.Fprintf(w, `{"name":"%s","attrs":{"class":"%s%s"},"children":[`, tag, r.CSSPrefix, class)
	if err != nil {
		logrus.WithError(err).Error("write tag start failed")
	}
}

func (r *Renderer) end(w io.Writer) {
	_, err := fmt.Fprint(w, `]}`)
	if err != nil {
		logrus.WithError(err).Error("write tag end failed")
	}
}

func (r *Renderer) empty(w io.Writer) {
	_, err := fmt.Fprint(w, `{}`)
	if err != nil {
		logrus.WithError(err).Error("write empty failed")
	}
}

func headingTag(level int) string {
	real := level
	if level <= 1 {
		real = 1
	}
	if level >= 6 {
		real = 6
	}
	return fmt.Sprintf("h%d", real)
}

// RenderNode is a renderer of a single node of a markdown syntax tree. For
// block nodes it will be called twice: first time with entering=true, second
// time with entering=false, so that it could know when it's working on an open
// tag and when on close. It writes the result to w.
//
// The return value is a way to tell the calling walker to adjust its walk
// pattern: e.g. it can terminate the traversal by returning Terminate. Or it
// can ask the walker to skip a subtree of this node by returning SkipChildren.
// The typical behavior is to return GoToNext, which asks for the usual
// traversal to the next node.
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	if entering && node.Prev != nil {
		// json list sep, if you want skip some node, write {} to w, for right json syntax.
		w.Write([]byte(","))
	}
	switch node.Type {
	case bf.Text:
		r.text(w, node.Literal)
	case bf.Softbreak:
		break
	case bf.Hardbreak:
		r.tag(w, "br")
	case bf.Emph:
		if entering {
			r.start(w, "em")
		} else {
			r.end(w)
		}
	case bf.Strong:
		if entering {
			r.start(w, "strong")
		} else {
			r.end(w)
		}
	case bf.Del:
		if entering {
			r.start(w, "del")
		} else {
			r.end(w)
		}
	case bf.HTMLSpan:
		// can not support html code
		r.empty(w)
		return bf.SkipChildren
	case bf.Link:
		// just output the text in weapp, mark link with a css class
		if entering {
			r.start(w, "a")
		} else {
			r.end(w)
		}
	case bf.Image:
		if !isValidLink(node.LinkData.Destination) {
			// image always have a alt text children
			// break this node will bring alt text up a level
			break
		}
		if entering {
			_, err := fmt.Fprintf(w, `{"name":"img","attrs":{"class":"%smd_img","src":"%s","alt":"%s"}}`,
				r.CSSPrefix, node.LinkData.Destination, node.LinkData.Title)
			if err != nil {
				logrus.WithError(err).Error("write tag start failed")
			}
		}
		return bf.SkipChildren
	case bf.Code:
		r.start(w, "code")
		r.text(w, node.Literal)
		r.end(w)
	case bf.Document:
		break
	case bf.Paragraph:
		if entering {
			r.start(w, "p")
		} else {
			r.end(w)
		}
	case bf.BlockQuote:
		if entering {
			r.start(w, "blockquote")
		} else {
			r.end(w)
		}
	case bf.HTMLBlock:
		//  原样输出 html 片段
		r.text(w, node.Literal)
	case bf.Heading:
		headingLevel := node.HeadingData.Level + r.HeadingOffset
		tag := headingTag(headingLevel)
		if entering {
			class := "md_" + tag
			if node.IsTitleblock {
				class = "md_title"
			}
			r.startWithClass(w, tag, class)
		} else {
			r.end(w)
		}
	case bf.HorizontalRule:
		r.tag(w, "hr")
	case bf.List:
		tag := "ul"
		if node.ListFlags&bf.ListTypeOrdered != 0 {
			tag = "ol"
		}
		if node.ListFlags&bf.ListTypeDefinition != 0 {
			tag = "dl"
		}
		if entering {
			r.start(w, tag)
		} else {
			r.end(w)
		}
	case bf.Item:
		tag := "li"
		if node.ListFlags&bf.ListTypeDefinition != 0 {
			tag = "dd"
		}
		if node.ListFlags&bf.ListTypeTerm != 0 {
			tag = "dt"
		}
		if entering {
			r.start(w, tag)
		} else {
			r.end(w)
		}
	case bf.CodeBlock:
		if len(node.Info) > 0 {
			// for prismjs, class must be language-xxx
			_, err := fmt.Fprintf(w, `{"name":"code","attrs":{"class":"language-%s"},"children":[`, node.Info)
			if err != nil {
				logrus.WithError(err).Error("write tag start failed")
			}
		} else {
			r.start(w, "code")
		}
		lines := bytes.Split(node.Literal, []byte("\n"))
		for i, line := range lines {
			if i != 0 {
				w.Write([]byte(","))
			}
			r.text(w, line)
		}
		r.end(w)
	case bf.Table:
		if entering {
			r.start(w, "table")
		} else {
			r.end(w)
		}
	case bf.TableCell:
		tag := "td"
		if node.IsHeader {
			tag = "th"
		}
		if entering {
			r.start(w, tag)
		} else {
			r.end(w)
		}
	case bf.TableHead:
		if entering {
			r.start(w, "thead")
		} else {
			r.end(w)
		}
	case bf.TableBody:
		if entering {
			r.start(w, "tbody")
		} else {
			r.end(w)
		}
	case bf.TableRow:
		if entering {
			r.start(w, "tr")
		} else {
			r.end(w)
		}
	default:
		logrus.Errorf("Unknown node type: %s", node.Type.String())
	}
	return bf.GoToNext
}

// RenderHeader writes rich text nodes front [
func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
	_, err := fmt.Fprint(w, "[")
	if err != nil {
		logrus.WithError(err).Error("write start failed")
	}
}

// RenderFooter writes rich text nodes latter ]
func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
	_, err := fmt.Fprint(w, "]")
	if err != nil {
		logrus.WithError(err).Error("write end failed")
	}
}
