package appbuilder

type appBuilderSearchResponse struct {
	Search string `example:"Fox Plus 1.6 8V 2009"`
	Cars   []car
}

type expandedSearch struct {
	Search  string
	Year    int
	HasYear bool
}
type yearCandidate struct {
	Start int
	End   int
	Value string
}

type appBuilderSearchRequest struct {
	Search []string `json:"search" example:"Fox Plus 1.6 8V 2009"`
}
