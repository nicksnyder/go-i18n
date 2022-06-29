# go-i18n ![Build status](https://github.com/nicksnyder/go-i18n/workflows/Build/badge.svg) [![Report card](https://goreportcard.com/badge/github.com/nicksnyder/go-i18n)](https://goreportcard.com/report/github.com/nicksnyder/go-i18n) [![codecov](https://codecov.io/gh/nicksnyder/go-i18n/branch/master/graph/badge.svg)](https://codecov.io/gh/nicksnyder/go-i18n) [![Sourcegraph](https://sourcegraph.com/github.com/nicksnyder/go-i18n/-/badge.svg)](https://sourcegraph.com/github.com/nicksnyder/go-i18n?badge)

go-i18n 是一个帮助您将 Go 程序翻译成多种语言的 Go [包](#package-i18n) 和 [命令](#command-goi18n)。

- 支持 [Unicode Common Locale Data Repository (CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html) 中所有 200 多种语言的 [复数字符](http://cldr.unicode.org/index/cldr-spec/plural-rules)。
  - 代码和测试是从 [CLDR 数据](http://cldr.unicode.org/index/downloads) 中 [自动生成](https://github.com/nicksnyder/go-i18n/tree/main/v2/internal/plural/codegen) 的。
- 使用 [text/template](http://golang.org/pkg/text/template/) 语法支持带有命名变量的字符串。
- 支持任何格式的消息文件（例如：JSON、TOML、YAML）。

<strong align="center">
<samp>

[**English**](../README.md) · [**简体中文**](README.zh-Hans.md)

</samp>
</strong>

## Package i18n

[![GoDoc](https://godoc.org/github.com/nicksnyder/go-i18n?status.svg)](https://godoc.org/github.com/nicksnyder/go-i18n/v2/i18n)

i18n 包支持根据一组语言环境首选项查找消息。

```go
import "github.com/nicksnyder/go-i18n/v2/i18n"
```

创建一个 Bundle 以在应用程序的整个生命周期中使用。

```go
bundle := i18n.NewBundle(language.English)
```

在初始化期间将翻译加载到您的 bundle 中。

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

创建一个 Localizer 以用于一组语言首选项。

```go
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

使用 Localizer 查找消息。

```go
localizer.Localize(&i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID: "PersonCats",
        One: "{{.Name}} has {{.Count}} cat.",
        Other: "{{.Name}} has {{.Count}} cats.",
    },
    TemplateData: map[string]interface{}{
        "Name": "Nick",
        "Count": 2,
    },
    PluralCount: 2,
}) // Nick has 2 cats.
```

## goi18n 命令

[![GoDoc](https://godoc.org/github.com/nicksnyder/go-i18n?status.svg)](https://godoc.org/github.com/nicksnyder/go-i18n/v2/goi18n)

goi18n 命令管理 i18n 包使用的消息文件。

```
go get -u github.com/nicksnyder/go-i18n/v2/goi18n
goi18n -help
```

### 提取消息

使用 `goi18n extract` 将 Go 源文件中的所有 i18n.Message 结构文字提取到消息文件中以进行翻译。

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### 翻译一种新语言

1. 为您要添加的语言创建一个空消息文件（例如：`translate.es.toml`）。
2. 运行 `goi18n merge active.en.toml translate.es.toml` 以填充 `translate.es.toml` 要翻译的消息。

   ```toml
   # translate.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hello {{.Name}}"
   ```

3. 翻译完成 `translate.es.toml` 后，将其重命名为 `active.es.toml``。

   ```toml
   # active.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hola {{.Name}}"
   ```

4. 加载 `active.es.toml` 到您的 bundle 中。

   ```go
   bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
   bundle.LoadMessageFile("active.es.toml")
   ```

### 翻译新消息

如果您在程序中添加了新消息：

1. 运行 `goi18n extract` 以使用新消息更新 `active.en.toml`。
2. 运行 `goi18n merge active.*.toml` 以生成更新的 `translate.*.toml` 文件。
3. 翻译 `translate.*.toml` 文件中的所有消息。
4. 运行 `goi18n merge active.*.toml translate.*.toml` 将翻译后的消息合并到 active 消息文件中。

## 有关更多信息和示例：

- 阅读 [文档](https://godoc.org/github.com/nicksnyder/go-i18n/v2)。
- 查看 [代码示例](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go) 和 [测试](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/localizer_test.go)。
- 查看一个示例 [程序](https://github.com/nicksnyder/go-i18n/tree/main/v2/example)。

## 许可证

go-i18n 在 MIT 许可下可用。有关更多信息，请参阅 [许可证](LICENSE) 文件。
