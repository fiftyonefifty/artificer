package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

const successTemplate = `
<html>
    <head>
        <meta charset="utf-8">
        <meta
            name="viewport"
            content="width=device-width, initial-scale=1"
        >
        <title>Artificer</title>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.7.2/css/bulma.min.css">
        <script
            defer
            src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"
        ></script>
    </head>
    <body>
        <section class="hero is-light is-bold is-fullheight">
            <div class="hero-head"></div>
            <div class="hero-body">
                <div class="container has-text-centered">
                    <h1 class="title">Artificer</h1>
                    <h2 class="subtitle">A crafter of Tokens</h2>
                    <p><a href="/health">Health</a>.</p>
                    <p><a href="/.well-known/openid-configuration">OpenId Configuration</a>.</p>
                </div>
               
            </div>
        </section>
    </body>
</html>`

// HealthCheck - Healthcheck Handler
func Index(c echo.Context) error {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	t, err := template.New("webpage").Parse(successTemplate)
	check(err)

	data := struct {
		Title string
	}{
		Title: "Success",
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, data)
	return c.HTMLBlob(http.StatusOK, buf.Bytes())

	return err
}
