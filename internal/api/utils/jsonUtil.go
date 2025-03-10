package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func ParseJsonRequest(c *gin.Context) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})

	bodyAsByteArray, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyAsByteArray, &jsonMap)

	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
