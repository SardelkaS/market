package api_http_model

import "auth/internal/common"

type VersionResponse struct {
	Version string `json:"version"`
	common.Response
}
