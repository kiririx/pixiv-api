package pixiv_api

import "github.com/kiririx/pixiv_api/param"

// SearchParam 搜索参数
//
// Word 搜索关键字
//
// SearchTarget 搜索类型
//
// 可选值：
//
// - SearchTargetPartialMatchForTags 标签部分一致
//
// - SearchTargetExactMatchForTags 标签完全一致
//
// - SearchTargetTitleAndCaption 标题说明文
//
// Sort 排序
//
// 可选值：
//
// - SortDateDesc 按时间降序
//
// - SortDateAsc 按时间升序
//
// Offset 起始数据位置
//
// Auth 是否展示已经登录的数据
type SearchParam struct {
	Word         string
	SearchTarget param.SearchTargetType
	Sort         param.SortType
	Offset       int
	Auth         bool
	StartDate    string
	EndDate      string
	Duration     string
}

func init() {

}
