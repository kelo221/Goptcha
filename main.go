package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"image/png"
	"log"
	"net/http"
	"time"
)

var store *session.Store

func main() {

	app := fiber.New()

	store = session.New(session.Config{
		Expiration:     time.Minute * 1,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	app.Get("/captcha", getCaptcha)
	app.Post("/checker", restricted)

	err := app.Listen(":3000")
	if err != nil {
		return
	}

}

func getCaptcha(c *fiber.Ctx) error {

	captchaString, img := generateCaptcha()

	fmt.Println(captchaString)

	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set("captcha", captchaString)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	// Create a new buffer and encode an image to it
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set the content type header to image/jpeg
	c.Set(fiber.HeaderContentType, "image/png")

	// Write the image bytes to the response body
	if _, err := c.Write(buf.Bytes()); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func restricted(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		println("parsing error")
		return err
	}

	if data["captcha"] == "" {
		c.Status(400)
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing captcha!",
		})
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	captcha := sess.Get("captcha")

	if captcha != data["captcha"] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect captcha!"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"Success": "Correct captcha!"})
}
