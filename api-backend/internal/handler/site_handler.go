package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	siteRepo *repository.SiteRepository
}

func NewSiteHandler() *SiteHandler {
	return &SiteHandler{siteRepo: repository.NewSiteRepository()}
}

func (h *SiteHandler) List(c *gin.Context) {
	isAdmin := middleware.GetIsAdmin(c)
	var tenantID uint
	if isAdmin == 2 {
		tenantID = 0
	} else {
		tenantID = middleware.GetTenantID(c)
	}
	keyword := c.Query("keyword")

	sites, err := h.siteRepo.List(tenantID, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, sites)
}

func (h *SiteHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type"`
		URL      string `json:"url"`
		Image    string `json:"image"`
		IsDomain int8   `json:"is_domain"`
		Status   int8   `json:"status"`
		Sort     int    `json:"sort"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	site := model.Site{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Type:            req.Type,
		URL:             req.URL,
		Image:           req.Image,
		IsDomain:        req.IsDomain,
		Status:          req.Status,
		Sort:            req.Sort,
		Remark:          req.Remark,
	}

	if err := h.siteRepo.Create(&site); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.SuccessWithMessage(c, "创建成功", site)
}

func (h *SiteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		URL      string `json:"url"`
		Image    string `json:"image"`
		IsDomain int8   `json:"is_domain"`
		Status   int8   `json:"status"`
		Sort     int    `json:"sort"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	site, err := h.siteRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 404, "站点不存在")
		return
	}

	if req.Name != "" {
		site.Name = req.Name
	}
	site.Type = req.Type
	site.URL = req.URL
	site.Image = req.Image
	site.IsDomain = req.IsDomain
	site.Status = req.Status
	site.Sort = req.Sort
	site.Remark = req.Remark

	if err := h.siteRepo.Update(site); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.SuccessWithMessage(c, "更新成功", site)
}

func (h *SiteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.siteRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
