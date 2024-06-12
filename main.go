package main

import (
	"log"
	"net/http"
	"os"
	"post-bee/app"

	"github.com/fsnotify/fsnotify"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mymmrac/telego"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func main() {
	app.InitConfig()
	db := app.InitDB()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	app.WatcherDirs(db, watcher)

	botToken := viper.GetString("TELEGRAM_BOT_TOKEN")
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	// 启动一个goroutine来处理文件事件
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						// 处理错误
						log.Println("出现错误", err)
						return
					}
					if fileInfo.IsDir() {
						watcher.Add(event.Name)
						return
					}

					log.Println("创建文件：", event.Name)
					app.NotifyCreat(db, event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("错误：", err)
			}
		}
	}()

	api := app.ApiHandle{DB: db, Watcher: watcher}

	scheduleHandle := app.Schedule{DB: db, Bot: bot}
	crontab := cron.New()
	crontab.AddFunc("* * * * *", scheduleHandle.Notify)
	crontab.Start()

	e.Static("/", "dist")
	e.File("/", "dist/index.html")

	e.POST("/api/login", api.Login)

	auth := e.Group("")
	auth.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString("USER")),
	}))

	auth.POST("/api/media", api.CreateMedia)
	auth.GET("/api/media", api.ListMedia)
	auth.PUT("/api/media/:id", api.UpdateMedia)
	auth.DELETE("/api/media/:id", api.DeleteMedia)
	auth.POST("/api/monitorDirs", api.CreateMonitorDir)
	auth.GET("/api/monitorDirs", api.ListMonitorDir)
	auth.PUT("/api/monitorDirs/:id", api.UpdateMonitorDir)
	auth.DELETE("/api/monitorDirs/:id", api.DeleteMonitorDir)

	// 启动服务器
	e.Start(":" + viper.GetString("PORT"))
}
