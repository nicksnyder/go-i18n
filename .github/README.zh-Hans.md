# go-i18n ![Build status](https://github.com/nicksnyder/go-i18n/workflows/Build/badge.svg) [![Report card](https://goreportcard.com/badge/github.com/nicksnyder/go-i18n)](https://goreportcard.com/report/github.com/nicksnyder/go-i18n) [![codecov](https://codecov.io/gh/nicksnyder/go-i18n/branch/master/graph/badge.svg)](https://codecov.io/gh/nicksnyder/go-i18n) [![Sourcegraph](https://sourcegraph.com/github.com/nicksnyder/go-i18n/-/badge.svg)](https://sourcegraph.com/github.com/nicksnyder/go-i18n?badge)

go-i18n 是一个帮助您将 Go 程序翻译成多种语言的 Go [包](#package-i18n)和[命令](#command-goi18n)。

- 支持 [Unicode Common Locale Data Repository (CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html)
  中所有 200 多种语言的[复数字符串](http://cldr.unicode.org/index/cldr-spec/plural-rules)。
  - 代码和测试是基于 [CLDR 数据](http://cldr.unicode.org/index/downloads)[自动生成](https://github.com/nicksnyder/go-i18n/tree/main/internal/plural/codegen)的。
- 使用 [text/template](http://golang.org/pkg/text/template/) 语法支持带有命名变量的字符串。
- 支持所有格式的消息文件（例如：JSON、TOML、YAML）。

<strong align="center">
<samp>

[**English**](../README.md) · [**简体中文**](README.zh-Hans.md)

</samp>
</strong>

## i18n 包

[![GoDoc](https://pkg.go.dev/github.com/nicksnyder/go-i18n?status.svg)](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2/i18n)

i18n 包支持根据一组语言环境首选项来查找消息。

```go
import "github.com/nicksnyder/go-i18n/v2/i18n"
```

创建一个 Bundle 以在应用程序的整个生命周期中使用。

```go
bundle := i18n.NewBundle(language.English)
```

在初始化时，将翻译加载到你的 Bundle 中。

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

```go
// 如果使用 go:embed
//go:embed locale.*.toml
var LocaleFS embed.FS

bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFileFS(LocaleFS, "locale.es.toml")
```

创建一个 Localizer 以便用于一组首选语言。

```go
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

使用此 Localizer 查找消息。

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
}) // Nick 有两只猫
```

## goi18n 命令

[![GoDoc](https://pkg.go.dev/github.com/nicksnyder/go-i18n?status.svg)](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2/goi18n)

goi18n 命令管理 i18n 包所使用的消息文件。

```
go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n -help
```

### 提取消息

使用 `goi18n extract` 将 Go 源文件中的所有 i18n.Message 结构中的文字提取到消息文件中以进行翻译。

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### 翻译一种新语言

1. 为你要添加的语言创建一个空的消息文件（例如：`translate.es.toml`）。
2. 运行 `goi18n merge active.en.toml translate.es.toml` 以将要翻译的消息填充到 `translate.es.toml` 中。

   ```toml
   # translate.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hello {{.Name}}"
   ```

3. 完成 `translate.es.toml` 的翻译之后，将其重命名为 `active.es.toml`。

   ```toml
   # active.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hola {{.Name}}"
   ```

4. 加载 `active.es.toml` 到你的 Bundle 中。

   ```go
   bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
   bundle.LoadMessageFile("active.es.toml")
   ```

### 翻译新消息

如果你在程序中添加了新消息：

1. 运行 `goi18n extract` 以将新的消息更新到 `active.en.toml`。
2. 运行 `goi18n merge active.*.toml` 以生成更新后的 `translate.*.toml` 文件。
3. 翻译 `translate.*.toml` 文件中的所有消息。
4. 运行 `goi18n merge active.*.toml translate.*.toml` 将翻译后的消息合并到活跃消息文件
   （Active Message Files）中。

## 进一步的信息和示例：

- 阅读[文档](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2)。
- 查看[代码示例](https://github.com/nicksnyder/go-i18n/blob/main/i18n/example_test.go)和
  [测试](https://github.com/nicksnyder/go-i18n/blob/main/i18n/localizer_test.go)。
- 查看示例[程序](https://github.com/nicksnyder/go-i18n/tree/main/example)。

## 许可证

go-i18n 使用在 MIT 许可来提供。更多的相关信息，请参 [LICENSE](LICENSE) 文件。
