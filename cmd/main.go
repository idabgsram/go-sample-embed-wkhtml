package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load Env File
	e := godotenv.Load()
	if e != nil {
		log.Println(".env not found, using global variable")
	}

	// Set Fiber Port (default to 4000)
	fiberPort := os.Getenv("APP_PORT")
	if fiberPort == "" {
		fiberPort = "4000"
	}

	// Fiber Initialization
	app := fiber.New(fiber.Config{})

	app.Post("/build-pdf", HandleBuildPdf)

	// Start Fiber HTTP Service
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", fiberPort)); err != nil {
			log.Fatal(err)
		}
	}()

	// Create a channel to listen for OS signals
	interrupt := make(chan os.Signal, 1)
	// Register this channel to listen for SIGTERM and SIGINT signals
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	// Wait for one of these signals to be received
	<-interrupt
}

func HandleBuildPdf(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles("assets/sample_content.html")
	if err != nil {
		fmt.Printf("Failed to parse content file, err %s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.StatusBadRequest,
			"meta":   map[string]any{},
			"data":   nil,
			"messages": map[string]any{
				"default": "PDF Source Parse Failed, requested content not valid",
			},
		})
	}

	var parser bytes.Buffer

	// If need to replace any content, do here in the map[string]any
	err = tmpl.Execute(&parser, map[string]any{})
	if err != nil {
		fmt.Printf("Failed to embed content to html file, err %s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.StatusBadRequest,
			"meta":   map[string]any{},
			"data":   nil,
			"messages": map[string]any{
				"default": "PDF Build Failed, requested content not valid",
			},
		})
	}

	err = BuildPDF(parser, "output/test.pdf")
	if err != nil {
		fmt.Printf("Failed to build pdf file, err %s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
			"meta":   map[string]any{},
			"data":   nil,
			"messages": map[string]any{
				"default": "PDF Build Failed",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"meta":   map[string]any{},
		"data":   nil,
		"messages": map[string]any{
			"default": "PDF Build Success",
		},
	})
}

func BuildPDF(parse bytes.Buffer, outputName string) error {
	pdf := wkhtmltopdf.NewPDFPreparer()
	res := wkhtmltopdf.NewPageReader(&parse)
	res.DisableExternalLinks.Set(false)
	res.EnableLocalFileAccess.Set(true)

	// Set Header And Footer if Any
	res.HeaderHTML.Set("assets/sample_header.html")
	res.FooterHTML.Set("assets/sample_footer.html")

	pdf.AddPage(res)

	// Set PDF Margin
	pdf.MarginLeft.Set(0)
	pdf.MarginBottom.Set(51)
	pdf.MarginRight.Set(0)

	js, err := pdf.ToJSON()

	if err != nil {
		return err
	}

	pdfFromJson, err1 := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewReader(js))
	if err1 != nil {
		return err
	}

	err = pdfFromJson.Create()
	if err != nil {
		return err
	}

	return pdfFromJson.WriteFile(outputName)
}
