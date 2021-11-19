package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"

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

func ResponseLogger(c *gin.Context) {
	blw := &bodyLogWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: c.Writer,
	}
	c.Writer = blw
	c.Next()

	dst := &bytes.Buffer{}
	json.Indent(dst, blw.body.Bytes(), "", "  ")

	fmt.Println("Response body: " + string(dst.String()))
}
