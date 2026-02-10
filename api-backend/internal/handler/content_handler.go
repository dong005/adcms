package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/storage"
	"adcms/pkg/utils"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{categoryRepo: repository.NewCategoryRepository()}
}

type CreateCategoryRequest struct {
	ParentID    uint   `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int8   `json:"status"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	category := model.Category{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		ParentID:        req.ParentID,
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		Sort:            req.Sort,
		Status:          req.Status,
	}

	if err := h.categoryRepo.Create(&category); err != nil {
		utils.ServerError(c, "创建分类失败")
		return
	}
	utils.Success(c, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	category, err := h.categoryRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 6001, "分类不存在")
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	category.ParentID = req.ParentID
	category.Name = req.Name
	category.Slug = req.Slug
	category.Description = req.Description
	category.Sort = req.Sort
	category.Status = req.Status

	if err := h.categoryRepo.Update(category); err != nil {
		utils.ServerError(c, "更新分类失败")
		return
	}
	utils.Success(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	hasChildren, _ := h.categoryRepo.HasChildren(uint(id))
	if hasChildren {
		utils.Fail(c, 6002, "请先删除子分类")
		return
	}
	if err := h.categoryRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除分类失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *CategoryHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	categories, err := h.categoryRepo.FindAll(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, categories)
}

type TagHandler struct {
	tagRepo *repository.TagRepository
}

func NewTagHandler() *TagHandler {
	return &TagHandler{tagRepo: repository.NewTagRepository()}
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

func (h *TagHandler) Create(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	tag := model.Tag{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Slug:            req.Slug,
	}

	if err := h.tagRepo.Create(&tag); err != nil {
		utils.ServerError(c, "创建标签失败")
		return
	}
	utils.Success(c, tag)
}

func (h *TagHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tag, err := h.tagRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 7001, "标签不存在")
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tag.Name = req.Name
	tag.Slug = req.Slug

	if err := h.tagRepo.Update(tag); err != nil {
		utils.ServerError(c, "更新标签失败")
		return
	}
	utils.Success(c, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.tagRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除标签失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *TagHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	tags, err := h.tagRepo.List(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, tags)
}

type ArticleHandler struct {
	articleRepo *repository.ArticleRepository
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{articleRepo: repository.NewArticleRepository()}
}

type CreateArticleRequest struct {
	CategoryID uint   `json:"category_id"`
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Cover      string `json:"cover"`
	Status     int8   `json:"status"`
	TagIDs     []uint `json:"tag_ids"`
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	article := model.Article{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		CategoryID:      req.CategoryID,
		UserID:          userID,
		Title:           req.Title,
		Slug:            req.Slug,
		Summary:         req.Summary,
		Content:         req.Content,
		Cover:           req.Cover,
		Status:          req.Status,
	}

	if err := h.articleRepo.Create(&article); err != nil {
		utils.ServerError(c, "创建文章失败")
		return
	}

	if len(req.TagIDs) > 0 {
		h.articleRepo.AssignTags(article.ID, req.TagIDs)
	}

	utils.Success(c, article)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	article, err := h.articleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 8001, "文章不存在")
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	article.CategoryID = req.CategoryID
	article.Title = req.Title
	article.Slug = req.Slug
	article.Summary = req.Summary
	article.Content = req.Content
	article.Cover = req.Cover
	article.Status = req.Status

	if err := h.articleRepo.Update(article); err != nil {
		utils.ServerError(c, "更新文章失败")
		return
	}

	if req.TagIDs != nil {
		h.articleRepo.AssignTags(article.ID, req.TagIDs)
	}

	utils.Success(c, article)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.articleRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除文章失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ArticleHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	keyword := c.Query("keyword")

	articles, total, err := h.articleRepo.List(tenantID, page, pageSize, uint(categoryID), int8(status), keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, articles, total, page, pageSize)
}

func (h *ArticleHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	article, err := h.articleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 8001, "文章不存在")
		return
	}
	utils.Success(c, article)
}

func (h *ArticleHandler) Publish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.articleRepo.UpdateStatus(uint(id), 1); err != nil {
		utils.ServerError(c, "发布失败")
		return
	}
	utils.SuccessWithMessage(c, "发布成功", nil)
}

func (h *ArticleHandler) Draft(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.articleRepo.UpdateStatus(uint(id), 0); err != nil {
		utils.ServerError(c, "操作失败")
		return
	}
	utils.SuccessWithMessage(c, "已存为草稿", nil)
}

type MediaHandler struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaHandler() *MediaHandler {
	return &MediaHandler{mediaRepo: repository.NewMediaRepository()}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}

	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	ext := strings.ToLower(filepath.Ext(file.Filename))
	mediaType := "file"
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
		mediaType = "image"
	case ".mp4", ".avi", ".mov", ".wmv":
		mediaType = "video"
	case ".mp3", ".wav", ".flac":
		mediaType = "audio"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx":
		mediaType = "document"
	}

	// 使用 Storage 接口上传
	dir := fmt.Sprintf("%d", tenantID)
	fileInfo, err := storage.Default.Upload(file, dir)
	if err != nil {
		utils.ServerError(c, "文件上传失败: "+err.Error())
		return
	}

	media := model.Media{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		UserID:          userID,
		Name:            fileInfo.Name,
		Path:            fileInfo.Path,
		Type:            mediaType,
		Size:            fileInfo.Size,
		MimeType:        fileInfo.MimeType,
	}

	if err := h.mediaRepo.Create(&media); err != nil {
		utils.ServerError(c, "保存记录失败")
		return
	}

	// 返回时附带访问URL
	media.Path = fileInfo.URL

	utils.Success(c, media)
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	media, err := h.mediaRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 9001, "文件不存在")
		return
	}

	// 使用 Storage 接口删除文件
	_ = storage.Default.Delete(media.Path)

	if err := h.mediaRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *MediaHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	mediaType := c.Query("type")

	medias, total, err := h.mediaRepo.List(tenantID, page, pageSize, mediaType)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	for i := range medias {
		medias[i].Path = storage.Default.GetURL(medias[i].Path)
	}
	utils.SuccessWithPage(c, medias, total, page, pageSize)
}
