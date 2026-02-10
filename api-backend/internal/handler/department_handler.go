package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	deptRepo *repository.DepartmentRepository
}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{deptRepo: repository.NewDepartmentRepository()}
}

type CreateDepartmentRequest struct {
	ParentID uint   `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	dept := model.Department{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		ParentID:        req.ParentID,
		Name:            req.Name,
		Code:            req.Code,
		Leader:          req.Leader,
		Phone:           req.Phone,
		Email:           req.Email,
		Sort:            req.Sort,
		Status:          req.Status,
	}

	if err := h.deptRepo.Create(&dept); err != nil {
		utils.ServerError(c, "创建部门失败")
		return
	}
	utils.Success(c, dept)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dept, err := h.deptRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 5001, "部门不存在")
		return
	}

	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	dept.ParentID = req.ParentID
	dept.Name = req.Name
	dept.Code = req.Code
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email
	dept.Sort = req.Sort
	dept.Status = req.Status

	if err := h.deptRepo.Update(dept); err != nil {
		utils.ServerError(c, "更新部门失败")
		return
	}
	utils.Success(c, dept)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	hasChildren, _ := h.deptRepo.HasChildren(uint(id))
	if hasChildren {
		utils.Fail(c, 5002, "请先删除子部门")
		return
	}

	hasUsers, _ := h.deptRepo.HasUsers(uint(id))
	if hasUsers {
		utils.Fail(c, 5003, "该部门下还有用户，不能删除")
		return
	}

	if err := h.deptRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除部门失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	depts, err := h.deptRepo.FindAll(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, depts)
}

func (h *DepartmentHandler) Tree(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	depts, err := h.deptRepo.FindAll(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	tree := repository.BuildDepartmentTree(depts, 0)
	utils.Success(c, tree)
}
