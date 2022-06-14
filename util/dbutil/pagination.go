package dbutil

import (
	"demo-curd/dto/request"
	"demo-curd/util/constant"
	"gorm.io/gorm"
)

func Pagination(page request.Page) func(db *gorm.DB) *gorm.DB {
	if page.Size == 0 {
		page.Size = constant.DefaultPageSize
	}
	if page.Sort == "" {
		page.Sort = constant.DefaultPageSort
	}
	if page.Page == 0 {
		page.Page = constant.DefaultPage
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page.Page - 1) * page.Size).
			Limit(page.Size).
			Order(page.Sort)
	}
}
