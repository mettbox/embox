package middleware

import "github.com/gin-gonic/gin"

// LanguageMiddleware sets the language preference based on the Accept-Language header.
// If the header is not set, it defaults to "en".
// Usage: lang, _ := c.Get("lang")
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		c.Set("lang", lang)
		c.Next()
	}
}
