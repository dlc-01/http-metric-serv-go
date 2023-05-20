package html

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"html/template"
)

func ShowMetrics(gin *gin.Context) {
	gin.Writer.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")

	page := ""
	for _, n := range storage.GetAll() {
		page += fmt.Sprintf("<h1>	%s</h1>", n)
	}
	tmpl, _ := template.New("data").Parse("<h1>AVAILABLE METRICS</h1>{{range .}}<h3>{{ .}}</h3>{{end}}")
	if err := tmpl.Execute(gin.Writer, storage.GetAll()); err != nil {
		return
	}
	gin.Writer.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
}
