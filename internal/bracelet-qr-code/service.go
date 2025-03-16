package braceletqrcode

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/skip2/go-qrcode"

	"bracelet-ticket-system-be/internal/config"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/utils"
)

type braceletQrCodeService struct {
	cfg *config.Config
}

func BraceletQrCodeService(cfg *config.Config) domain.BraceletQrCodeService {
	return &braceletQrCodeService{cfg: cfg}
}

func (b *braceletQrCodeService) GenerateBraceletQrCode(braceletQrCodeData []domain.BraceletQrCodeData, folderName string) (string, error) {

	folderPath := fmt.Sprintf("%s/%s", "qr-code-output", folderName)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed make folder: %v", err)
	}

	for _, data := range braceletQrCodeData {

		dataJSON, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed conversion to JSON: %v", err)
		}

		encryptedData, err := utils.EncryptAESGCM(string(dataJSON), b.cfg.EncriptionKey)
		if err != nil {
			return "", fmt.Errorf("failed encript the data: %v", err)
		}

		fileName := fmt.Sprintf("%s.png", data.NoTicket)
		filePath := filepath.Join(folderPath, fileName)

		err = qrcode.WriteFile(encryptedData, qrcode.Medium, 1080, filePath)
		if err != nil {
			return "", fmt.Errorf("failed make QR Code: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	zipFilePath := fmt.Sprintf("%s.zip", folderPath)
	err = zipFolder(folderPath, zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed zip the folder: %v", err)
	}

	err = os.RemoveAll(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed to remove folder: %v", err)
	}

	return folderName + ".zip", nil
}

func zipFolder(source, destination string) error {
	zipFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		writer, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
