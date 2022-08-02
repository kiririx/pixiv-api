package pixiv_api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kiririx/krutils/algox"
	"github.com/kiririx/krutils/filex"
	"github.com/kiririx/krutils/httpx"
	"github.com/kiririx/krutils/logx"
	"github.com/kiririx/krutils/strx"
	"io"
	"os"
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
	ext, _ := filex.GetUrlFileExt(url, photoExt)
	fileName := algox.MD5(url) + "." + ext
	_, err := os.Stat("./photo/" + fileName)
	if err == nil {
		return fileName, nil
	}
	referer := "https://app-api.pixiv.net/"
	resp, err := p.dwClient().Headers(map[string]string{
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

var photoExt = []string{"jpg", "png", "jpeg", "JPG", "JPEG", "PNG"}

func (p *PixivClient) Auth() error {
	if p.accessToken != "" && time.Now().Unix() < p.expireTime {
		return nil
	}
	json, err := httpx.Client().
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

// IllustsRank 插画排行
func (p *PixivClient) IllustsRank() ([]string, error) {
	host := "https://app-api.pixiv.net"
	url := host + "/v1/illust/ranking"
	mode := "day_male_r18"
	filter := "for_ios"
	headers, err := p.getHeaders()
	if err != nil {
		return nil, err
	}
	json, err := p.apiClient().Headers(headers).GetJSON(url, map[string]string{
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

// IllustsRecommend 插画推荐
func (p *PixivClient) IllustsRecommend() (map[string]interface{}, error) {
	req := apiHosts + "/v1/illust/recommended"
	headers, err := p.getHeaders()
	headers[`include_ranking_label`] = "true"
	jsonMap, err := p.apiClient().Headers(headers).GetJSON(req, map[string]string{})
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
	json, err := p.apiClient().Headers(headers).GetJSON(url, map[string]string{
		"word":          param.Word,
		"search_target": string(param.SearchTarget),
		"sort":          string(param.Sort),
		"filter":        "for_ios",
		"offset":        strx.ToStr(param.Offset),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return json, nil
}

// UserDetail 获取用户详情
func (p PixivClient) UserDetail() {

}

// UserIllusts 获取用户插画列表
func (p PixivClient) UserIllusts() {

}

// UserBookmarksIllusts 用户收藏作品列表
func (p PixivClient) UserBookmarksIllusts() {

}

// FollowIllusts 获取用户关注的插画
func (p PixivClient) FollowIllusts() {

}

// IllustsDetail 作品详情
func (p PixivClient) IllustsDetail() {

}

// IllustComments 获取作品评论
func (p PixivClient) IllustComments() {

}

// IllustRelated 获取相关作品
func (p PixivClient) IllustRelated() {

}

// SearchNovel 搜索小说
func (p PixivClient) SearchNovel() {

}

// SearchUser 搜索用户
func (p PixivClient) SearchUser() {

}

// IllustsBookmarkDetail 作品收藏详情
func (p PixivClient) IllustsBookmarkDetail() {

}

// IllustsBookmarkAdd 新增收藏
func (p PixivClient) IllustsBookmarkAdd() {

}

// UserFollowAdd 关注用户
func (p PixivClient) UserFollowAdd() {

}

// UserFollowDelete 取消关注用户
func (p PixivClient) UserFollowDelete() {

}

// UserBookmarkTagsIllusts 用户收藏标签列表
func (p PixivClient) UserBookmarkTagsIllusts() {

}

// UserFollowing 关注的用户列表
func (p PixivClient) UserFollowing() {

}

// UserFollower 被关注的用户列表
func (p PixivClient) UserFollower() {

}

// UserMyPixiv 我的pixiv朋友
func (p *PixivClient) UserMyPixiv() {

}

// UserBlack 黑名单用户
func (p *PixivClient) UserBlack() {

}

// UgoiraMetadata 获取ugoira信息
func (p *PixivClient) UgoiraMetadata() {

}

// UserNovels 用户小说列表
func (p *PixivClient) UserNovels() {

}

// NovelSeries 小说系列详情
func (p *PixivClient) NovelSeries() {

}

// NovelDetail 小说详情
func (p *PixivClient) NovelDetail() {

}

// NovelText 小说文本
func (p *PixivClient) NovelText() {

}

// IllustsNew 大家的新作
func (p *PixivClient) IllustsNew() {

}

// ShowcaseArticle 特辑详情
func (p *PixivClient) ShowcaseArticle() {

}
