package response

type NetworkTokenSupport struct {
	Network string   `json:"network"`
	Tokens  []string `json:"tokens"`
}

type SupportedAssetsResponse struct {
	Supports []NetworkTokenSupport `json:"supports"`
}
