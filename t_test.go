package pixiv_api

import "testing"

func TestClient(t *testing.T) {
	GetPixivClient(PixivClient{
		ProxyURL:        "",
		RefreshToken:    "",
		accessToken:     "",
		expireTime:      0,
		DownloadTimeout: 0,
		APITimeout:      0,
		Login:           false,
	})
}
