package utils

import (
	"fmt"
	"os"
)

func CheckAndRemoveDir(dirPath string) error {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("ошибка проверки папки: %w", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s не является папкой", dirPath)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка удаления папки: %w", err)
	}

	return nil
}
