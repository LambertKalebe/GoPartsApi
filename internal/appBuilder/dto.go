package appbuilder

type appBuilderSearchResponse struct {
	MainSearchIndex int    `json:"mainSearchIndex" example:"1"`
	Search          string `json:"search" example:"Fox Plus 1.6 8V 2009"`
	Cars            []car  `json:"cars"`
	CarIds          []int  `json:"-"`
	Total           int    `json:"total"`
}

type appBuilderResponse struct {
	Results []appBuilderSearchResponse `json:"results"`
	CarIDs  []int                      `json:"carIds"`
	Total   int                        `json:"total"`
}

type expandedSearch struct {
	MainSearchIndex int
	Search          string
	Year            int
	HasYear         bool
}
type yearCandidate struct {
	Start int
	End   int
	Value string
}

type appBuilderSearchRequest struct {
	Search []string `json:"search" example:"Fox Plus 1.6 8V 2009"`
}
