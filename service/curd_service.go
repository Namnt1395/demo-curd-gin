package service

import (
	"demo-curd/dao"
	"demo-curd/dto/request"
	"demo-curd/dto/response"
	"demo-curd/model"
	"demo-curd/util"
	"github.com/jinzhu/copier"
)

type CurdService struct {
	CurdDao *dao.CurdDao
}

func (s *CurdService) Create(dto *request.CurdDTO) (*response.CurdDTO, error) {
	var curd model.Curd
	if err1 := dto.Validate(); err1 != nil {
		return nil, err1
	}
	util.Must(copier.Copy(&curd, &dto))
	_, err1 := s.CurdDao.Create(&curd)
	util.Must(err1)
	var res response.CurdDTO
	util.Must(copier.Copy(&res, &curd))
	return &res, nil
}
