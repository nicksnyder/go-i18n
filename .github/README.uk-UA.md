# go-i18n
![Build status](https://github.com/nicksnyder/go-i18n/workflows/Build/badge.svg) [![Report card](https://goreportcard.com/badge/github.com/nicksnyder/go-i18n/v2)](https://goreportcard.com/report/github.com/nicksnyder/go-i18n/v2) [![codecov](https://codecov.io/gh/nicksnyder/go-i18n/graph/badge.svg?token=A9aMfR9vxG)](https://codecov.io/gh/nicksnyder/go-i18n) [![Sourcegraph](https://sourcegraph.com/github.com/nicksnyder/go-i18n/-/badge.svg)](https://sourcegraph.com/github.com/nicksnyder/go-i18n?badge)

go-i18n — це Go [пакет](#package-i18n) та [інструмент](#command-goi18n), які допомагають перекладати Go програми на різні мови.

- Підтримує [множинні форми](http://cldr.unicode.org/index/cldr-spec/plural-rules) для всіх 200+ мов у [Unicode Common Locale Data Repository (CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html).
  - Код і тести [автоматично генеруються](https://github.com/nicksnyder/go-i18n/tree/main/internal/plural/codegen) з даних [CLDR](http://cldr.unicode.org/index/downloads).
- Підтримує рядки з іменованими змінними, використовуючи синтаксис [text/template](http://golang.org/pkg/text/template/).
- Підтримує файли повідомлень у будь-якому форматі (наприклад, JSON, TOML, YAML).

## Пакет i18n

[![Go Reference](https://pkg.go.dev/badge/github.com/nicksnyder/go-i18n/v2/i18n.svg)](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2/i18n)

Пакет i18n забезпечує підтримку пошуку повідомлень відповідно до набору мовних уподобань.

```go
import "github.com/nicksnyder/go-i18n/v2/i18n"
```

Створіть Bundle, який використовуватимете протягом усього терміну служби вашої програми.

```go
bundle := i18n.NewBundle(language.English)
```

Завантажуйте переклади у ваш пакет під час ініціалізації.

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

```go
// Якщо використовуєте go:embed
//go:embed locale.*.toml
var LocaleFS embed.FS

bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFileFS(LocaleFS, "locale.es.toml")
```

Створіть Localizer, який використовуватимете для набору мовних уподобань.

```go
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

Використовуйте Localizer для пошуку повідомлень.

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

## Команда goi18n

[![Go Reference](https://pkg.go.dev/badge/github.com/nicksnyder/go-i18n/v2/goi18n.svg)](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2/goi18n)

Команда goi18n управляє файлами повідомлень, що використовуються пакетом i18n.

```
go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n -help
```

### Витяг повідомлень

Використовуйте команду `goi18n extract`, щоб витягнути всі літерали структури i18n.Message із Go-файлів у файл повідомлень для перекладу.

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### Переклад нової мови

1. Створіть порожній файл повідомлень для мови, яку ви хочете додати (наприклад, translate.uk.toml).
2. Виконайте команду `goi18n merge active.en.toml translate.es.toml`, щоб заповнити `translate.es.toml` повідомленнями для перекладу.

   ```toml
   # translate.uk.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hello {{.Name}}"
   ```

3. Після перекладу файлу `translate.es.toml` перейменуйте його на `active.es.toml`.

   ```toml
   # active.uk.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Вітаю {{.Name}}"
   ```

4. Завантажте файл `active.es.toml` у свій пакет.

   ```go
   bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
   bundle.LoadMessageFile("active.es.toml")
   ```

### Переклад нових повідомлень

Якщо ви додали нові повідомлення до своєї програми:

1.	Виконайте `goi18n extract`, щоб оновити файл `active.en.toml` новими повідомленнями.
2.	Виконайте `goi18n merge active.*.toml`, щоб згенерувати оновлені файли `translate.*.toml`.
3.	Перекладіть усі повідомлення у файлах `translate.*.toml`.
4.	Виконайте `goi18n merge active.*.toml translate.*.toml`, щоб об’єднати перекладені повідомлення з активними файлами повідомлень.

## Для отримання додаткової інформації та прикладів:

- Ознайомтеся з [документацією](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2).
- Подивіться [приклади коду](https://github.com/nicksnyder/go-i18n/blob/main/i18n/example_test.go) та [тести](https://github.com/nicksnyder/go-i18n/blob/main/i18n/localizer_test.go).
- Перегляньте приклад [додатку](https://github.com/nicksnyder/go-i18n/tree/main/example).

## Переклади цього документа

Переклади цього документа, зроблені спільнотою, можна знайти в папці [.github](.github).

Ці переклади підтримуються спільнотою і не підтримуються автором цього проєкту.  
Немає гарантії, що вони є точними або актуальними.

## Ліцензія

go-i18n доступний під ліцензією MIT. Див. файл [LICENSE](LICENSE) для отримання додаткової інформації.
