package types

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

// Define the structures based on your YAML file

type Route struct {
	Name  string `yaml:"name"`
	Paths string `yaml:"paths"`
}

type Service struct {
	Host       string  `yaml:"host"`
	Name       string  `yaml:"name"`
	Port       string  `yaml:"port"`
	Protocol   string  `yaml:"protocol"`
	Premession string  `yaml:"premession"`
	Routes     []Route `yaml:"routes"`
}
type RateLimiterMemoryStoreConfig struct {
	Rate         int           `yaml:"rate"`
	Burst        int           `yaml:"burst"`
	ExpiresIn    time.Duration `yaml:"expires_in"`
	BlacklistIPs []string      `yaml:"blacklist_ips"`
}
type Config struct {
	FormatVersion          string                       `yaml:"_format_version"`
	Transform              bool                         `yaml:"_transform"`
	RateLimiterMemoryStore RateLimiterMemoryStoreConfig `yaml:"rate_limiter_memory_store"`
	Services               []Service                    `yaml:"services"`
}

func (s *Service) GetPremessionType() echo.MiddlewareFunc {

	switch s.Premession {
	case "require_record_auth":
		return apis.RequireRecordAuth()
	case "require_guest_only":
		return apis.RequireGuestOnly()
	case "require_admin_auth":
		return apis.RequireAdminAuth()
	case "require_admin_or_record_auth":
		return apis.RequireAdminOrRecordAuth()

	default:
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}

		}

	}

}
