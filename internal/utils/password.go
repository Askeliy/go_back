package utils

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(p string) string {
	// Генерируем хеш пароля с использованием bcrypt и дефолтной сложности (bcrypt.DefaultCost).
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	// Возвращаем хеш пароля как строку.
	return string(hash)
}

// ComparePassword сравнивает хешированный пароль с введенным паролем.
func ComparePassword(hashedPassword, password string) bool {
	// Сравниваем хешированный пароль с введенным паролем, используя bcrypt.
	// Если хеши совпадают, функция возвращает nil (ошибки нет).
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	// Если ошибок нет (т.е. пароли совпадают), возвращаем true, иначе — false.
	return err == nil
}
