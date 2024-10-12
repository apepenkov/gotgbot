package images_manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
	"time"
)

// In MTPROTO files would become invalid after some time.
// I am not sure, how does it work in HTTP api, but better be safe than sorry.
// So, I will reupload all images from time to time.

type FileName string

var allImages = []FileName{}

type ImagesManager struct {
	CachedImages  map[string]*tgbotapi.FileID
	lock          *sync.RWMutex
	bot           *tgbotapi.BotAPI
	uploadChannel int64
	directory     string
}

func (i *ImagesManager) GetImage(name FileName) tgbotapi.RequestFileData {
	i.lock.RLock()
	defer i.lock.RUnlock()

	if fileID, ok := i.CachedImages[string(name)]; ok {
		return fileID
	}

	return tgbotapi.FilePath(i.directory + string(name))
}

func (i *ImagesManager) uploadImage(name FileName) (tgbotapi.FileID, error) {
	sentMessage, err := i.bot.Send(tgbotapi.NewPhoto(i.uploadChannel, tgbotapi.FilePath(i.directory+string(name))))
	if err != nil {
		return "", err
	}
	return tgbotapi.FileID(sentMessage.Photo[0].FileID), nil
}

func (i *ImagesManager) uploadAllImages() error {
	for _, img := range allImages {
		fileID, err := i.uploadImage(img)
		if err != nil {
			return err
		}
		i.lock.Lock()
		i.CachedImages[string(img)] = &fileID
		i.lock.Unlock()
		time.Sleep(1 * time.Second)
	}
	return nil
}
func (i *ImagesManager) uploader() {
	for {
		err := i.uploadAllImages()
		if err != nil {
			log.Printf("Error uploading images: %v", err)
		}
		time.Sleep(6 * time.Hour)
	}
}

func NewUploader(bot *tgbotapi.BotAPI, uploadChannel int64, directory string) *ImagesManager {
	im := &ImagesManager{
		CachedImages:  make(map[string]*tgbotapi.FileID),
		lock:          &sync.RWMutex{},
		bot:           bot,
		uploadChannel: uploadChannel,
		directory:     directory,
	}
	go im.uploader()
	return im
}
