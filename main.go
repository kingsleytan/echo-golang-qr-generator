package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/url"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/labstack/echo"
)

// pass in Querystring by http://localhost:8080/?data=1JAQcb8Q2RQXniwX4We1wT6rfAAkBER5tP&size=300

func main() {
	e := echo.New()
	e.GET("/", QrGenerator)
	e.Logger.Fatal(e.Start(":8080"))
}

func QrGenerator(c echo.Context) error {
	fmt.Println("first")
	data := c.QueryParam("data")
	if data == "" {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot get data",
		})
	}

	s, err := url.QueryUnescape(data)
	if err != nil {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot get query string",
		})
	}

	code, err := qr.Encode(s, qr.L, qr.Auto)
	if err != nil {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot generate QR code",
		})
	}

	size := c.QueryParam("size")
	if size == "" {
		size = "250"
	}
	intsize, err := strconv.Atoi(size)
	if err != nil {
		intsize = 250
	}

	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, intsize, intsize)
	if err != nil {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot scale the QR code",
		})
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, code); err != nil {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot encode png",
		})
	}

	fmt.Println("last")
	fmt.Println("code:", code)

	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := c.Response().Write(buffer.Bytes()); err != nil {
		return c.JSON(422, map[string]interface{}{
			"message": "Cannot generate QR code png",
		})
	}

	return nil
}
