package api_http_model

import "market_auth/internal/common"

type VersionResponse struct {
	Version string `json:"version"`
	common.Response
}
