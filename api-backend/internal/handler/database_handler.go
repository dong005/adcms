package handler

import (
	"adcms/pkg/database"
	"adcms/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DatabaseHandler struct{}

func NewDatabaseHandler() *DatabaseHandler {
	return &DatabaseHandler{}
}

type TableInfo struct {
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Rows      int64  `json:"rows"`
	DataSize  string `json:"data_size"`
	IndexSize string `json:"index_size"`
	Comment   string `json:"comment"`
}

func (h *DatabaseHandler) Tables(c *gin.Context) {
	var tables []TableInfo
	database.DB.Raw(`
		SELECT 
			TABLE_NAME as name,
			ENGINE as engine,
			TABLE_ROWS as 'rows',
			CONCAT(ROUND(DATA_LENGTH/1024, 2), ' KB') as data_size,
			CONCAT(ROUND(INDEX_LENGTH/1024, 2), ' KB') as index_size,
			TABLE_COMMENT as comment
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME
	`).Scan(&tables)

	utils.Success(c, tables)
}

type ColumnInfo struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Nullable  string  `json:"nullable"`
	Key       string  `json:"key"`
	Default   *string `json:"default"`
	Comment   string  `json:"comment"`
}

func (h *DatabaseHandler) Columns(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		utils.BadRequest(c, "参数错误")
		return
	}

	var columns []ColumnInfo
	database.DB.Raw(`
		SELECT 
			COLUMN_NAME as name,
			COLUMN_TYPE as type,
			IS_NULLABLE as nullable,
			COLUMN_KEY as 'key',
			COLUMN_DEFAULT as 'default',
			COLUMN_COMMENT as comment
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`, tableName).Scan(&columns)

	utils.Success(c, columns)
}
