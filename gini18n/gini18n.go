package gini18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ContextKey default context key
var ContextKey = "WW_LOCALE"

const (
	// ParamName default parameter key name
	ParamName = "__locale"

	// HeaderName default http header name
	HeaderName = "X-Accept-Language"

	// CookieName default cookie name
	CookieName = "WW_LOCALE"
)

// Localizer localizer middleware
type Localizer struct {
	Locales []string

	ParamName          string
	HeaderName         string
	CookieName         string
	FromAcceptLanguage bool
}

// Default create a default localizer
// = NewLocalizer()
func Default() *Localizer {
	return NewLocalizer()
}

// NewLocalizer create a default Localizer
func NewLocalizer(locales ...string) *Localizer {
	if len(locales) == 0 {
		locales = []string{"en"}
	}

	return &Localizer{
		Locales:            locales,
		HeaderName:         HeaderName,
		ParamName:          ParamName,
		CookieName:         CookieName,
		FromAcceptLanguage: true,
	}
}

// GetLocale get locale from gin.Context
func GetLocale(c *gin.Context) string {
	if v, ok := c.Get(ContextKey); ok {
		return v.(string)
	}
	return ""
}

// SetLocale set locale to gin.Context
func SetLocale(c *gin.Context, loc string) {
	c.Set(ContextKey, loc)
}

// Handler returns the gin.HandlerFunc
func (ll *Localizer) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ll.handle(c)
	}
}

// handle process gin request
func (ll *Localizer) handle(c *gin.Context) {
	loc := ""

	if ll.ParamName != "" {
		loc = ll.getLocaleFromParameter(c, ll.ParamName)
	}

	if loc == "" && ll.HeaderName != "" {
		loc = ll.getLocaleFromHeader(c, ll.HeaderName)
	}

	if loc == "" && ll.CookieName != "" {
		loc = ll.getLocaleFromCookie(c, ll.CookieName)
	}

	if loc == "" && ll.FromAcceptLanguage {
		loc = ll.getLocaleFromHeader(c, "Accept-Language")
	}

	if loc == "" {
		loc = ll.Locales[0]
	}

	SetLocale(c, loc)

	c.Next()
}

func (ll *Localizer) getLocaleFromHeader(c *gin.Context, k string) string {
	loc := c.GetHeader(k)
	qls := strings.FieldsFunc(loc, func(r rune) bool {
		return strings.ContainsRune(",; ", r)
	})
	for _, ql := range qls {
		if ll.acceptable(ql) {
			return ql
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromParameter(c *gin.Context, k string) string {
	if loc, ok := c.GetPostForm(k); ok {
		if ll.acceptable(loc) {
			return loc
		}
	}
	if loc, ok := c.GetQuery(k); ok {
		if ll.acceptable(loc) {
			return loc
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromCookie(c *gin.Context, k string) string {
	if loc, err := c.Cookie(k); err == nil {
		if ll.acceptable(loc) {
			return loc
		}
	}
	return ""
}

func (ll *Localizer) acceptable(loc string) bool {
	if loc != "" {
		for _, al := range ll.Locales {
			if strings.HasPrefix(loc, al) {
				return true
			}
		}
	}
	return false
}
