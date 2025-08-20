package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	// tạo đường dẫn file ghi log
	logPath := "../../internal/logs/http.log"

	// sử dụng package zerolog
	// Sử dụng lumberjack để tạo file log
	// xóa file khi đạt dung lượng cho phép
	// tránh file log quá nặng
	logger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1, // megabytes
		MaxBackups: 5, // lưu file backup => nén file, tối đa là 5
		// nếu tạo thêm thì sẽ xóa file cũ nhất và tạo file mới
		MaxAge:    5,    //5 ngày sẽ xóa file
		Compress:  true, // disabled by default => nén file
		LocalTime: true,
	}).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		start := time.Now()

		// Các loại content-type:
		// application/json
		// application/x-www-form-urlencoded
		// multipart/form-data

		contentType := ctx.GetHeader("Content-Type")

		requestBody := make(map[string]any)
		var formFiles []map[string]any

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := ctx.Request.ParseMultipartForm(32 << 20); err == nil && ctx.Request.MultipartForm != nil {
				// for value
				for key, vals := range ctx.Request.MultipartForm.Value {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}

				// for file
				for field, files := range ctx.Request.MultipartForm.File {
					for _, f := range files {
						formFiles = append(formFiles, map[string]any{
							"filed":        field,
							"filename":     f.Filename,
							"size":         formatFileSize(f.Size),
							"content-type": f.Header.Get("Content-Type"),
						})
					}
				}
				if len(formFiles) > 0 {
					requestBody["form_files"] = formFiles
				}
			}
			log.Println("multipart/form-data")
		} else {
			// multipart/form-data
			// lấy tất cả thông tin gửi lên từ body
			body, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to read request body")
			}

			// phải gán lại giá trị cho body
			// vì nếu ReadAll thì body lúc này ko còn giá trị gì nữa
			// nên sau khi qua handler sẽ bị lỗi
			// vì ko còn body để đọc
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if strings.HasPrefix(contentType, "application/json") {
				// application/json
				_ = json.Unmarshal(body, &requestBody)
			} else {
				// application/x-www-form-urlencoded
				values, _ := url.ParseQuery(string(body))
				for key, vals := range values {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}
			}
		}

		customWriter := &CustomResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}

		ctx.Writer = customWriter

		// trước ctx.Next là bắt đầu
		// sau ctx.Next là kết thúc

		ctx.Next()

		duration := time.Since(start)

		statusCode := ctx.Writer.Status()

		responseContentType := ctx.Writer.Header().Get("Content-Type")
		responseBodyRaw := customWriter.body.String()
		var responseBodyParsed interface{}
		if strings.HasPrefix(responseContentType, "image/") {
			// image
			responseBodyParsed = "[BINARY DATA]"
		} else if strings.HasPrefix(responseContentType, "application/json") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "{") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "[") {
			// json
			if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParsed); err != nil {
				responseBodyParsed = responseBodyRaw
			}
		} else {
			responseBodyParsed = responseBodyRaw
		}
		log.Printf("%s", responseBodyRaw)

		logEvent := logger.Info()
		if statusCode >= 500 {
			logEvent = logger.Error()
		} else if statusCode >= 400 {
			logEvent = logger.Warn()
		}

		// log.Str => có nghĩa là chuỗi
		// log Int64 => số nguyên 64 bit
		logEvent.
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.RawQuery).
			Str("query", ctx.Request.Method).
			Str("client_ip", ctx.ClientIP()).
			// user_agent: sử dụng trình duyệt nào
			Str("user_agent", ctx.Request.UserAgent()).
			// referer : khi người dùng click link từ fb hoặc zalo...
			// thì nó sẽ chuyển hướng vào api
			Str("referer", ctx.Request.Referer()).
			Str("protocol", ctx.Request.Proto).
			Str("host", ctx.Request.Host).
			// remote_address: lấy địa chỉ nếu user đó đi qua proxy
			Str("remote_addr", ctx.Request.RemoteAddr).
			Str("request_uri", ctx.Request.RequestURI).
			Int64("content_length", ctx.Request.ContentLength).
			Interface("header", ctx.Request.Header).
			Interface("request_body", requestBody).
			Interface("response_body", responseBodyParsed).
			Int("status_code", statusCode).
			Int64("duration_ms", duration.Microseconds()).
			// để có thể ghi lại log thì cần có msg
			Msg("HTTP request log")
	}
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d byte", size)
	}
}
