package googleDrive

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	serviceAccountKey = "internal/app/googleDrive/service-account.json"
	folderID          = "1IM10QuHDyHTH1xgUd5yEz4WRJl82O-Rc"
)

var (
	DriveService *drive.Service
)

func InitGoogleDrive() {
	ctx := context.Background()

	key, err := os.ReadFile(serviceAccountKey)
	if err != nil {
		log.Fatalf("Ошибка чтения ключа: %v", err)
	}

	credentials, err := google.CredentialsFromJSON(ctx, key, drive.DriveScope)
	if err != nil {
		log.Fatalf("Ошибка аутентификации: %v", err)
	}

	DriveService, err = drive.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}
}

func FileExistsInFolder(service *drive.Service, fileName string, folderID string) (string, bool, error) {
	query := fmt.Sprintf("name='%s' and trashed=false", fileName)
	if folderID != "" {
		query += fmt.Sprintf(" and '%s' in parents", folderID)
	}

	res, err := service.Files.List().
		Q(query).
		Spaces("drive").
		Fields("files(id, name)").
		PageSize(1).
		Do()

	if err != nil {
		return "", false, fmt.Errorf("ошибка поиска файла: %v", err)
	}

	if len(res.Files) > 0 {
		return res.Files[0].Id, true, nil
	}
	return "", false, nil
}
