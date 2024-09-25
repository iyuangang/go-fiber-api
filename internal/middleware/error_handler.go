package middleware

import (
    "github.com/gofiber/fiber/v2"
    "log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    msg := "Internal Server Error"
    
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
        msg = e.Message
    }

    log.Printf("Error: %s, Path: %s, StatusCode: %d", err.Error(), c.Path(), code)

    return c.Status(code).JSON(fiber.Map{
        "status":  code,
        "message": msg,
    })
}
