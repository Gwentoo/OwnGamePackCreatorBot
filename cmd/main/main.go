package main

import (
	"OwnGamePack/config"
	"OwnGamePack/internal/app/bot"
	"OwnGamePack/internal/app/googleDrive"
	"OwnGamePack/internal/app/handlers"
	"OwnGamePack/internal/storage"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	//Подключение к бд
	cfg, err := config.LoadConfig()

	handlers.StorageDB, err = storage.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	googleDrive.InitGoogleDrive()

	defer func() {
		if err2 := handlers.StorageDB.Close(); err2 != nil {
			log.Printf("Ошибка при ответе: %v", err2)
		}
	}()
	//Подключение бота
	b, err := bot.NewBot(cfg, handlers.StorageDB)
	if err != nil {
		log.Fatal("Bot init error:", err)
	}

	//pack, err := json.JsonToPack("D:\\ProgsGo\\OwnGamePack\\3794095286597811567.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Конвертируем в xml
	//xmlData, err := xml.ConvertPackToXML(b.Bot, *pack)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Сохраняем в файл
	//if err := os.WriteFile("content.xml", xmlData, 0644); err != nil {
	//	log.Fatal(err)
	//}

	log.Println("Bot started")
	b.Bot.Start()

	//ID, err2 := googleDrive.ExtractFileIDFromURL("https://drive.google.com/file/d/1zIY0PyhpXTZSD-T2LzP304KkkzZoNa-p/view")
	//if err2 != nil {
	//	fmt.Println(err2.Error())
	//}
	//file, err3 := googleDrive.DownloadFileByID(ID)
	//if err3 != nil {
	//	fmt.Println(err3.Error())
	//}
	//fmt.Println(string(file))
}
