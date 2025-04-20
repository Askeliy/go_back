package utils

import (
	"log"
	"net/smtp"
)

// EmailConfig содержит настройки для отправки электронной почты
type EmailConfig struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

// getEmailConfig возвращает настройки электронной почты
func getEmailConfig() EmailConfig {
	return EmailConfig{
		From:     "task-flow@mail.ru",
		Password: "YfJfi6yfF9u5Vd0MPzPa",
		SMTPHost: "SMTP.mail.ru", // SMTP сервер
		SMTPPort: "587",          // Порт SMTP
	}
}

// SendPasswordResetEmail отправляет электронное письмо с подтверждением смены пароля
func SendPasswordResetEmail(to string, code string) error {
	config := getEmailConfig()

	// Создание сообщения
	subject := "Код для изменения пароля"
	body := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Смена пароля</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f9f9f9;
				color: #333;
				padding: 20px;
			}
			.container {
				background-color: #fff;
				border: 1px solid #eaeaea;
				border-radius: 5px;
				padding: 20px;
				max-width: 400px;
				margin: auto;
				box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
			}
			h2 {
				color: #ff6600;
			}
			.code {
				font-size: 24px;
				font-weight: bold;
				color: #ff6600;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Смена пароля</h2>
			<p>Ваш код для изменения пароля:</p>
			<p class="code">` + code + `</p>
			<p>Пожалуйста, введите этот код в приложении для завершения смены пароля.</p>
		</div>
	</body>
	</html>
	`

	// Форматирование сообщения
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	// Аутентификация
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)
	log.Println("Mail sender authorized")
	// Отправка письма
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.From, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}

// SendEmailVerification отправляет электронное письмо с подтверждением регистрации
func SendEmailVerification(to string, code string) error {
	config := getEmailConfig()

	// Создание сообщения
	subject := "Код подтверждения регистрации"
	body := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Подтверждение регистрации</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f9f9f9;
				color: #333;
				padding: 20px;
			}
			.container {
				background-color: #fff;
				border: 1px solid #eaeaea;
				border-radius: 5px;
				padding: 20px;
				max-width: 400px;
				margin: auto;
				box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
			}
			h2 {
				color: #ff6600;
			}
			.code {
				font-size: 24px;
				font-weight: bold;
				color: #ff6600;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Подтверждение регистрации </h2>
			<p>Ваш код для подтверждения регистрации:</p>
			<p class="code">` + code + `</p>
			<p>Пожалуйста, введите этот код в приложении для завершения регистрации.</p>
		</div>
	</body>
	</html>
	`

	// Форматирование сообщения
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	// Аутентификация
	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)
	log.Println("Mail sender authorized")
	// Отправка письма
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.From, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}
