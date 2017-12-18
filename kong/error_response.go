package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ErrorFromResponse(response *http.Response, errorResponse map[string]interface{}) error {
	bytes, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("unexpected status (%v) received: %v", response.Status, errorResponse)
	} else {
		return fmt.Errorf("unexpected status (%v) received: %v", response.Status, string(bytes))
	}
}
