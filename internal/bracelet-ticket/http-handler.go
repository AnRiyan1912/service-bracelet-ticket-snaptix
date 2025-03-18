package braceletticket

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/middleware"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type httpBraceletTicketHandler struct {
	braceletTicketService domain.BraceletTicketService
}

func NewHttpHandler(r fiber.Router, braceletTicketService domain.BraceletTicketService) {
	handler := &httpBraceletTicketHandler{
		braceletTicketService: braceletTicketService,
	}
	r.Get("/total/:eventID", handler.GetTotalBraceletAndTotalCheckInBraceletTicketByEventID)
	r.Get("/download-qr-code/:fileName", handler.GetBraceletTicketQrCodeFile)
	r.Post("/generate", middleware.ValidationRequest[domain.GenerateBraceletTicketReq](), handler.GenerateBraceletTicket)
	r.Post("/check-in", middleware.ValidationRequest[domain.CheckInBraceletTicketRequest](), handler.CheckInBraceletTicket)
}

func (h *httpBraceletTicketHandler) CheckInBraceletTicket(c *fiber.Ctx) error {
	logger := xlogger.Logger
	var requestBody domain.CheckInBraceletTicketRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}
	logger.Info().Msgf("CheckInBraceletTicketRequest: %v", requestBody)

	response, err := h.braceletTicketService.CheckInBraceletTicket(requestBody.EventID, requestBody.QrData)
	if err != nil {
		return err
	}

	return c.Status(response.StatusCode).JSON(response)
}

func (h *httpBraceletTicketHandler) GetBraceletTicketQrCodeFile(c *fiber.Ctx) error {
	fileName := c.Params("fileName")

	if fileName == "" || filepath.Ext(fileName) != ".zip" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file name",
		})
	}

	filePath := fmt.Sprintf("qr-code-output/%s", fileName)
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	return c.SendFile(filePath)
}

func (h *httpBraceletTicketHandler) GenerateBraceletTicket(c *fiber.Ctx) error {
	var requestBody domain.GenerateBraceletTicketReq
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}

	go func() {
		err := h.braceletTicketService.GenerateBraceletQrCode(requestBody.EventID, requestBody.BraceletCategoryID, requestBody.Total, requestBody.Sessions)
		if err != nil {
			log.Printf("failed to generate bracelet QR code: %v", err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponseWithaoutData{
		StatusCode: fiber.StatusOK,
		Error:      false,
		Message:    "Success generate bracelet ticket, wait a moment",
	})
}

func (h *httpBraceletTicketHandler) GetTotalBraceletAndTotalCheckInBraceletTicketByEventID(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	if eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	res, err := h.braceletTicketService.GetTotalBraceletAndTotalCheckInBraceletTicketByEventID(eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponse{
		Error:   false,
		Message: "Success",
		Data:    res,
	})
}
