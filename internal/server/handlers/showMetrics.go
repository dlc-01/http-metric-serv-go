package handlers

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"html/template"
)

func (s stor) ShowMetrics(gin *gin.Context) {
	gin.Writer.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")

	page := ""

	metric, err := s.GetAll(gin)
	if err != nil {
		logging.Errorf("cannot get all metric : %s", err)
	}

	for _, n := range metric {
		page += fmt.Sprintf("<h1>	%s</h1>", n)
	}

	tmpl, _ := template.New("data").Parse("<h1>AVAILABLE METRICS</h1>{{range .}}<h3>{{ .}}</h3>{{end}}")
	err = tmpl.Execute(gin.Writer, metric)
	if err != nil {
		logging.Errorf("cannot execute template : %s", err)
		return
	}

	gin.Writer.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
}
