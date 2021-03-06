// simple Page struct layout for RESTFul API pagination

package dto

// Page is a struct for RESTful API pagination
type Page struct {
	List       interface{} `json:"list"`
	Pagination *Pagination `json:"pagination"`
	Meta       interface{} `json:"meta,omitempty"`
}

// Pagination respects the pagination of the current page
type Pagination struct {
	Total    int64 `json:"total"`    // 总条数
	PageSize int   `json:"pageSize"` // 页大小
	Current  int   `json:"current"`  // 当前页码
}

const (
	defaultMinPageSize = 20
	defaultMaxPageSize = 100
)

// NewPage create page instance
func NewPage() *Page {
	return &Page{
		List: []interface{}{},
		Pagination: &Pagination{
			Total:    0,
			PageSize: 0,
			Current:  1,
		},
	}
}

// PageForm is a struct for Page request
type PageForm struct {
	Page     int `json:"page" form:"page" query:"page"`
	PageSize int `json:"pageSize" form:"pageSize" query:"pageSize"`
}

// GetPage get current page number
func (p *PageForm) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

// GetPageSize get current page size
func (p *PageForm) GetPageSize() int {
	if p.PageSize == 0 {
		return defaultMinPageSize
	}

	if p.PageSize > defaultMaxPageSize {
		return defaultMaxPageSize
	}

	return p.PageSize
}
