package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/nssteinbrenner/anitogo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func isVideoFile(fileName string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".rmvb", ".3gp", ".mpg", ".m4v", ".vob"}
	ext := strings.ToLower(fileName[strings.LastIndex(fileName, "."):])
	for _, extension := range videoExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

func NotifyCreat(db *gorm.DB, path string) {
	dir := filepath.Dir(path)

	var monitorDir MonitorDir

	db.Preload("Media").Where("dir = ?", dir).First(&monitorDir)

	if monitorDir.ID == 0 {
		return
	}

	// 如果任务已经存在且状态为0，则不创建新的任务
	var monitorDirCount int64
	db.Where("dir = ? AND status = ?", path, 0).Count(&monitorDirCount)
	if monitorDirCount != 0 {
		return
	}

	// 创建一条新的带处理任务
	db.Create(&NotificationTask{Dir: path, Status: 0, MonitorDirID: monitorDir.ID})
}

func NotifyTelegram(db *gorm.DB, path string, bot *telego.Bot) {
	log.Println("创建事件：", path)
	// 获取文件名
	filename := filepath.Base(path)

	// 只触发视频文件的创建
	isVideo := isVideoFile(filename)
	if !isVideo {
		return
	}

	parsed := anitogo.Parse(filename, anitogo.DefaultOptions)

	// Accessing the elements directly
	fmt.Println("Anime Title:", parsed.AnimeTitle)
	fmt.Println("Anime Year:", parsed.AnimeYear)
	fmt.Println("Episode Number:", parsed.EpisodeNumber)
	fmt.Println("Release Group:", parsed.ReleaseGroup)
	fmt.Println("File Checksum:", parsed.FileChecksum)

	dir := filepath.Dir(path)

	var monitorDir MonitorDir

	db.Preload("Media").Where("dir = ?", dir).First(&monitorDir)

	if monitorDir.ID == 0 {
		return
	}

	var err error
	season := monitorDir.Season
	if season == 0 {
		season, err = strconv.Atoi(parsed.AnimeSeason[0])
		if err != nil {
			season = 0
		}
	}

	episodeNumber, err := strconv.Atoi(parsed.EpisodeNumber[0])
	if err != nil {
		episodeNumber = 0
	}

	caption := fmt.Sprintf("名称：<b>%s (%d) S%02dE%02d %s</b>  \n\n", monitorDir.Media.Name, monitorDir.Media.Year, season, episodeNumber, monitorDir.Remark)
	desc := fmt.Sprintf("简介：%s \n\n", monitorDir.Media.Desc)
	tags := fmt.Sprintf("标签：%s \n\n", monitorDir.Media.Tags)
	link := fmt.Sprintf("链接：\n %s \n\n", monitorDir.Link)

	adStr := viper.GetString("TELEGRAM_BOT_AD")
	caption = caption + desc + tags + link + "\n" + adStr
	params := telego.SendPhotoParams{
		ChatID: telego.ChatID{Username: monitorDir.NotifyTgChat},
		Photo: telego.InputFile{
			URL: monitorDir.Media.PosterUrl,
		},
		Caption:   caption,
		ParseMode: "html",
	}

	go func() {
		bot.SendPhoto(&params)
	}()
}
