package pager

//The PageDetails struct is designed to be embeded in client request objects asking for a paged response, as well as embedded in the response to identify to clients what the current page and size of the potential response is
type PageDetails struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

//Offset takes the current pagesize multiplied by the provided pagenumber to derive the item starting index from 0
func (p *PageDetails) Offset() int {
	if p.PageNumber < 1 {
		p.PageNumber = 1
	}
	return p.PageSize * (p.PageNumber - 1)
}

//Result takes a total item count and builds a PagedResult
func (p *PageDetails) Result(totalItemCount int) PagedResult {
	result := PagedResult{
		PageDetails:    *p,
		TotalItemCount: totalItemCount,
		TotalPageCount: 0,
	}

	if totalItemCount > 0 && p.PageSize > 0 {
		result.TotalPageCount = ((totalItemCount - 1) / p.PageSize) + 1
	}

	return result
}

//The PagedResult embeds PageDetails and is intended to provide a uniform data structure to provide a paged response
type PagedResult struct {
	PageDetails
	TotalItemCount int         `json:"totalItemCount"`
	TotalPageCount int         `json:"totalPageCount"`
	Results        interface{} `json:"results"`
}
