package werich

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/russross/blackfriday"
)

func TestRun(t *testing.T) {
	input := []byte(`---
date: 2018-11-20 14:14
title: 架构
---

## 数据收集系统
* 政府公开数据
* 全网新闻
* 社交媒体

## 数据整合系统
* 数据清洗
* 语义识别
`)
	out := Run(input)
	t.Log(string(out))
	bf := blackfriday.New(blackfriday.WithNoExtensions())
	out2 := bf.Parse(input)
	t.Log(spew.Sdump(out2))
}
