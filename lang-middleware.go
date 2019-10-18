package langmiddleware

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/text/language"
)

type contextKey struct {
	name string
}

// LangContextKey is the key used to store Lang in context
var LangContextKey = &contextKey{"lang"}

// LangSource can be used to select possible source of language
type LangSource int

const (
	// Cookie indicates that only cookie will be used to extract language informations
	Cookie LangSource = 1 << iota
	// Header indicates that only header will be used to extract language informations
	Header
	// HeaderAndCookie indicates that both header and cookie will be used to extract language informations
	HeaderAndCookie = Cookie | Header
)

// LangMiddleware
type LangMiddleware struct {
	CookieName         string
	SupportedLanguages []string
	DefaultLanguage    string
	source             LangSource
}

var (
	// ErrHeaderParsing is returned by ParseLangHeader when accept language header is not conform
	ErrHeaderParsing = errors.New("An error occurred while parsing accept-language header")
	// ErrEmptyHeader is returned by ParseLangHeader when no tags are provided in header
	ErrEmptyHeader = errors.New("Accept-language header is empty")
)

// globalInit initialize LangMiddleware with common properties
func globalInit(pDefaultLang string, supportedLanguages []string) (*LangMiddleware, error) {
	defaultLang, err := language.ParseBase(pDefaultLang)
	if err != nil {
		return nil, nil
	}

	for i := 0; i < len(supportedLanguages); i++ {
		tmp, err := language.ParseBase(supportedLanguages[i])
		if err != nil {
			return nil, nil
		}
		supportedLanguages[i] = tmp.String()
	}

	return &LangMiddleware{
		SupportedLanguages: supportedLanguages,
		DefaultLanguage:    defaultLang.String(),
	}, nil
}

// NewCookieOnly returns a new LangMiddleware with a configuration to use only informations from cookie named
func NewCookieOnly(defaultLang string, supportedLanguages []string, cookieName string) (*LangMiddleware, error) {
	mid, err := globalInit(defaultLang, supportedLanguages)
	if err != nil {
		return nil, err
	}
	mid.source = Cookie
	mid.CookieName = cookieName
	return mid, nil
}

// NewHeaderOnly returns a new LangMiddleware with a configuration to use only informations from accept-language header
func NewHeaderOnly(defaultLang string, supportedLanguages []string) (*LangMiddleware, error) {
	mid, err := globalInit(defaultLang, supportedLanguages)
	if err != nil {
		return nil, err
	}
	mid.source = Header
	return mid, nil
}

// NewCookieAndHeader returns a new LangMiddleware with a configuration to use informations from both cookie and accept-language header
func NewCookieAndHeader(defaultLang string, supportedLanguages []string, cookieName string) (*LangMiddleware, error) {
	mid, err := globalInit(defaultLang, supportedLanguages)
	if err != nil {
		return nil, err
	}
	mid.source = HeaderAndCookie
	mid.CookieName = cookieName
	return mid, nil
}

// Extractor return http handler to use as middleware
func (i *LangMiddleware) Extractor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			lang := i.extract(r)

			ctx = context.WithValue(ctx, LangContextKey, lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func (i *LangMiddleware) extract(r *http.Request) string {
	if i.source&Cookie != 0 {
		c, err := r.Cookie(i.CookieName)
		if err == nil {
			return i.fromCookie(c)
		}
	}
	if i.source&Header != 0 {
		return i.fromLangHeader(r.Header.Get("Accept-Language"))
	}
	return i.DefaultLanguage
}

func (i *LangMiddleware) fromLangHeader(header string) string {
	tags, priorities, err := language.ParseAcceptLanguage(header)
	if err != nil {
		return i.DefaultLanguage
	}

	if len(tags) == 0 {
		return i.DefaultLanguage
	}

	for j := 0; j < len(tags) && j < len(priorities); j++ {
		tmp, _ := tags[j].Base()
		asStr := tmp.String()
		if asStr == "mul" {
			return i.DefaultLanguage
		} else if contains(i.SupportedLanguages, asStr) {
			return asStr
		}
	}

	return i.DefaultLanguage
}

func (i *LangMiddleware) fromCookie(c *http.Cookie) string {
	var b language.Base
	b, err := language.ParseBase(c.Value)
	if err != nil {
		return i.DefaultLanguage
	}
	asStr := b.String()
	if contains(i.SupportedLanguages, asStr) {
		return asStr
	}
	return i.DefaultLanguage
}

func contains(slice []string, search string) bool {
	for _, tmp := range slice {
		if tmp == search {
			return true
		}
	}
	return false
}
