# Lang Middleware

The Lang Middleware provide a simple way to obtain the prefed language of the client.

## How to use

Initialize a new middleware using one of the function provided.

* `NewCookieOnly` returns a new LangMiddleware configured to use only informations from cookie named 
* `NewHeaderOnly` returns a new LangMiddleware configured to use only informations from accept-language header
* `NewHeaderAndCookie` returns a new LangMiddleware configured to use informations from both cookie and accept-language header

Then use the `Extractor()` function to get the http handler.


## Example
```golang
package main

import (
	"net/http"
	"github.com/go-chi/chi"
	"github.com/jderail/langmiddleware"
)

func main() {

	lang, _ := langmiddleware.NewCookieOnly("en", []string{"en", "fr"}, "lang-cookie")
	r := chi.NewRouter()
	r.Use(lang.Extractor())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Context().Value(langmiddleware.LangContextKey).(string)))
	})
	http.ListenAndServe(":3000", r)
}
```