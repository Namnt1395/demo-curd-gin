package dao

import (
	"demo-curd/database"
	"demo-curd/model"
	"demo-curd/util"
	"errors"
	"gorm.io/gorm"
)

type CurdDao struct {
	Db *database.Database
}

func (r CurdDao) Create(curd *model.Curd) (*model.Curd, error) {
	err := r.Db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(curd).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return curd, nil
}

func (r CurdDao) UpdateDepartment(department *model.Curd) (*model.Curd, error) {
	if err := r.Db.DB.Save(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (r CurdDao) DeleteDepartment(department *model.Curd) (*model.Curd, error) {
	if err := r.Db.DB.Delete(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (r CurdDao) GetDepartmentDetail(id uint64) (*model.Curd, error) {
	var curd model.Curd
	if err := r.Db.DB.Where("id = ?", id).First(&curd).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			util.Must(err)
		}
	}
	return &curd, nil
}

func (r CurdDao) List(page int, size int, sort string) (*[]model.Curd, error) {
	var curds []model.Curd
	offset := (page - 1) * size
	result := r.Db.DB.Offset(offset).Limit(size).Order(sort).Model(&model.Curd{}).Find(&curds)
	if result.Error != nil {
		return nil, result.Error
	}
	return &curds, nil
}
