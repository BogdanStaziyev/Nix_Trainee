package domain

type UrlQueryParams struct {
	Page        uint
	PageSize    uint
	ShowDeleted bool
	LastError   error
	LastSymbols bool
}
