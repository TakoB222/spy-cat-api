package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

// Request logger middleware
type requestLogger struct {
	rc io.ReadCloser
	w  io.Writer
}

func (rc *requestLogger) Read(p []byte) (n int, err error) {
	n, err = rc.rc.Read(p)
	if n > 0 {
		if n, err := rc.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}

func (rc *requestLogger) Close() error {
	return rc.rc.Close()
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := bytes.Buffer{}
		body := &requestLogger{c.Request.Body, &buf}
		c.Request.Body = body
		c.Next()
		if buf.Len() > 0 {
			logrus.Println(map[string]string{
				"request": buf.String(),
			})
		}
	}
}

// Response logger middleware
type responseLogger struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseLogger) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := &responseLogger{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = buf
		c.Next()
		if buf.body.Len() > 0 {
			logrus.Println(map[string]string{
				"response": buf.body.String(),
			})
		}
	}
}
