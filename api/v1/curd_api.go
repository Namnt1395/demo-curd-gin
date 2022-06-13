package v1

import (
	"demo-curd/dto/request"
	"demo-curd/dto/response"
	"demo-curd/service"
	"demo-curd/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CurdV1Api struct {
	CurdService *service.CurdService
}

// Create
// @Summary Create new curd
// @Description Create new curd
// @Tags CURD
// @Security ApiKeyAuth
// @Accept json
// @Param body body request.curdDTO true "JSON body"
// @Success 200 {object} response.curdDTO
// @Failure 500 {object} interface{} "{"error_code": "<Mã lỗi>", "error_msg": "<Nội dung lỗi>"}"
// @Router /api/v1/curd [post]
func (r CurdV1Api) Create(c *gin.Context) {
	var curdDTO request.CurdDTO
	util.Must(c.BindJSON(&curdDTO))
	res, err := r.CurdService.Create(&curdDTO)
	util.Must(err)
	c.JSON(http.StatusOK, response.Response{
		Data: res,
	})
}
