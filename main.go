package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/mect/go-escpos"
)

var printMutex = sync.Mutex{}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Mvm Voeding Printer")
	})

	e.POST("/print", handleVoedingPrint)
	e.POST("/eenmaligen", handleEenmaligenPrint)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleVoedingPrint(c echo.Context) error {
	data := VoedingRequest{}
	c.Bind(&data)

	printMutex.Lock()
	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinter(0, 0) // auto discover USB
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}
	p.Init()

	p.Size(4, 4)
	p.PrintLn(fmt.Sprintf("%d", data.TicketCount))

	if data.PrintType != "" && data.PrintType != "Gewoon" {
		p.PrintLn(data.PrintType)
	}

	p.PrintLn(data.Doelgroepnummer)
	p.Size(2, 2)
	p.PrintLn(fmt.Sprintf("%s %s", data.Voornaam, data.Naam))
	if data.TypeVoeding != "gewoon" {
		p.Align(escpos.AlignRight)
	}
	p.PrintLn(data.TypeVoeding)
	p.Align(escpos.AlignLeft)

	p.PrintLn(data.Code)

	p.PrintLn(fmt.Sprintf("Volwassenen: %d", data.Volwassenen))
	p.PrintLn(fmt.Sprintf("Kinderen: %d", data.Kinderen))

	if data.SpecialeVoeding != "" {
		p.PrintLn(data.SpecialeVoeding)
	}

	if data.NeedsMelkpoeder {
		p.PrintLn("MELKPOEDER")
	}
	if data.NeedsVerjaardag {
		p.PrintLn("VERJAARDAG")
	}

	p.Cut()
	p.End()

	log.Println(data.Doelgroepnummer, data.TypeVoeding)

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func handleEenmaligenPrint(c echo.Context) error {
	data := VoedingRequestEenmaligen{}
	c.Bind(&data)

	printMutex.Lock()
	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinter(0, 0) // auto discover USB
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}
	p.Init()

	p.Size(4, 4)
	p.PrintLn("VR")

	p.Size(2, 2)
	p.PrintLn(data.EenmaligenNummer)

	p.PrintLn(data.Naam)
	if data.TypeVoeding != "gewoon" {
		p.Align(escpos.AlignRight)
	}
	p.PrintLn(data.TypeVoeding)
	p.Align(escpos.AlignLeft)

	p.PrintLn(data.Grootte)

	p.Cut()
	p.End()

	log.Println(data.EenmaligenNummer, data.TypeVoeding)

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
