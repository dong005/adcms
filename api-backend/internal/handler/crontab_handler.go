package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CrontabHandler struct {
	crontabRepo *repository.CrontabRepository
}

func NewCrontabHandler() *CrontabHandler {
	return &CrontabHandler{crontabRepo: repository.NewCrontabRepository()}
}

func (h *CrontabHandler) List(c *gin.Context) {
	isAdmin := middleware.GetIsAdmin(c)
	var tenantID uint
	if isAdmin == 2 {
		tenantID = 0
	} else {
		tenantID = middleware.GetTenantID(c)
	}

	crontabs, err := h.crontabRepo.List(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, crontabs)
}

func (h *CrontabHandler) Create(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		Expression string `json:"expression" binding:"required"`
		Command    string `json:"command" binding:"required"`
		Status     int8   `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	crontab := model.Crontab{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Expression:      req.Expression,
		Command:         req.Command,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	if err := h.crontabRepo.Create(&crontab); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", crontab)
}

func (h *CrontabHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name       string `json:"name"`
		Expression string `json:"expression"`
		Command    string `json:"command"`
		Status     int8   `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	crontab, err := h.crontabRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "定时任务不存在")
		return
	}

	if req.Name != "" {
		crontab.Name = req.Name
	}
	if req.Expression != "" {
		crontab.Expression = req.Expression
	}
	if req.Command != "" {
		crontab.Command = req.Command
	}
	crontab.Status = req.Status
	crontab.Remark = req.Remark

	if err := h.crontabRepo.Update(crontab); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", crontab)
}

func (h *CrontabHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.crontabRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
