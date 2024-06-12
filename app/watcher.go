package app

import (
	"log"
	"os/exec"

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
		// err := filepath.Walk(item.Dir, func(path string, info os.FileInfo, err error) error {
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if info.IsDir() {
		// 		err = watcher.Add(path)
		// 		if err != nil {
		// 			log.Println("无法添加路径到监视器:", err)
		// 		}
		// 	}
		// 	return nil
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }
		exec.Command("ls " + item.Dir)
		err := watcher.Add(item.Dir)
		if err != nil {
			log.Println(item.Dir, "无法添加路径到监视器:", err)
		} else {
			log.Println(item.Dir, "添加路径到监视器成功")
		}
	}
}
