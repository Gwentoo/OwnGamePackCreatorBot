package storage

import (
	"OwnGamePack/internal/app/googleDrive"
	"OwnGamePack/internal/structs"
	"fmt"
	"log"
)

func (s *Storage) GetPack(userID int64, name string) (string, error) {
	var link string
	err := s.db.QueryRow(`
	SELECT pack FROM packs
	WHERE user_id = $1 AND pack_name = $2`, userID, name).Scan(&link)

	if err != nil {
		return "", err
	}
	return link, err
}

func (s *Storage) GetPacksName(userID int64) (map[string]bool, error) {
	names := make(map[string]bool)
	var name string
	var public bool
	rows, err := s.db.Query(`
					SELECT pack_name, public FROM packs
					WHERE user_id = $1`,
		userID,
	)
	defer func() {
		if err1 := rows.Close(); err1 != nil {
			log.Printf("Ошибка при ответе: %v", err1)
		}
	}()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err2 := rows.Scan(&name, &public); err2 != nil {
			return nil, fmt.Errorf("scan error: %v", err2)
		}
		names[name] = public
	}
	if err3 := rows.Err(); err3 != nil {
		return nil, fmt.Errorf("rows error: %v", err3)
	}
	return names, nil
}

func (s *Storage) SavePack(pack *structs.Pack, public bool) error {

	link, err := googleDrive.UploadFileToDrive(pack)
	if err != nil {
		return err
	}

	_, err2 := s.db.Exec(
		`INSERT INTO packs (user_id, pack_id, pack, created_at, updated_at, public, pack_name) 
        VALUES ($1, $2, $3, NOW(), NOW(), $4, $5) 
        ON CONFLICT (pack_id)
        DO UPDATE SET 
        	pack = EXCLUDED.pack,
        	updated_at = NOW(),
        	public = excluded.public`,
		pack.UserID,
		pack.PackID,
		link,
		public,
		pack.PackName,
	)
	if err2 != nil {
		err3 := googleDrive.DeleteFromGoogleDrive(link)
		if err3 != nil {
			return err3
		}
		return err2
	}

	return nil
}
