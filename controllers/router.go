package controllers

import (
	"net/http"

	"goweb/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"goweb/container"
)

const (
	StaticDir    = "static"
	StaticPrefix = "files"
)

func BuildRouter(c *container.Container) {
	// Static files with proper cache control
	// funcmap.File() should be used in templates to append a cache key to the URL in order to break cache
	// after each server restart
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.MaxAge.StaticFile)).
		Static(StaticPrefix, StaticDir)

	// Middleware
	g := c.Web.Group("",
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.RequestID(),
		echomw.Recover(),
		echomw.Gzip(),
		echomw.Logger(),
		middleware.LogRequestID(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
		middleware.PageCache(c.Cache),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup: "form:csrf",
		}),
	)

	// Base controller
	ctr := NewController(c)

	// Error handler
	err := Error{Controller: ctr}
	c.Web.HTTPErrorHandler = err.Get

	// Routes
	navRoutes(g, ctr)
	userRoutes(g, ctr)
}

func navRoutes(g *echo.Group, ctr Controller) {
	home := Home{Controller: ctr}
	g.GET("/", home.Get).Name = "home"

	about := About{Controller: ctr}
	g.GET("/about", about.Get).Name = "about"

	contact := Contact{Controller: ctr}
	g.GET("/contact", contact.Get).Name = "contact"
	g.POST("/contact", contact.Post).Name = "contact.post"
}

func userRoutes(g *echo.Group, ctr Controller) {
	login := Login{Controller: ctr}
	g.GET("/user/login", login.Get).Name = "login"
	g.POST("/user/login", login.Post).Name = "login.post"

	register := Register{Controller: ctr}
	g.GET("/user/register", register.Get).Name = "register"
	g.POST("/user/register", register.Post).Name = "register.post"
}