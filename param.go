package pixiv_api

// SearchParam 搜索参数
// Word: 搜索关键字
// SearchTarget 搜索类型:
// - partial_match_for_tags 标签部分一致
// - exact_match_for_tags 标签完全一致
// - title_and_caption 标题说明文
// Sort 排序:
// - date_desc 按时间降序
// - date_asc 按时间升序
// Offset 起始数据位置
// Auth 是否展示已经登录的数据
type SearchParam struct {
	Word         string
	SearchTarget string
	Sort         string
	Offset       int
	Auth         bool
	StartDate    string
	EndDate      string
	Duration     string
}

func init() {
}
