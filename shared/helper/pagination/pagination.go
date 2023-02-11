package pagination

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Param struct {
	Limit      uint
	Page       uint
	Sort       []ParamSort
	TotalRows  int64
	TotalPages uint
}

type ParamSort struct {
	Column string
	Order  string
}

const (
	OrderDesc    = "DESC"
	OrderAsc     = "ASC"
	OrderDefault = "id DESC"
)

func (p *Param) GetOffset() uint {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Param) GetLimit() uint {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Param) GetPage() uint {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Param) GetSort() string {
	var sortBy string
	if len(p.Sort) == 0 {
		return OrderDefault
	}

	var sortBys []string
	for _, s := range p.Sort {
		if strings.ToUpper(s.Order) != OrderAsc && strings.ToUpper(s.Order) != OrderDesc {
			continue
		}

		sortBys = append(sortBys, fmt.Sprintf("%s %s", s.Column, s.Order))
	}

	if len(sortBys) == 0 {
		return OrderDefault
	}

	sortBy = strings.Join(sortBys, ", ")

	return sortBy
}

func Paginate(value interface{}, param *Param, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	param.TotalRows = totalRows
	totalPages := uint(math.Ceil(float64(totalRows) / float64(param.Limit)))
	param.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(param.GetOffset())).Limit(int(param.GetLimit())).Order(param.GetSort())
	}
}
