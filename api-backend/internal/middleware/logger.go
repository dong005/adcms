package middleware

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"adcms/pkg/logcfg"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var reqBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			reqBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		duration := time.Since(startTime).Milliseconds()

		userID := GetUserID(c)
		tenantID := GetTenantID(c)

		if userID > 0 && logcfg.IsLogEnabled("log_operation_enabled") {
			log := model.OperationLog{
				TenantID:  tenantID,
				UserID:    userID,
				Module:    "",
				Action:    "",
				Method:    c.Request.Method,
				Path:      c.Request.URL.Path,
				Params:    reqBody,
				Response:  blw.body.String(),
				IP:        c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Duration:  duration,
				CreatedAt: time.Now(),
			}

			go database.DB.Create(&log)
		}
	}
}
