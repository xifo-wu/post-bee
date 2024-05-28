package app

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ApiHandle struct {
	DB      *gorm.DB
	Watcher *fsnotify.Watcher
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *ApiHandle) Login(c echo.Context) error {
	var data LoginBody
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	// Throws unauthorized error
	if data.Username != viper.GetString("USER") || data.Password != viper.GetString("PASS") {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		viper.GetString("USER"),
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("LoginDuration")))),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("USER")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

var (
	lock = sync.Mutex{}
)

func (h *ApiHandle) CreateMedia(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	var media Media
	err := c.Bind(&media)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	h.DB.Create(&media)
	if media.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": "创建失败"})
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": media})
}

func (h *ApiHandle) UpdateMedia(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var media Media
	media.ID = uint(id)
	err := c.Bind(&media)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	err = h.DB.Save(&media).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": media})
}

func (h *ApiHandle) ListMedia(c echo.Context) error {
	var media []Media
	h.DB.Find(&media)

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": media})
}

func (h *ApiHandle) DeleteMedia(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	h.DB.Delete(&Media{}, id)
	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

func (h *ApiHandle) CreateMonitorDir(c echo.Context) error {
	var monitorDir MonitorDir
	err := c.Bind(&monitorDir)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	h.DB.Create(&monitorDir)
	if monitorDir.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": "创建失败"})
	}

	err = h.Watcher.Add(monitorDir.Dir)
	if err != nil {
		log.Println(monitorDir.Dir, "无法添加路径到监视器:", err)
	}
	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": monitorDir})
}

func (h *ApiHandle) UpdateMonitorDir(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var monitorDir MonitorDir
	monitorDir.ID = uint(id)
	err := c.Bind(&monitorDir)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	err = h.DB.Save(&monitorDir).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": monitorDir})
}

func (h *ApiHandle) ListMonitorDir(c echo.Context) error {
	var monitorDir []MonitorDir
	orm := h.DB
	mediaID := c.QueryParam("mediaId")
	if mediaID != "" {
		orm = orm.Where("media_id = ?", mediaID)
	}
	orm.Find(&monitorDir)

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": monitorDir})
}

func (h *ApiHandle) DeleteMonitorDir(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	h.DB.Delete(&MonitorDir{}, id)
	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
