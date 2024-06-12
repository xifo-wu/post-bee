package app

import (
	"log"

	"github.com/mymmrac/telego"
	"gorm.io/gorm"
)

type Schedule struct {
	DB  *gorm.DB
	Bot *telego.Bot
}

func (s *Schedule) Notify() {
	var notificationTasks []NotificationTask

	s.DB.Where("status = ?", 0).Find(&notificationTasks)

	for _, notificationTask := range notificationTasks {
		var notificationTask = notificationTask
		log.Println(notificationTask)
		NotifyTelegram(s.DB, notificationTask.Dir, s.Bot)
		s.DB.Model(&notificationTask).Update("status", 1)
	}
}
