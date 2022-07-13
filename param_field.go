package pixiv_api

var (
	SearchTarget _searchTarget
	Sort         _sort
)

func init() {
	SearchTarget = _searchTarget{
		PartialMatchForTags: "partial_match_for_tags",
		ExactMatchForTags:   "exact_match_for_tags",
		TitleAndCaption:     "title_and_caption",
	}
	Sort = _sort{
		DateDesc:    "date_desc",
		DateAsc:     "date_asc",
		PopularDesc: "popular_desc",
	}
}

type _searchTarget struct {
	PartialMatchForTags string
	ExactMatchForTags   string
	TitleAndCaption     string
}

// "date_desc", "date_asc", "popular_desc"
type _sort struct {
	DateDesc    string
	DateAsc     string
	PopularDesc string
}
