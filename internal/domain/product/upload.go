/*
|------------------------------------------------
| File: internal/domain/product/upload.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package product

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
)

// Endpoint para upload: POST /upload
func UploadImageHandler(c *fiber.Ctx) error {
	// Leitura das variáveis de ambiente (.env)
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-2"
	}
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		bucket = "sabiosystem-produtos"
	}

	fmt.Println("AWS_ACCESS_KEY_ID:", os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Println("AWS_SECRET_ACCESS_KEY:", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	fmt.Println("AWS_S3_BUCKET:", bucket)
	fmt.Println("AWS_REGION:", region)

	// Recebe o arquivo do formulário
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Arquivo não encontrado")
	}
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Falha ao abrir arquivo")
	}
	defer src.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Printf("Erro AWS Session: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao inicializar sessão AWS")
	}
	uploader := s3manager.NewUploader(sess)

	key := fmt.Sprintf("products/%s", file.Filename)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")), // Corrige o Content-Type
	})

	if err != nil {
		fmt.Printf("Erro upload S3: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Falha ao enviar para S3: %v", err))
	}

	fmt.Printf("Upload realizado com sucesso: %s\n", result.Location)
	// Retorna a URL pública
	return c.JSON(fiber.Map{"url": result.Location})
}
