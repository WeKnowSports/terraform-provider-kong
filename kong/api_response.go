package kong

import (
	"encoding/json"
)

// APIResponse : Kong API response object structure
type APIResponse struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name"`
	UpstreamURL            string `json:"upstream_url"`
	StripURI               bool   `json:"strip_uri"`
	PreserveHost           bool   `json:"preserve_host"`
	Retries                int    `json:"retries,omitempty"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int    `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int    `json:"upstream_read_timeout,omitempty"`
	HTTPSOnly              bool   `json:"https_only"`
	HTTPIfTerminated       bool   `json:"http_if_terminated"`

	// These will be set in our UnmarshalJSON function.
	Hosts   []string `json:"-"`
	Methods []string `json:"-"`
	Uris    []string `json:"-"`
}

func (r *APIResponse) UnmarshalJSON(data []byte) error {
	type Alias APIResponse
	wrapped := &struct {
		RawHosts   interface{} `json:"hosts,omitempty"`
		RawMethods interface{} `json:"methods,omitempty"`
		RawUris    interface{} `json:"uris,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	if hosts, ok := wrapped.RawHosts.([]interface{}); ok {
		for _, host := range hosts {
			r.Hosts = append(r.Hosts, host.(string))
		}
	}

	if methods, ok := wrapped.RawMethods.([]interface{}); ok {
		for _, method := range methods {
			r.Methods = append(r.Methods, method.(string))
		}
	}

	if uris, ok := wrapped.RawUris.([]interface{}); ok {
		for _, uri := range uris {
			r.Uris = append(r.Uris, uri.(string))
		}
	}

	return nil
}
