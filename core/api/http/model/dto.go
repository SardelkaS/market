package api_http_model

import "core/internal/common"

type VersionResponse struct {
	Version string `json:"version"`
	common.Response
}
