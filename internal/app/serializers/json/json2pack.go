package json

import (
	"OwnGamePack/internal/structs"
	"encoding/json"
	"fmt"
)

func DataToPack(data []byte) (*structs.Pack, error) {

	var pack structs.Pack
	if err := json.Unmarshal(data, &pack); err != nil {
		return nil, fmt.Errorf("ошибка парсинга json: %w", err)
	}
	return &pack, nil
}
