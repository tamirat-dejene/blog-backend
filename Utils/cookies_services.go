package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieOptions struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// SetCookie sets a cookie on the Gin context with custom options
func SetCookie(c *gin.Context, opts CookieOptions) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     opts.Name,
		Value:    opts.Value,
		Path:     opts.Path,
		Domain:   opts.Domain,
		MaxAge:   opts.MaxAge,
		Expires:  time.Now().Add(time.Duration(opts.MaxAge) * time.Second),
		Secure:   opts.Secure,
		HttpOnly: opts.HttpOnly,
		SameSite: opts.SameSite,
	})
}

// GetCookie retrieves a cookie from the Gin context by name
func GetCookie(c *gin.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

// DeleteCookie deletes a cookie from the Gin context
func DeleteCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}
