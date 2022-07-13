package pixiv_api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kiririx/krutils/algo_util"
	"github.com/kiririx/krutils/http_util"
	"github.com/kiririx/krutils/logx"
	"github.com/kiririx/krutils/str_util"
	"io"
	"os"
	"path"
	"regexp"
	"time"
)

var (
	clientId     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	clientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	hashSecret   = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
	apiHosts     = "https://app-api.pixiv.net"
)

// PixivClient pixiv客户端
//
// ProxyURL: 代理服务器地址
//
// APITimeout: API超时时间
//
// RefreshToken: 刷新令牌
//
// accessToken: 访问令牌
//
// expireTime: 访问令牌过期时间
//
// Login: 如果为false，则以游客身份访问，如果为true，会获取accessToken，以用户身份访问
type PixivClient struct {
	ProxyURL        string
	RefreshToken    string
	accessToken     string
	expireTime      int64
	DownloadTimeout time.Duration
	APITimeout      time.Duration
	Login           bool
}

func (p *PixivClient) getHeaders() (map[string]string, error) {
	localTime := time.Now().Format(time.RFC3339)
	px := make(map[string]string)
	px["Accept-Language"] = "en-us"
	px["X-Client-Time"] = localTime
	px["X-Client-Hash"] = genClientHash(localTime)
	px["User-Agent"] = "PixivAndroidApp/5.0.115 (Android 6.0)"
	if p.Login {
		err := p.Auth()
		if err != nil {
			return nil, err
		}
		px["Authorization"] = "Bearer " + p.accessToken
	}
	return px, nil
}

func (p *PixivClient) DownloadImg(url string) (string, error) {
	ext, _ := GetFileExt(url)
	fileName := algo_util.MD5(url) + "." + ext
	_, err := os.Stat("./photo/" + fileName)
	if err == nil {
		return fileName, nil
	}
	referer := "https://app-api.pixiv.net/"
	resp, err := http_util.Client().Timeout(p.APITimeout).Proxy(p.ProxyURL).Headers(map[string]string{
		"Referer": referer,
	}).Get(url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	f, err := os.Create("./photo/" + fileName)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	io.Copy(f, resp.Body)
	f.Close()
	return fileName, nil
}

// GetFileExt 获取文件的扩展名
func GetFileExt(fileAddr string) (string, error) {
	fileName := path.Base(fileAddr)
	reg, err := regexp.Compile("\\.(" + "jpg|png|jpeg|JPG|JPEG|PNG" + ")")
	if err != nil {
		return "", err
	}
	matchedExtArr := reg.FindAllString(fileName, -1)
	if matchedExtArr != nil && len(matchedExtArr) > 0 {
		ext := matchedExtArr[len(matchedExtArr)-1]
		return ext[1:], nil
	}
	return "", errors.New("获取文件扩展名失败")
}

func (p *PixivClient) Auth() error {
	if p.accessToken != "" && time.Now().Unix() < p.expireTime {
		return nil
	}
	json, err := http_util.Client().
		Timeout(p.APITimeout).
		Proxy(p.ProxyURL).
		PostFormGetJSON("https://oauth.secure.pixiv.net/auth/token", map[string]string{
			"client_id":      clientId,
			"client_secret":  clientSecret,
			"grant_type":     "refresh_token",
			"get_secure_url": "1",
			"refresh_token":  p.RefreshToken,
		})
	if err != nil {
		logx.ERR(err)
		return err
	}
	accessToken := json["access_token"].(string)
	expireTime := time.Now().Unix() + int64(json["expires_in"].(float64))
	if accessToken == "" || expireTime == 0 {
		return errors.New("accessToken or expireTime is empty, {accessToken:" + accessToken + ", expireTime:" + fmt.Sprintf("%d", expireTime) + "}")
	}
	p.accessToken = accessToken
	p.expireTime = expireTime
	return nil
}

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, hashSecret)
	return hex.EncodeToString(h.Sum(nil))
}

func (p *PixivClient) Rank() ([]string, error) {
	host := "https://app-api.pixiv.net"
	url := host + "/v1/illust/ranking"
	mode := "day_male_r18"
	filter := "for_ios"
	headers, err := p.getHeaders()
	if err != nil {
		return nil, err
	}
	json, err := http_util.Client().Timeout(p.APITimeout).Proxy(p.ProxyURL).Headers(headers).GetJSON(url, map[string]string{
		"mode":   mode,
		"filter": filter,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	photos := make([]string, 0)
	illusts := json["illusts"].([]interface{})
	for _, illust := range illusts {
		image := illust.(map[string]interface{})["image_urls"].(map[string]interface{})
		if image["large"] != "" {
			photos = append(photos, image["large"].(string))
		} else if image["medium"] != "" {
			photos = append(photos, image["medium"].(string))
		} else if image["square_medium"] != "" {
			photos = append(photos, image["square_medium"].(string))
		}
	}
	return photos, nil
}

func (p *PixivClient) Recommend() (map[string]interface{}, error) {
	req := apiHosts + "/v1/illust/recommended"
	headers, err := p.getHeaders()
	headers[`include_ranking_label`] = "true"
	jsonMap, err := http_util.Client().Timeout(p.APITimeout).Proxy(p.ProxyURL).Headers(headers).GetJSON(req, map[string]string{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return jsonMap, nil
}

// SearchIllust 搜索插画
func (p *PixivClient) SearchIllust(param SearchParam) (map[string]any, error) {
	url := apiHosts + "/v1/search/illust"
	headers, err := p.getHeaders()
	if err != nil {
		return nil, err
	}
	json, err := http_util.Client().Timeout(p.APITimeout).Proxy(p.ProxyURL).Headers(headers).GetJSON(url, map[string]string{
		"word":          param.Word,
		"search_target": string(param.SearchTarget),
		"sort":          string(param.Sort),
		"filter":        "for_ios",
		"offset":        str_util.ToStr(param.Offset),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return json, nil
}
