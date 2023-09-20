package controller

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

type HandleController struct {
	CommonInfo CommonInfo
	Version    string
	ConfigFile string
	TokensFile string
}

func NewHandleController(config *HandleController) *HandleController {
	return config
}

func (c *HandleController) Register(rootDir string, engine *gin.Engine) {
	assets := filepath.Join(rootDir, "assets")
	_, err := os.Stat(assets)
	if err != nil && !os.IsExist(err) {
		assets = "./assets"
	}

	engine.Delims("${", "}")
	engine.LoadHTMLGlob(filepath.Join(assets, "templates/*"))
	engine.Static("/static", filepath.Join(assets, "static"))
	engine.GET("/lang.json", c.MakeLangFunc())
	engine.GET(LoginUrl, c.MakeLoginFunc())
	engine.POST(LoginUrl, c.MakeLoginFunc())
	engine.GET(LogoutUrl, c.MakeLogoutFunc())

	var group *gin.RouterGroup
	if len(c.CommonInfo.AdminUser) != 0 {
		group = engine.Group("/", c.BasicAuth())
	} else {
		group = engine.Group("/")
	}
	group.GET("/", c.MakeIndexFunc())
	group.POST("/add", c.MakeAddProxyFunc())
	group.POST("/update", c.MakeUpdateProxyFunc())
	group.POST("/remove", c.MakeRemoveProxyFunc())
	group.GET("/proxy/*serverApi", c.MakeProxyFunc())
}
