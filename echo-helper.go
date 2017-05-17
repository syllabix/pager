package pager

import (
	"strconv"

	"github.com/labstack/echo"
)

//GetPagingDetailsFromContext is a helper function for use with Echo
func GetPagingDetailsFromContext(c echo.Context) *PageDetails {
	pageRequest := new(PageDetails)
	pSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		pSize = 0
	}
	pNum, err := strconv.Atoi(c.QueryParam("pageNumber"))
	if err != nil {
		pNum = 0
	}
	pageRequest.PageSize = pSize
	pageRequest.PageNumber = pNum
	return pageRequest
}
