package subnet

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// CheckSubnet — middleware that use for checking ip for http request "X-Real-IP".
func CheckSubnet(subnet string) gin.HandlerFunc {
	return func(gin *gin.Context) {
		if subnet != "" {
			realIP := gin.Request.Header.Get("X-Real-IP")
			_, ipnet, err := net.ParseCIDR(subnet)
			if err != nil {
				logging.Errorf("can't parse ip %s", err)
				gin.AbortWithStatus(http.StatusBadRequest)
			}

			if !ipnet.Contains(net.ParseIP(realIP)) {
				gin.AbortWithStatus(http.StatusForbidden)
			}
		}

		gin.Next()
	}
}
