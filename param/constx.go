package param

type SearchTargetType string
type SortType string

var (
	// SearchTargetPartialMatchForTags 标签部分一致
	SearchTargetPartialMatchForTags = SearchTargetType("partial_match_for_tags")
	// SearchTargetExactMatchForTags 标签完全一致
	SearchTargetExactMatchForTags = SearchTargetType("exact_match_for_tags")
	// SearchTargetTitleAndCaption 标题说明文
	SearchTargetTitleAndCaption = SearchTargetType("title_and_caption")
	// SortDateDesc 按时间降序
	SortDateDesc = SortType("date_desc")
	// SortDateAsc 按时间升序
	SortDateAsc = SortType("date_asc")
)
