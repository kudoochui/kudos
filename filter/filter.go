package filter

type Filter interface {
	Before(route string, msgReq interface{})
	After(route string, msgResp interface{})
}
