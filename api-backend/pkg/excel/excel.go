package excel

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ColumnDef 列定义
type ColumnDef struct {
	Header string // 表头名称
	Field  string // 结构体字段名
	Width  float64
}

// Export 通用导出函数，将数据导出为 Excel 并写入 HTTP 响应
func Export(c *gin.Context, filename string, sheetName string, columns []ColumnDef, data interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := sheetName
	if sheet == "" {
		sheet = "Sheet1"
	}
	f.SetSheetName("Sheet1", sheet)

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "D9D9D9", Style: 1},
			{Type: "right", Color: "D9D9D9", Style: 1},
			{Type: "top", Color: "D9D9D9", Style: 1},
			{Type: "bottom", Color: "D9D9D9", Style: 1},
		},
	})

	// 写入表头
	for i, col := range columns {
		cell := fmt.Sprintf("%s1", colName(i))
		f.SetCellValue(sheet, cell, col.Header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		if col.Width > 0 {
			f.SetColWidth(sheet, colName(i), colName(i), col.Width)
		} else {
			f.SetColWidth(sheet, colName(i), colName(i), 15)
		}
	}

	// 写入数据
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	for row := 0; row < dataValue.Len(); row++ {
		item := dataValue.Index(row)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		for col, colDef := range columns {
			cell := fmt.Sprintf("%s%d", colName(col), row+2)
			val := getFieldValue(item, colDef.Field)
			f.SetCellValue(sheet, cell, val)
		}
	}

	// 设置响应头
	encodedName := url.QueryEscape(filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	return f.Write(c.Writer)
}

// Import 通用导入函数，从上传的 Excel 文件读取数据
func Import(c *gin.Context, fieldName string, columns []ColumnDef) ([]map[string]string, error) {
	file, _, err := c.Request.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %v", err)
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("解析Excel失败: %v", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel文件没有数据行")
	}

	// 建立表头到字段的映射
	headerMap := make(map[int]string)
	for i, header := range rows[0] {
		for _, col := range columns {
			if col.Header == header {
				headerMap[i] = col.Field
				break
			}
		}
	}

	// 读取数据行
	var result []map[string]string
	for _, row := range rows[1:] {
		record := make(map[string]string)
		for i, cell := range row {
			if field, ok := headerMap[i]; ok {
				record[field] = cell
			}
		}
		if len(record) > 0 {
			result = append(result, record)
		}
	}

	return result, nil
}

// ExportTemplate 导出空模板（仅表头）
func ExportTemplate(c *gin.Context, filename string, sheetName string, columns []ColumnDef) error {
	return Export(c, filename, sheetName, columns, []struct{}{})
}

// colName 将列索引转换为 Excel 列名 (0->A, 1->B, ..., 25->Z, 26->AA)
func colName(index int) string {
	name := ""
	for {
		name = string(rune('A'+index%26)) + name
		index = index/26 - 1
		if index < 0 {
			break
		}
	}
	return name
}

// getFieldValue 通过字段名获取结构体字段值
func getFieldValue(v reflect.Value, field string) interface{} {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	f := v.FieldByName(field)
	if !f.IsValid() {
		// 尝试通过 json tag 匹配
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")
			if tag == field || (len(tag) > 0 && tag[:len(field)] == field) {
				f = v.Field(i)
				break
			}
		}
	}

	if !f.IsValid() {
		return ""
	}

	// 处理特殊类型
	switch f.Interface().(type) {
	case time.Time:
		t := f.Interface().(time.Time)
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	case *time.Time:
		t := f.Interface().(*time.Time)
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	}

	return f.Interface()
}

// WriteResponse 简化的响应写入（用于自定义导出场景）
func WriteResponse(c *gin.Context, filename string, f *excelize.File) {
	encodedName := url.QueryEscape(filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Status(http.StatusOK)
	f.Write(c.Writer)
}
