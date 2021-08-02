package echox

import (
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	privateNet = make([]*net.IPNet, 0, 4)
)

func init() {
	// refer: https://en.wikipedia.org/wiki/Private_network
	privateCidr := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fd00::/8",
	}
	for _, it := range privateCidr {
		_, block, _ := net.ParseCIDR(it)
		privateNet = append(privateNet, block)
	}
}

// IntranetOnly middleware for only allow intranet access
func IntranetOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := net.ParseIP(c.RealIP())
		if ip.IsLoopback() {
			return next(c)
		}

		for _, b := range privateNet {
			if b.Contains(ip) {
				return next(c)
			}
		}

		return c.String(http.StatusForbidden, fmt.Sprint("Forbidden, ip not allowed ", ip))
	}
}
