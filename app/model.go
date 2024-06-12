package app

import "time"

type Media struct {
	ID          uint         `gorm:"primary_key" json:"id"`
	Name        string       `json:"name"`
	PosterUrl   string       `json:"posterUrl"`
	Year        int          `json:"year"`
	Tags        string       `json:"tags"`
	Desc        string       `json:"desc"`
	MonitorDirs []MonitorDir `json:"monitorDirs"`
}

type MonitorDir struct {
	ID                uint               `gorm:"primary_key" json:"id"`
	Dir               string             `gorm:"unique_index" json:"dir"`
	Link              string             `json:"link"`
	Remark            string             `json:"remark"`
	NotifyTgChat      string             `json:"notifyTgChat"`
	Season            int                `json:"season"`
	Media             Media              `json:"media"`
	MediaID           uint               `json:"mediaId"`
	NotificationTasks []NotificationTask `json:"notificationTasks"`
}

type NotificationTask struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Dir          string    `gorm:"unique_index" json:"dir"`
	Status       int       `gorm:"status"`
	MonitorDirID uint      `json:"monitorDirId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
