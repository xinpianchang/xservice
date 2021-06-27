package dao

import (
	"fmt"

	"gorm.io/gorm"

	"{{.Module}}/internal/dto"
)

func ListPage(db *gorm.DB, form dto.PageForm, list interface{}) (*dto.Page, error) {
	current := form.GetPage()
	size := form.GetPageSize()

	page := dto.NewPage()
	page.Pagination.PageSize = size
	page.Pagination.Current = current

	db.Model(list).Count(&page.Pagination.Total)

	if page.Pagination.Total == 0 {
		return page, nil
	}

	offset := (current - 1) * size
	if err := db.Offset(offset).Limit(size).Find(list).Error; err != nil {
		return nil, err
	}

	page.List = list

	return page, nil
}

func ListPageRawSQL(db *gorm.DB, form dto.PageForm, list interface{}, sql string, params ...interface{}) (*dto.Page, error) {
	current := form.GetPage()
	size := form.GetPageSize()

	page := dto.NewPage()
	page.Pagination.PageSize = size
	page.Pagination.Current = current

	_ = db.
		Raw(fmt.Sprintf("select count(*) from (%s) all_for_count_", sql), params...).
		Row().Scan(&page.Pagination.Total)

	if page.Pagination.Total == 0 {
		return page, nil
	}

	offset := (current - 1) * size
	err := db.Raw(fmt.Sprintf("%s limit ?, ?", sql), append(append(make([]interface{}, 0, len(params)+2), params...), offset, size)...).Find(list).Error
	if err != nil {
		return nil, err
	}

	page.List = list

	return page, nil
}
