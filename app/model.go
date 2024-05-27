package app

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
	ID           uint   `gorm:"primary_key" json:"id"`
	Dir          string `gorm:"unique_index" json:"dir"`
	Link         string `json:"link"`
	Remark       string `json:"remark"`
	NotifyTgChat string `json:"notifyTgChat"`
	Season       int    `json:"season"`
	Media        Media  `json:"media"`
	MediaID      uint   `json:"mediaId"`
}
