package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	tableName   string
	moduleName  string
	outputDir   string
	withFrontend bool
)

func init() {
	flag.StringVar(&tableName, "table", "", "数据库表名 (必填, 如: orders)")
	flag.StringVar(&moduleName, "module", "", "模块名 (可选, 默认使用表名单数形式)")
	flag.StringVar(&outputDir, "output", ".", "输出根目录 (默认当前目录)")
	flag.BoolVar(&withFrontend, "frontend", false, "是否同时生成前端页面")
}

// TemplateData 模板数据
type TemplateData struct {
	TableName      string // orders
	ModuleName     string // order
	ModelName      string // Order
	ModuleNamePlural string // orders
	HandlerName    string // OrderHandler
	RepoName       string // OrderRepository
	RoutePath      string // /orders
}

func main() {
	flag.Parse()

	if tableName == "" {
		fmt.Println("用法: go run cmd/gen/main.go -table=<表名> [-module=<模块名>] [-output=<输出目录>] [-frontend]")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  go run cmd/gen/main.go -table=orders")
		fmt.Println("  go run cmd/gen/main.go -table=product_categories -module=product_category")
		fmt.Println("  go run cmd/gen/main.go -table=orders -frontend")
		os.Exit(1)
	}

	if moduleName == "" {
		moduleName = strings.TrimSuffix(tableName, "s")
	}

	data := TemplateData{
		TableName:      tableName,
		ModuleName:     moduleName,
		ModelName:      toPascalCase(moduleName),
		ModuleNamePlural: tableName,
		HandlerName:    toPascalCase(moduleName) + "Handler",
		RepoName:       toPascalCase(moduleName) + "Repository",
		RoutePath:      "/" + tableName,
	}

	fmt.Printf("生成代码: table=%s module=%s model=%s\n", data.TableName, data.ModuleName, data.ModelName)

	generateFile("model", data, modelTemplate)
	generateFile("repository", data, repoTemplate)
	generateFile("handler", data, handlerTemplate)

	fmt.Println("\n=== 生成完成 ===")
	fmt.Println("请手动完成以下步骤:")
	fmt.Printf("1. 编辑 internal/model/%s.go 添加字段\n", data.ModuleName)
	fmt.Printf("2. 在 pkg/database/migrate.go 的 AutoMigrate 中添加 &model.%s{}\n", data.ModelName)
	fmt.Printf("3. 在 internal/router/router.go 中注册路由:\n")
	fmt.Printf("   %sHandler := handler.New%s()\n", toCamelCase(data.ModuleName), data.HandlerName)
	fmt.Printf("   %s := protected.Group(\"%s\")\n", data.ModuleNamePlural, data.RoutePath)
	fmt.Printf("   {\n")
	fmt.Printf("       %s.GET(\"\", %sHandler.List)\n", data.ModuleNamePlural, toCamelCase(data.ModuleName))
	fmt.Printf("       %s.GET(\"/:id\", %sHandler.Detail)\n", data.ModuleNamePlural, toCamelCase(data.ModuleName))
	fmt.Printf("       %s.POST(\"\", %sHandler.Create)\n", data.ModuleNamePlural, toCamelCase(data.ModuleName))
	fmt.Printf("       %s.PUT(\"/:id\", %sHandler.Update)\n", data.ModuleNamePlural, toCamelCase(data.ModuleName))
	fmt.Printf("       %s.DELETE(\"/:id\", %sHandler.Delete)\n", data.ModuleNamePlural, toCamelCase(data.ModuleName))
	fmt.Printf("   }\n")
}

func generateFile(fileType string, data TemplateData, tmplStr string) {
	tmpl, err := template.New(fileType).Parse(tmplStr)
	if err != nil {
		fmt.Printf("模板解析失败 [%s]: %v\n", fileType, err)
		return
	}

	var dir, filename string
	switch fileType {
	case "model":
		dir = filepath.Join(outputDir, "internal/model")
		filename = data.ModuleName + ".go"
	case "repository":
		dir = filepath.Join(outputDir, "internal/repository")
		filename = data.ModuleName + "_repo.go"
	case "handler":
		dir = filepath.Join(outputDir, "internal/handler")
		filename = data.ModuleName + "_handler.go"
	}

	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  [跳过] %s 已存在\n", path)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("创建文件失败 [%s]: %v\n", path, err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("模板渲染失败 [%s]: %v\n", path, err)
		return
	}

	fmt.Printf("  [生成] %s\n", path)
}

// toPascalCase 下划线转大驼峰: order_item -> OrderItem
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase 下划线转小驼峰: order_item -> orderItem
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ===== 模板定义 =====

var modelTemplate = `package model

// {{.ModelName}} {{.TableName}}表
type {{.ModelName}} struct {
	TenantBaseModel
	Name   string ` + "`" + `gorm:"size:100;not null" json:"name"` + "`" + `
	Status int8   ` + "`" + `gorm:"default:1" json:"status"` + "`" + `
	// TODO: 添加更多字段
}

func ({{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
`

var repoTemplate = `package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type {{.RepoName}} struct {
	db *gorm.DB
}

func New{{.RepoName}}() *{{.RepoName}} {
	return &{{.RepoName}}{db: database.DB}
}

func (r *{{.RepoName}}) Create(item *model.{{.ModelName}}) error {
	return r.db.Create(item).Error
}

func (r *{{.RepoName}}) Update(item *model.{{.ModelName}}) error {
	return r.db.Save(item).Error
}

func (r *{{.RepoName}}) Delete(id uint) error {
	return r.db.Delete(&model.{{.ModelName}}{}, id).Error
}

func (r *{{.RepoName}}) FindByID(id uint) (*model.{{.ModelName}}, error) {
	var item model.{{.ModelName}}
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *{{.RepoName}}) List(tenantID uint, page, pageSize int, keyword string) ([]model.{{.ModelName}}, int64, error) {
	var items []model.{{.ModelName}}
	var total int64

	query := r.db.Where("tenant_id = ?", tenantID)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Model(&model.{{.ModelName}}{}).Count(&total)

	err := query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}
`

var handlerTemplate = `package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type {{.HandlerName}} struct {
	repo *repository.{{.RepoName}}
}

func New{{.HandlerName}}() *{{.HandlerName}} {
	return &{{.HandlerName}}{repo: repository.New{{.RepoName}}()}
}

type Create{{.ModelName}}Request struct {
	Name   string ` + "`" + `json:"name" binding:"required"` + "`" + `
	Status int8   ` + "`" + `json:"status"` + "`" + `
	// TODO: 添加更多字段
}

func (h *{{.HandlerName}}) Create(c *gin.Context) {
	var req Create{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	item := model.{{.ModelName}}{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Status:          req.Status,
	}

	if err := h.repo.Create(&item); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.Success(c, item)
}

func (h *{{.HandlerName}}) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "记录不存在")
		return
	}

	var req Create{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	item.Name = req.Name
	item.Status = req.Status

	if err := h.repo.Update(item); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.Success(c, item)
}

func (h *{{.HandlerName}}) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.repo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *{{.HandlerName}}) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	items, total, err := h.repo.List(tenantID, page, pageSize, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, items, total, page, pageSize)
}

func (h *{{.HandlerName}}) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "记录不存在")
		return
	}
	utils.Success(c, item)
}
`
