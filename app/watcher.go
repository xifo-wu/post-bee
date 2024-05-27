package app

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"gorm.io/gorm"
)

type FileEvent struct {
	Name string      // 文件名称
	Op   fsnotify.Op // 文件操作类型
}

func WatcherDirs(db *gorm.DB, watcher *fsnotify.Watcher) {
	var dirs []MonitorDir
	// Get all records
	db.Where("dir <> ?", "").Find(&dirs)

	for _, item := range dirs {
		err := watcher.Add(item.Dir)
		if err != nil {
			log.Println(item.Dir, "无法添加路径到监视器:", err)
		}
	}
}
