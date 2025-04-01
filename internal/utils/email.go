package utils

import (
	"log"
	"net/smtp"
)

// SendEmail отправляет электронное письмо с подтверждением смены пароля
func SendEmail(to string, code string) error {
	from := "viteek123@mail.ru"        // Замените на ваш email
	password := "8Czxetq6gnvVm3zfWxt0" // Замените на ваш пароль

	// Настройки SMTP
	smtpHost := "SMTP.mail.ru" // Замените на ваш SMTP сервер
	smtpPort := "587"          // Порт SMTP

	// Создание сообщения
	subject := "Код для изменения пароля"
	body := "Ваш код для изменения пароля: " + code

	message := []byte("Subject: " + subject + "\r\n" + body)

	// Аутентификация
	auth := smtp.PlainAuth("", from, password, smtpHost)
	log.Println("Mail sender authorized")
	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}
