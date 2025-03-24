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
	r.Post("/download-exel", middleware.ValidationRequest[domain.GetBraceletTicketExelReq](), handler.GetBraceletTicketExcelFile)
	r.Post("/generate", middleware.ValidationRequest[domain.GenerateBraceletTicketReq](), handler.GenerateBraceletTicket)
	r.Post("/check-in-online", middleware.ValidationRequest[domain.CheckInBraceletTicketOnlineRequest](), handler.CheckInBraceletTicketOnline)
	r.Post("/check-in-offline", middleware.ValidationRequest[domain.CheckInBraceletTicketOfflineRequest](), handler.CheckInBraceletTicketOffline)
	r.Get("/get-list-filename-exel/:eventID", handler.GetListFileNameBraceletTicketExelByEventID)
	r.Post("/check-in-online-manual", middleware.ValidationRequest[domain.CheckInBraceletTicketWithSerialNumberOnlineRequest](), handler.CheckInBraceletTicketOnlineWithSerialNumber)
	r.Post("/check-in-offline-manual", middleware.ValidationRequest[domain.CheckInBraceletTicketWithSerialNumberOfflineRequest](), handler.CheckInBraceletTicketOfflineWithSerialNumber)
}

func (h *httpBraceletTicketHandler) CheckInBraceletTicketOnline(c *fiber.Ctx) error {
	logger := xlogger.Logger
	var requestBody domain.CheckInBraceletTicketOnlineRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}
	logger.Info().Msgf("CheckInBraceletTicketRequest: %v", requestBody)

	response, err := h.braceletTicketService.CheckInBraceletTicketOnline(requestBody.EventID, requestBody.QrData, requestBody.DeviceID, requestBody.DeviceName)
	if err != nil {
		return err
	}

	return c.Status(response.StatusCode).JSON(response)
}

func (h *httpBraceletTicketHandler) CheckInBraceletTicketOffline(c *fiber.Ctx) error {
	logger := xlogger.Logger
	var requestBody domain.CheckInBraceletTicketOfflineRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}
	logger.Info().Msgf("CheckInBraceletTicketRequest: %v", requestBody)

	go func() {
		err := h.braceletTicketService.CheckInBraceletTicketOffline(requestBody.Data)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to check in bracelet ticket offline")
		}
	}()

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponseWithaoutData{
		StatusCode: fiber.StatusOK,
		Error:      false,
		Message:    "Success syncronize bracelet ticket, wait a moment",
	})
}

func (h *httpBraceletTicketHandler) GetBraceletTicketExcelFile(c *fiber.Ctx) error {
	var requestBody domain.GetBraceletTicketExelReq

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusBadRequest,
			Error:      true,
			Message:    "Failed to parse request body",
		})
	}

	if requestBody.FileName == "" || requestBody.EventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusBadRequest,
			Error:      true,
			Message:    "File name and event ID are required",
		})
	}

	if filepath.Ext(requestBody.FileName) != ".xlsx" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusBadRequest,
			Error:      true,
			Message:    "Invalid file extension, only .xlsx files are allowed",
		})
	}

	filePath := fmt.Sprintf("folder-bracelet-ticket-excel/%s/%s", requestBody.EventID, requestBody.FileName)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponseWithaoutData{
				StatusCode: fiber.StatusNotFound,
				Error:      true,
				Message:    "File not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusInternalServerError,
			Error:      true,
			Message:    "Error checking file",
		})
	}

	// Set Content-Type agar dikenali sebagai Excel file
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", requestBody.FileName))

	// Kirim file ke client
	return c.SendFile(filePath)
}

func (h *httpBraceletTicketHandler) GenerateBraceletTicket(c *fiber.Ctx) error {
	var requestBody domain.GenerateBraceletTicketReq
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusBadRequest,
			Error:      true,
			Message:    "failed to parse request body",
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

func (h *httpBraceletTicketHandler) GetListFileNameBraceletTicketExelByEventID(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	if eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&domain.ApiResponseWithaoutData{
			StatusCode: fiber.StatusBadRequest,
			Error:      true,
			Message:    "Invalid event ID",
		})
	}

	response, err := h.braceletTicketService.GetListFileNameExelBaceletTicketByEventID(eventID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&domain.ApiResponse{
		Error:   false,
		Message: "Success",
		Data:    response,
	})
}

func (h *httpBraceletTicketHandler) CheckInBraceletTicketOnlineWithSerialNumber(c *fiber.Ctx) error {
	logger := xlogger.Logger
	var requestBody domain.CheckInBraceletTicketWithSerialNumberOnlineRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}
	logger.Info().Msgf("CheckInBraceletTicketRequest: %v", requestBody)

	response, err := h.braceletTicketService.CheckInBraceletTicketOnlineManual(requestBody.EventID, requestBody.SerialNumber, requestBody.DeviceID, requestBody.DeviceName)
	if err != nil {
		return fmt.Errorf("failed to check in bracelet ticket online manual: %v", err)
	}

	return c.Status(response.StatusCode).JSON(response)
}

func (h *httpBraceletTicketHandler) CheckInBraceletTicketOfflineWithSerialNumber(c *fiber.Ctx) error {
	logger := xlogger.Logger
	var requestBody domain.CheckInBraceletTicketWithSerialNumberOfflineRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponseWithaoutData{
			Error:   true,
			Message: "failed to parse request body",
		})
	}
	logger.Info().Msgf("CheckInBraceletTicketRequest: %v", requestBody)

	go func() {
		err := h.braceletTicketService.CheckInBraceletTicketOfflineManual(requestBody.Data)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to check in bracelet ticket offline manual")
		}
	}()

	return c.Status(fiber.StatusOK).JSON(domain.ApiResponseWithaoutData{
		StatusCode: fiber.StatusOK,
		Error:      false,
		Message:    "Success syncronize bracelet ticket manual, wait a moment",
	})
}
