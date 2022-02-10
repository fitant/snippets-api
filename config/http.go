package config

import "fmt"

type HTTPServerConfig struct {
	host         string
	port         int
	CORS         string
	baseURL      string
	BaseEndpoint string
}

func (h HTTPServerConfig) GetListenAddr() string {
	return fmt.Sprintf("%s:%d", h.host, h.port)
}

func (h HTTPServerConfig) GetBaseURL() string {
	return fmt.Sprintf("http://%s:%d%s/%s", h.host, h.port, h.BaseEndpoint, h.baseURL)
}
