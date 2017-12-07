package kong

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse map[string]interface {
}

func ErrorFromResponse(response *http.Response, errorResponse *ErrorResponse) error {
	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		return fmt.Errorf("unexpected status (%v) received: %v", response.Status, errorResponse)
	} else {
		return fmt.Errorf("unexpected status (%v) received: %v", response.Status, string(bytes))
	}
}
