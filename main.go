package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"github.com/mect/go-escpos"
)

var printMutex = sync.Mutex{}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Mvm Voeding Printer")
	})

	e.POST("/print", handleVoedingPrint)
	e.POST("/eenmaligen", handleEenmaligenPrint)
	e.POST("/sinterklaas", handleSinterklaasPrint)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleVoedingPrint(c echo.Context) error {
	data := VoedingRequest{}
	c.Bind(&data)

	printMutex.Lock()

	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	p.Barcode(strings.Replace(data.Doelgroepnummer, "MVM", "", -1), escpos.BarcodeTypeCODE39)
	p.PrintLn("")
	p.PrintLn("")

	p.Align(escpos.AlignLeft)

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
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

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

func handleSinterklaasPrint(c echo.Context) error {
	data := SinterklaasRequest{}
	c.Bind(&data)

	printMutex.Lock()
	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	defer p.Close()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	// p.Size(3, 3)
	// p.PrintLn(data.Speelgoed.MVMNummer)

	// p.Size(2, 2)
	// p.PrintLn(data.Speelgoed.Naam)
	// p.PrintLn("")

	// p.PrintLn("Sinterklaas")

	// for _, entry := range data.Speelgoed.Paketten {
	// 	p.PrintLn("")
	// 	p.PrintLn("----------------")
	// 	p.PrintLn("")

	// 	p.PrintLn(entry.Naam)
	// 	p.PrintLn(entry.Geslacht)
	// 	p.PrintLn(fmt.Sprintf("%.1f jaar", entry.Leeftijd))
	// 	p.PrintLn(entry.Opmerking)

	// 	p.PrintLn("")
	// 	p.PrintLn("----------------")
	// 	p.PrintLn("")
	// }

	// p.Cut()

	p.Size(3, 3)
	p.PrintLn(data.Snoep.MVMNummer)

	p.Size(2, 2)
	p.PrintLn(data.Snoep.Naam)
	p.PrintLn("")

	p.PrintLn("Sinterklaas Snoep")
	p.PrintLn("")
	p.PrintLn(fmt.Sprintf("%d personen", data.Snoep.Personen))

	p.Cut()

	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
