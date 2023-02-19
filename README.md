# Goptcha

### A proof of consept captcha solution done in Golang.

<img src="https://user-images.githubusercontent.com/61495413/219882869-114165e9-f1fb-4486-90e2-871c1e3c2bb4.png" width="256" height="256" />


<hr>



## Example using Fiber

```go
package main

import (
	"github.com/kelo221/goptcha"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"image/png"
	"log"
	"net/http"
	"time"
)

var capchaStore *session.Store

func main() {

	app := fiber.New()

	capchaStore = session.New(session.Config{
		Expiration:     time.Minute * 1,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	Goptcha.Configure(&Goptcha.Config{
		ImageSizeMultiplier: 4,
		CharSet:             "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		CharacterCount:      8,
		Opacity:             100,
	})

	app.Get("/captcha", getCaptcha)
	app.Post("/checker", restricted)

	err := app.Listen(":3000")
	if err != nil {
		return
	}

}

func getCaptcha(c *fiber.Ctx) error {

	captchaString, img := Goptcha.GenerateCaptcha()

	fmt.Println(captchaString)

	sess, err := capchaStore.Get(c)
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

	sess, err := capchaStore.Get(c)
	if err != nil {
		log.Println(err)
	}

	captcha := sess.Get("captcha")

	log.Println(captcha, data["captcha"])

	if captcha != data["captcha"] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect captcha!"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"Success": "Correct captcha!"})
}

```
Here localhost:3000/captcha returns an randomly generated image, which contains text that is passed to the localhost:3000/checker endpoint for verification.

![image](https://user-images.githubusercontent.com/61495413/218850589-9e30b6dd-4f69-4260-83fc-809644e5e6db.png)

*Sample captcha*
