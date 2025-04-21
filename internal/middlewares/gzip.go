package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipReader - обертка для gzip.Reader
type GzipReader struct {
	io.ReadCloser
	reader *gzip.Reader
}

func (g *GzipReader) Read(p []byte) (int, error) {
	return g.reader.Read(p)
}

// GunzipMiddleware - middleware для распаковки запросов
func (m Middleware) GunzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gz.Close()

			c.Request.Body = &GzipReader{c.Request.Body, gz}
		}
		c.Next()
	}
}

// GzipWriter - обертка для gzip.Writer
type GzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *GzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// GzipMiddleware - middleware для сжатия ответов
func (m Middleware) GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()

			c.Writer = &GzipWriter{c.Writer, gz}
			c.Header("Content-Encoding", "gzip")
		}
		c.Next()
	}
}
