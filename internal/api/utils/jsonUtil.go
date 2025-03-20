package utils

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func ParseJSONRequest(c *gin.Context) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})

	bodyAsByteArray, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, ErrReadingRequestBody
	}

	err = json.Unmarshal(bodyAsByteArray, &jsonMap)
	if err != nil {
		return nil, ErrUnmarshalJSON
	}

	return jsonMap, nil
}
