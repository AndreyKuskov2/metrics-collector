package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusInternalServerError)
				return
			}

			defer func(gzReader *gzip.Reader) {
				err := gzReader.Close()
				if err != nil {
					http.Error(w, "Error close gzReader", http.StatusInternalServerError)
				}
			}(gzReader)
			r.Body = gzReader

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")

				gw := gzip.NewWriter(w)
				defer func(gw *gzip.Writer) {
					err := gw.Close()
					if err != nil {
						http.Error(w, "Error close gzWriter", http.StatusInternalServerError)
					}
				}(gw)

				gzipWriter := gzipResponseWriter{Writer: gw, ResponseWriter: w}
				h.ServeHTTP(gzipWriter, r)
				return
			}
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && strings.Contains(r.Header.Get("Accept"), "text/html") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Type", "text/html")
			gw := gzip.NewWriter(w)
			defer func(gw *gzip.Writer) {
				err := gw.Close()
				if err != nil {
					http.Error(w, "Error close gzWriter", http.StatusInternalServerError)
				}
			}(gw)

			gzipWriter := gzipResponseWriter{Writer: gw, ResponseWriter: w}
			h.ServeHTTP(gzipWriter, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
