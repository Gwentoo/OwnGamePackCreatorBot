package googleDrive

import (
	"OwnGamePack/internal/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/api/drive/v3"
)

func UploadFileToDrive(pack *structs.Pack) (string, error) {

	jsonData, err := json.Marshal(&pack)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации: %v", err)
	}

	fileName := fmt.Sprintf("%d.json", pack.PackID)

	fileID, exists, err := FileExistsInFolder(DriveService, fileName, "1IM10QuHDyHTH1xgUd5yEz4WRJl82O-Rc")
	if err != nil {
		return "", fmt.Errorf("ошибка проверки файла: %v", err)
	}
	if exists {
		_, err = DriveService.Files.Update(fileID, &drive.File{
			Name: fileName,
		}).
			Media(bytes.NewReader(jsonData)).
			Do()

		if err != nil {
			return "", fmt.Errorf("ошибка обновления: %v", err)
		}

		return fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID), nil

	} else {
		file := &drive.File{
			Name:     fileName,
			MimeType: "application/json",
			Parents:  []string{folderID},
		}

		res, err := DriveService.Files.Create(file).
			Media(bytes.NewReader(jsonData)).
			Do()

		if err != nil {
			return "", fmt.Errorf("ошибка загрузки: %v", err)
		}

		return fmt.Sprintf("https://drive.google.com/file/d/%s/view", res.Id), nil
	}

}
