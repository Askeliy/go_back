package controllers

import (
	"go_back/internal/initializers"
	"go_back/internal/models"
	"go_back/internal/utils"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword проверяет правильность пароля
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Генерация случайного кода
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999) // Генерируем код от 0 до 999999
	return strconv.Itoa(code)
}

// CreateUser  создает нового пользователя
func CreateUser(c *fiber.Ctx) error {
	var payload *models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Create:User  BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser := models.User{
		ID:       uuid.New(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: utils.GeneratePassword(payload.Password),
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "(SQLSTATE 23505)") {
			log.Printf("Create:User  Email already exists: %s", payload.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exists"})
		}
		log.Printf("Create:User  Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	log.Printf("Create:User  User created successfully: %v", newUser)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": newUser}})
}

// LoginUser  авторизует пользователя
func LoginUser(c *fiber.Ctx) error {
	var payload *models.LoginUserSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Login:User  BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", payload.Email)
	if result.Error != nil {
		log.Printf("Login:User  Invalid credentials for email: %s", payload.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	if err := VerifyPassword(user.Password, payload.Password); err != nil {
		log.Printf("Login:User  Invalid password for email: %s", payload.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	// Генерация токена
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		log.Printf("Login:User  Could not generate token for email: %s, error: %v", payload.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not generate token"})
	}

	log.Printf("Login:User  User logged in successfully: %s", payload.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user, "token": token}})
}

// ChangePassword изменяет пароль пользователя с подтверждением по почте
func ChangePassword(c *fiber.Ctx) error {
	var payload models.ChangePasswordSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("ChangePassword: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid input"})
	}

	// Валидация данных
	if payload.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Email is required"})
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", payload.Email)
	if result.Error != nil {
		log.Printf("ChangePassword: User not found for email: %s", payload.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User  not found"})
	}

	// Генерация кода для подтверждения
	code := generateCode()
	log.Println("Generated code:", code)

	// Сохранение кода и времени истечения
	user.PasswordResetCode = code
	user.CodeExpiry = time.Now().Add(10 * time.Minute) // Код действителен 10 минут
	user.CodeUsed = false                              // Код еще не использован
	initializers.DB.Save(&user)

	// Отправка email с кодом
	log.Println("Sending email...")
	if err := utils.SendEmail(user.Email, code); err != nil {
		log.Printf("ChangePassword: Failed to send email to: %s, error: %v", user.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to send email"})
	}

	log.Printf("ChangePassword: Confirmation email sent to: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Confirmation email sent"})
}

// ConfirmChangePassword подтверждает смену пароля
func ConfirmChangePassword(c *fiber.Ctx) error {
	var payload models.ConfirmChangePasswordSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("ConfirmChangePassword: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid input"})
	}

	// Валидация данных
	if payload.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Email is required"})
	}
	if payload.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Code is required"})
	}
	if payload.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "New password is required"})
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", payload.Email)
	if result.Error != nil {
		log.Printf("ConfirmChangePassword: User not found for email: %s", payload.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User  not found"})
	}

	// Проверка кода, его срока действия и статуса использования
	if user.PasswordResetCode != payload.Code || time.Now().After(user.CodeExpiry) || user.CodeUsed {
		log.Printf("ConfirmChangePassword: Invalid or expired code for email: %s", payload.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid or expired code"})
	}

	// Хеширование нового пароля
	user.Password = utils.GeneratePassword(payload.NewPassword)
	user.PasswordResetCode = ""   // Очистка кода после успешной смены пароля
	user.CodeExpiry = time.Time{} // Очистка времени истечения
	user.CodeUsed = true          // Установка статуса использования кода
	initializers.DB.Save(&user)

	log.Printf("ConfirmChangePassword: Password changed successfully for email: %s", payload.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Password changed successfully"})
}

// DeleteUser   удаляет пользователя
func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	result := initializers.DB.Delete(&models.User{}, "id = ?", userId)

	if result.RowsAffected == 0 {
		log.Printf("Delete:User  No user found with ID: %s", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that Id exists"})
	} else if result.Error != nil {
		log.Printf("Delete:User  Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	log.Printf("Delete:User  User deleted successfully with ID: %s", userId)
	return c.SendStatus(fiber.StatusNoContent)
}

// FindUsers возвращает список пользователей с пагинацией
func FindUsers(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil {
		log.Printf("FindUsers: Database error: %v", results.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	log.Printf("FindUsers: Retrieved %d users", len(users))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "users": users})
}

// FindUser ById возвращает пользователя по ID
func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if result.Error != nil {
		log.Printf("FindUser ById: User not found with ID: %s", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User   not found"})
	}

	log.Printf("FindUser ById: User found with ID: %s", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}
