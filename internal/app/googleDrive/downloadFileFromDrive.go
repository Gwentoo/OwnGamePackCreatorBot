package googleDrive

import (
	"fmt"
	"io"
	"log"
	"regexp"
)

func ExtractFileIDFromURL(url string) (string, error) {
	re := regexp.MustCompile(`/file/d/([^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("неверный формат URL")
	}
	return matches[1], nil
}

func DownloadFileByID(fileID string) ([]byte, error) {
	res, err := DriveService.Files.Get(fileID).Download()
	if err != nil {
		return nil, fmt.Errorf("не удалось скачать файл: %v", err)
	}

	defer func() {
		if err2 := res.Body.Close(); err2 != nil {
			log.Printf("Ошибка при ответе: %v", err2)
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных: %v", err)
	}

	return data, nil
}
