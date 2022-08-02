package pixiv_api

import "github.com/kiririx/krutils/httpx"

// apiClient api client
func (p *PixivClient) apiClient() *httpx.HttpClient {
	c := httpx.Client().Timeout(p.APITimeout)
	if p.ProxyURL != "" {
		c.Proxy(p.ProxyURL)
	}
	return c
}

// dwClient download client
func (p *PixivClient) dwClient() *httpx.HttpClient {
	c := httpx.Client().Timeout(p.DownloadTimeout)
	if p.ProxyURL != "" {
		c.Proxy(p.ProxyURL)
	}
	return c
}
