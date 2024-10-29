package logger

import (
	"log"
	"os"
)

// Log является глобальной переменной для логгера
var Log *log.Logger

// InitLogger инициализирует логгер для записи в файл
func InitLogger() {
	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Не удалось открыть файл для логов:", err)
	}

	Log = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
