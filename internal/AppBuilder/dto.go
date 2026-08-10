package appbuilder

type carSearchResponse struct {
	Cars []car
}

type carSearchRequest struct {
	Search string `query:"search" example:"Fox 2009"`
}
