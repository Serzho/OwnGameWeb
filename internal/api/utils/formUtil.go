package utils

import "github.com/gin-gonic/gin"

func ParseFormRequest(c *gin.Context, keys []string) map[string]interface{} {
	var result map[string]interface{}

	for _, key := range keys {
		result[key] = c.PostForm(key)
	}

	return result
}
