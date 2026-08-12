package transport

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/http/handlers"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/Alexander272/Identic/backend/internal/transport/ws"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/limiter"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/Alexander272/Identic/backend/web"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	keycloak     *auth.KeycloakClient
	services     *services.Services
	hub          *ws_hub.Hub
	pricesModule Modules
}

type Modules interface {
	Init(*gin.RouterGroup, *middleware.Middleware)
}

func NewHandler(keycloak *auth.KeycloakClient, services *services.Services, hub *ws_hub.Hub, pricesModule Modules) *Handler {
	return &Handler{
		keycloak:     keycloak,
		services:     services,
		hub:          hub,
		pricesModule: pricesModule,
	}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.New()

	// Отключаем редиректы для SPA
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			Skip: func(c *gin.Context) bool {
				path := c.Request.URL.Path
				if strings.HasPrefix(path, "/api") {
					return false
				}
				return c.Writer.Status() < http.StatusBadRequest
			},
		}),
		gin.CustomRecovery(h.ErrorHandler),
		securityHeaders(),
	)

	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router, conf)
	h.initStatic(router)

	return router
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self' ws: wss:; "+
				"frame-ancestors 'none'; base-uri 'self'; form-action 'self';")
		c.Header("Permissions-Policy",
			"camera=(), microphone=(), geolocation=(), gyroscope=(), "+
				"accelerometer=(), magnetometer=(), usb=(), payment=(), "+
				"display-capture=(), document-domain=()")
		c.Next()
	}
}

func (h *Handler) ErrorHandler(c *gin.Context, origErr any) {
	err := fmt.Errorf("unexpected error: %v", origErr)

	rawStack := string(debug.Stack())                        // 1. Получаем стек в виде байтов
	cleanStack := strings.ReplaceAll(rawStack, "\t", "    ") // 2. Заменяем все табуляции на 4 пробела для красоты
	stackLines := strings.Split(cleanStack, "\n")            // 3. Превращаем в срез строк, разделяя по символу \n

	response.SendError(c, err, gin.H{"PANIC": true, "Stack trace": stackLines})
	debug.PrintStack()
}

func (h *Handler) initAPI(router *gin.Engine, conf *config.Config) {
	middleware := middleware.NewMiddleware(h.services, conf.Auth, h.keycloak)
	handler := handlers.NewHandler(&handlers.Deps{Services: h.services, Conf: conf, Hub: h.hub, Middleware: middleware})
	wsHandler := ws.NewWsHandler(h.hub, conf.Http, h.services)

	api := router.Group("/api")
	api.Use(limiter.Limit(conf.Limiter.RPS, conf.Limiter.Burst, conf.Limiter.TTL))
	handler.Init(api)

	api.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleWS(c)
	})

	h.pricesModule.Init(api, middleware)
}

var appStartTime = time.Now()

const (
	frontendRoot = "frontend"
	indexFile    = "index.html"
	assetsPrefix = "assets/"
)

func (h *Handler) initStatic(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		h.serveFrontend(c, web.Frontend)
	})
}

func (h *Handler) serveFrontend(c *gin.Context, frontend fs.FS) {
	if strings.HasPrefix(c.Request.URL.Path, "/api") {
		c.Status(http.StatusNotFound)
		return
	}

	filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
	if filePath == "" {
		filePath = indexFile
	}
	filePath = path.Clean(filePath)

	var f fs.File
	var err error
	openPath := frontendRoot + "/" + filePath
	encoding := negotiateEncoding(c.Request.Header.Get("Accept-Encoding"))

	if encoding == "br" {
		f, err = frontend.Open(openPath + ".br")
		if err == nil {
			c.Header("Content-Encoding", "br")
		}
	}
	if f == nil && encoding == "gzip" {
		f, err = frontend.Open(openPath + ".gz")
		if err == nil {
			c.Header("Content-Encoding", "gzip")
		}
	}
	if f == nil {
		f, err = frontend.Open(openPath)
		if err != nil {
			// Для недостающих файлов (чанки, шрифты и т.п.) возвращаем 404,
			// а не index.html: иначе браузер получает text/html под видом JS
			// и падает с "MIME type error" / NS_ERROR_CORRUPTED_CONTENT.
			if path.Ext(filePath) != "" {
				c.Status(http.StatusNotFound)
				return
			}
			f, err = frontend.Open(frontendRoot + "/" + indexFile)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			filePath = indexFile
		}
	}
	defer f.Close()

	c.Header("Vary", "Accept-Encoding")

	if strings.HasPrefix(filePath, assetsPrefix) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		c.Header("Cache-Control", "no-cache")
	}

	if ctype := mime.TypeByExtension(path.Ext(filePath)); ctype != "" {
		c.Header("Content-Type", ctype)
	}

	if rs, ok := f.(io.ReadSeeker); ok {
		http.ServeContent(c.Writer, c.Request, path.Base(filePath), appStartTime, rs)
	} else {
		io.Copy(c.Writer, f)
	}
}

// negotiateEncoding parses Accept-Encoding and returns the best compression
// we can offer: "br", "gzip", or "" (no compression).
func negotiateEncoding(header string) string {
	if header == "" {
		return ""
	}

	type pref struct {
		name    string
		quality float64
	}
	var prefs []pref

	for _, field := range strings.Split(header, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		q := 1.0
		name := field

		if idx := strings.Index(field, ";"); idx != -1 {
			name = strings.TrimSpace(field[:idx])
			// find q-value
			if qi := strings.Index(field[idx:], "q="); qi != -1 {
				if parsed, err := strconv.ParseFloat(strings.TrimSpace(field[idx+qi+2:]), 64); err == nil {
					q = parsed
				}
			}
		}

		// identity;q=0 disables content-encoding
		if name == "identity" && q == 0 {
			return ""
		}
		if q > 0 && (name == "br" || name == "gzip" || name == "*") {
			prefs = append(prefs, pref{name, q})
		}
	}

	bestQ := 0.0
	best := ""
	for _, p := range prefs {
		if p.quality > bestQ {
			if p.name == "*" {
				best = "gzip"
			} else {
				best = p.name
			}
			bestQ = p.quality
		} else if p.quality == bestQ && p.name == "br" && best != "br" {
			best = "br"
		}
	}
	return best
}
