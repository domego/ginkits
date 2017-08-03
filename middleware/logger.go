// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/domego/ginkits/params"
	"github.com/domego/gokits/log"
	"github.com/gin-gonic/gin"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.Buffer.Write(data)
	return g.out.Write(data)
}

func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter
// By default gin.DefaultWriter = os.Stdout
func Logger() gin.HandlerFunc {
	return LoggerWithWriter()
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(notlogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		w := bufio.NewWriter(c.Writer)
		buff := bytes.Buffer{}
		newWriter := &bufferedWriter{c.Writer, w, buff}

		c.Writer = newWriter

		// You have to manually flush the buffer at the end
		defer func() {
			if strings.Contains(newWriter.Header().Get("Content-Type"), gin.MIMEJSON) {
				log.Tracef("response: %s", string(newWriter.Buffer.Bytes()))
			}

			r := c.Request
			path := r.URL.Path
			// Log only when path is not being skipped
			if _, ok := skip[path]; !ok {
				// Stop timer
				end := time.Now()
				latency := end.Sub(start)

				clientIP := c.ClientIP()
				//		userAgent := c.Request.UserAgent()
				statusCode := c.Writer.Status()

				log.Infof("\"%d\" \"%v\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\"",
					statusCode,
					latency,
					clientIP,
					r.Method,
					r.Host,
					r.RequestURI,
					r.UserAgent(),
				)
			}

			w.Flush()
		}()
		// Process request
		log.Infof("request: %+v", paramkits.GetParams(c))
		c.Next()

	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
