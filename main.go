package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Configuration
	cookieDomain := "" // Empty for localhost, set to your domain for production
	env := "development"
	frontendURL := "http://localhost:8083" // Update this for your frontend

	// CORS configuration for cross-domain requests
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8083,https://*.netlify.app,https://*.ngrok.io,https://*.ngrok-free.app",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Route to set HTTP-only cookie
	app.Post("/set-cookie", func(c *fiber.Ctx) error {
		isSecure := true
		employeeID := "EMP123456"
		accessToken := "secure_session_12345"

		log.Printf("[COOKIE] Setting access_token cookie for employee: %s", employeeID)
		log.Printf("[COOKIE] Cookie config - Domain: %s, Secure: %v, Environment: %s", cookieDomain, isSecure, env)

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HTTPOnly: true,
			Secure:   isSecure,
			SameSite: "None",
			Domain:   cookieDomain,
			Expires:  time.Now().Add(10 * time.Minute),
		})

		redirectURL := fmt.Sprintf("%s/auth/callback", frontendURL)
		log.Printf("[COOKIE] Cookie set successfully, redirecting to: %s", redirectURL)
		log.Printf("[COOKIE] ===== Cookie Setting Completed =====")

		return c.JSON(fiber.Map{
			"message":     "HTTP-only cookie set successfully",
			"success":     true,
			"redirectURL": redirectURL,
			"employee":    employeeID,
		})
	})

	// Route to verify cookie exists
	app.Get("/verify-cookie", func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")

		log.Printf("[COOKIE] Verifying access_token cookie...")

		if accessToken == "" {
			log.Printf("[COOKIE] No access_token cookie found")
			return c.Status(401).JSON(fiber.Map{
				"message": "No access_token cookie found",
				"success": false,
			})
		}

		log.Printf("[COOKIE] Access token verified successfully")
		return c.JSON(fiber.Map{
			"message": "Access token cookie verified",
			"success": true,
			"token":   accessToken,
		})
	})

	// Route to clear cookie
	app.Post("/clear-cookie", func(c *fiber.Ctx) error {
		log.Printf("[COOKIE] Clearing access_token cookie...")

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
			Domain:   cookieDomain,
		})

		log.Printf("[COOKIE] Access token cookie cleared successfully")
		return c.JSON(fiber.Map{
			"message": "Access token cookie cleared successfully",
			"success": true,
		})
	})

	// Route with redirect functionality (like SAML callback)
	app.Get("/auth/saml/callback", func(c *fiber.Ctx) error {
		isSecure := true
		employeeID := "EMP789012"
		accessToken := "saml_token_67890"

		log.Printf("[SAML] Setting access_token cookie for employee: %s", employeeID)
		log.Printf("[SAML] Cookie config - Domain: %s, Secure: %v, Environment: %s", cookieDomain, isSecure, env)

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HTTPOnly: true,
			Secure:   isSecure,
			SameSite: "None",
			Domain:   cookieDomain,
			Expires:  time.Now().Add(10 * time.Minute),
		})

		redirectURL := fmt.Sprintf("%s/auth/callback", frontendURL)
		log.Printf("[SAML] Cookie set successfully, redirecting to: %s", redirectURL)
		log.Printf("[SAML] ===== SAML Callback Completed =====")

		c.Status(fiber.StatusSeeOther)
		c.Location(redirectURL)
		return c.SendString("Redirecting...")
	})

	log.Println("Server starting on port 3001...")
	log.Fatal(app.Listen(":3001"))
}
