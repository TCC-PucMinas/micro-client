package helpers

import "fmt"

type Paginate struct {
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
	Query string `json:"count"`
}

func (p *Paginate) PaginateMounted() {
	switch p.Page {
	case 0:
		p.Page = 0
	case 1:
		p.Page = 0
	default:
		p.Page = (p.Page * 10) - 10
	}
}

func (p *Paginate) MountedQuery(table string) {
	p.Query = fmt.Sprintf("CAST(((SELECT COUNT(*) FROM `%v`)/10) * 10 as UNSIGNED) AS total", table)
}
