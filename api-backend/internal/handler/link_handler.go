package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	linkRepo *repository.LinkRepository
}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{linkRepo: repository.NewLinkRepository()}
}

func (h *LinkHandler) List(c *gin.Context) {
	isAdmin := middleware.GetIsAdmin(c)
	var tenantID uint
	if isAdmin == 2 {
		tenantID = 0
	} else {
		tenantID = middleware.GetTenantID(c)
	}
	keyword := c.Query("keyword")

	links, err := h.linkRepo.List(tenantID, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, links)
}

func (h *LinkHandler) Create(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		URL    string `json:"url" binding:"required"`
		Logo   string `json:"logo"`
		Desc   string `json:"desc"`
		Sort   int    `json:"sort"`
		Status int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	link := model.Link{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		URL:             req.URL,
		Logo:            req.Logo,
		Desc:            req.Desc,
		Sort:            req.Sort,
		Status:          req.Status,
	}

	if err := h.linkRepo.Create(&link); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", link)
}

func (h *LinkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Logo   string `json:"logo"`
		Desc   string `json:"desc"`
		Sort   int    `json:"sort"`
		Status int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	link, err := h.linkRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "友链不存在")
		return
	}

	if req.Name != "" {
		link.Name = req.Name
	}
	if req.URL != "" {
		link.URL = req.URL
	}
	link.Logo = req.Logo
	link.Desc = req.Desc
	link.Sort = req.Sort
	link.Status = req.Status

	if err := h.linkRepo.Update(link); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", link)
}

func (h *LinkHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.linkRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
