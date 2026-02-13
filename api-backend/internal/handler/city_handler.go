package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	cityRepo *repository.CityRepository
}

func NewCityHandler() *CityHandler {
	return &CityHandler{cityRepo: repository.NewCityRepository()}
}

func (h *CityHandler) List(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.DefaultQuery("pid", "0"), 10, 64)
	cities, err := h.cityRepo.ListByPID(uint(pid))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, cities)
}

func (h *CityHandler) Tree(c *gin.Context) {
	maxLevel, _ := strconv.ParseInt(c.DefaultQuery("max_level", "2"), 10, 8)
	tree, err := h.cityRepo.Tree(0, int8(maxLevel))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, tree)
}

func (h *CityHandler) Create(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	var req struct {
		PID      uint    `json:"pid"`
		Level    int8    `json:"level"`
		Name     string  `json:"name" binding:"required"`
		Citycode string  `json:"citycode"`
		Adcode   string  `json:"adcode"`
		PAdcode  string  `json:"p_adcode"`
		Lng      float64 `json:"lng"`
		Lat      float64 `json:"lat"`
		Sort     int     `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	city := model.City{
		PID:      req.PID,
		Level:    req.Level,
		Name:     req.Name,
		Citycode: req.Citycode,
		Adcode:   req.Adcode,
		PAdcode:  req.PAdcode,
		Lng:      req.Lng,
		Lat:      req.Lat,
		Sort:     req.Sort,
	}
	if err := h.cityRepo.Create(&city); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", city)
}

func (h *CityHandler) Update(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	city, err := h.cityRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "区域不存在")
		return
	}
	var req struct {
		Name     string  `json:"name"`
		Citycode string  `json:"citycode"`
		Adcode   string  `json:"adcode"`
		PAdcode  string  `json:"p_adcode"`
		Lng      float64 `json:"lng"`
		Lat      float64 `json:"lat"`
		Sort     int     `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if req.Name != "" {
		city.Name = req.Name
	}
	city.Citycode = req.Citycode
	city.Adcode = req.Adcode
	city.PAdcode = req.PAdcode
	city.Lng = req.Lng
	city.Lat = req.Lat
	city.Sort = req.Sort

	if err := h.cityRepo.Update(city); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", city)
}

func (h *CityHandler) Delete(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if err := h.cityRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
