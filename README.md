# Markdown 转微信小程序 rich-text

## 使用说明

用 Markdown 直接生成微信小程序富文本，传给前端后前端赋值到 `<rich-text>` 的 nodes 即可。

示例：

```go
package main

import (
    "fmt"
    "io/ioutil"

    "github.com/hyacinthus/werich"
)

func main() {
    input, err := ioutil.ReadFile("test.md")
    if err != nil {
        panic(err)
    }
    md := werich.NewMD(input)
    // rich
    rich := md.Rich()
    fmt.Println(string(rich))
}
```

## 自定义 CSS

生成的富文本节点 css 属性为 md_xxx ，我们提供了一份样例，复制到小程序在相应页面导入即可。
你可以修改其中的样式。如果对 class 名称不满意，我们也提供了添加前缀的功能，请参考 `werich.go` 文件，
自行使用 renderer 和 blackfriday 解析即可。

## 解析 Front matter

有很多写作软件和静态网站生成器在 Markdown 头部添加了元数据，我们支持解析 YAML 型的元数据，
别的类型您可以参考 `werich.go` 文件，自行编写。

```go
package main

import (
    "fmt"
    "io/ioutil"

    "github.com/hyacinthus/werich"
)

func main() {
    input, err := ioutil.ReadFile("test.md")
    if err != nil {
        panic(err)
    }
    md := werich.NewMD(input)
    // meta
    var meta = make(map[string]interface{})
    err = md.Meta(meta)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v",meta)
}
```

## Markdown 语法限制

为了能在小程序和网页完美显示，我们对 Markdown 规则做了如下限制：

- 链接在小程序中将只保留描述文字，在网页中正常，尽量避免使用链接。
- 图片只能使用我们域名下的链接或者 "./xxx" 的形式，直接拖到编辑器就会是后者的效果。
- 图片不支持 title 语法，不支持 alt 语法，就光写地址就好了。
- 不支持直接使用 html 语法
- 列表最多只支持两层
- 还没有支持代码高亮功能，最好不要在文章中包含代码。
