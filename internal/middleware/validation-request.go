package middleware

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidationRequest[V any]() fiber.Handler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(c *fiber.Ctx) error {
		var v V

		if err := c.BodyParser(&v); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := validate.Struct(v); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				// Ambil nama field dari tag jsonc
				fieldName := getJSONTagFieldName(v, err.Field())
				if fieldName == "" {
					fieldName = err.Field() // Default ke field Go jika tidak ada tag json
				}
				message := fmt.Sprintf("%s is %s", fieldName, err.Tag())
				if err.Param() != "" {
					message += " " + err.Param()
				}
				errors = append(errors, message)
			}
			return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
				"code":    fiber.StatusBadRequest,
				"errors":  errors,
				"message": "validation error",
			})
		}
		c.Locals("parser", &v)
		return c.Next()
	}
}

func getJSONTagFieldName(obj any, field string) string {
	r := reflect.TypeOf(obj)

	structField, ok := r.FieldByName(field)
	if !ok {
		return ""
	}
	return structField.Tag.Get("json")
}
