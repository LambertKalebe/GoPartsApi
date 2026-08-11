package appbuilder

type appBuilderSearchResponse struct {
	Search string
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

type debugResponse struct {
	Debug string
}

type appBuilderSearchRequest struct {
	Search []string `json:"search" example:"Fox 2009"`
}
