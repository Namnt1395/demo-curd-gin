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
	curdDao dao.CurdDao
}

func (s CurdService) Create(dto *request.CurdDTO) (*response.CurdDTO, error) {
	var curd model.Curd
	util.Must(copier.Copy(&curd, &dto))
	_, err1 := s.curdDao.Create(&curd)
	util.Must(err1)
	var res response.CurdDTO
	util.Must(copier.Copy(&res, &curd))
	return &res, nil
}
