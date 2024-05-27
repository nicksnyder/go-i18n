# Web Example

This directory contains an example project that uses go-i18n.

```
go run main.go
```

Then open http://localhost:8080 in your web browser.

You can customize the template data and locale via query parameters like this:
http://localhost:8080/?name=Nick&unreadEmailCount=2
http://localhost:8080/?name=Nick&unreadEmailCount=2&lang=es


# CLI Example

This CLI example shows how to use a `T()` function to translate simple strings in a compiled CLI application.

The [go-locale](github.com/Xuanwo/go-locale) library is used to detect the users preferred language.
The [cobra](github.com/spf13/cobr) library is used to setup the CLI itself.
the [embed](https://pkg.go.dev/embed) library is used to compile the translation json files into the CLI binary.
These libraries may need to be installed specifically to compile this cli application if needed.

```bash
go get github.com/spf13/cobra
go get github.com/Xuanwo/go-locale
```

To run the example for a given language, you can do the following:
```bash
$> LANGUAGE=en-US go run cli.go
Detected Locale is en-US
Hello World!

$> LANGUAGE=es-ES go run cli.go
La configuración regional detectada es es-ES
Hola Mundo!
```

To extract the string from this example:
>Assumes you are in the `example` directory.

```bash
goi18n extract -outdir ./ -format json --sourceLanguage en-US ./cli.go
```