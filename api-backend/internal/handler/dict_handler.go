package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DictHandler struct {
	dictRepo *repository.DictRepository
}

func NewDictHandler() *DictHandler {
	return &DictHandler{dictRepo: repository.NewDictRepository()}
}

// ========== DictType ==========

func (h *DictHandler) ListTypes(c *gin.Context) {
	isAdmin := middleware.GetIsAdmin(c)
	var tenantID uint
	if isAdmin == 2 {
		tenantID = 0
	} else {
		tenantID = middleware.GetTenantID(c)
	}
	keyword := c.Query("keyword")

	types, err := h.dictRepo.ListTypes(tenantID, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, types)
}

func (h *DictHandler) CreateType(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Code   string `json:"code" binding:"required"`
		Sort   int    `json:"sort"`
		Status int8   `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查编码唯一性
	if _, err := h.dictRepo.FindTypeByCode(req.Code); err == nil {
		utils.Fail(c, 3001, "字典类型编码已存在")
		return
	}

	tenantID := middleware.GetTenantID(c)
	dictType := model.DictType{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Code:            req.Code,
		Sort:            req.Sort,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	if err := h.dictRepo.CreateType(&dictType); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", dictType)
}

func (h *DictHandler) UpdateType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Sort   int    `json:"sort"`
		Status int8   `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	dictType, err := h.dictRepo.FindTypeByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "字典类型不存在")
		return
	}

	// 检查编码唯一性（排除自身）
	if req.Code != "" && req.Code != dictType.Code {
		if _, err := h.dictRepo.FindTypeByCode(req.Code); err == nil {
			utils.Fail(c, 3001, "字典类型编码已存在")
			return
		}
		dictType.Code = req.Code
	}

	if req.Name != "" {
		dictType.Name = req.Name
	}
	dictType.Sort = req.Sort
	dictType.Status = req.Status
	dictType.Remark = req.Remark

	if err := h.dictRepo.UpdateType(dictType); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", dictType)
}

func (h *DictHandler) DeleteType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.dictRepo.DeleteType(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ========== Dict ==========

func (h *DictHandler) ListDicts(c *gin.Context) {
	typeID, err := strconv.ParseUint(c.Query("type_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误：type_id 必填")
		return
	}

	dicts, err := h.dictRepo.ListDicts(uint(typeID))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, dicts)
}

func (h *DictHandler) CreateDict(c *gin.Context) {
	var req struct {
		DictTypeID uint   `json:"dict_type_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		Value      string `json:"value" binding:"required"`
		Sort       int    `json:"sort"`
		Status     int8   `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	dict := model.Dict{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		DictTypeID:      req.DictTypeID,
		Name:            req.Name,
		Value:           req.Value,
		Sort:            req.Sort,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	if err := h.dictRepo.CreateDict(&dict); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", dict)
}

func (h *DictHandler) UpdateDict(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Sort   int    `json:"sort"`
		Status int8   `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	dict, err := h.dictRepo.FindDictByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "字典数据不存在")
		return
	}

	if req.Name != "" {
		dict.Name = req.Name
	}
	if req.Value != "" {
		dict.Value = req.Value
	}
	dict.Sort = req.Sort
	dict.Status = req.Status
	dict.Remark = req.Remark

	if err := h.dictRepo.UpdateDict(dict); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", dict)
}

func (h *DictHandler) DeleteDict(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.dictRepo.DeleteDict(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// GetDictsByCode 根据编码获取字典数据（公共接口）
func (h *DictHandler) GetDictsByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		utils.BadRequest(c, "参数错误")
		return
	}

	dicts, err := h.dictRepo.GetDictsByCode(code)
	if err != nil {
		utils.Fail(c, 404, "字典类型不存在")
		return
	}
	utils.Success(c, dicts)
}
