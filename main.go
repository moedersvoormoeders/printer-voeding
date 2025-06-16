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
	e.POST("/print-markt", handleMarktPrint)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleVoedingPrint(c echo.Context) error {
	data := VoedingRequest{}
	c.Bind(&data)

	printMutex.Lock()

	defer printMutex.Unlock()
	p, err := escpos.NewUSBPrinterByPath("") // auto discover USB
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}
	defer p.Close()

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	if data.Doelgroepnummer != "" {
		p.Barcode(strings.Replace(data.Doelgroepnummer, "MVM", "", -1), escpos.BarcodeTypeCODE39)
		p.PrintLn("")
		p.PrintLn("")
	}

	p.Align(escpos.AlignLeft)

	p.Size(4, 4)
	p.PrintLn(fmt.Sprintf("%d", data.TicketCount))

	if data.PrintType != "" && data.PrintType != "Gewoon" {
		p.PrintLn(data.PrintType)
	}

	if data.Doelgroepnummer != "" {
		p.PrintLn(data.Doelgroepnummer)
	}
	if data.EenmaligenNummer != "" {
		p.PrintLn(data.EenmaligenNummer)
	}
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
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": err.Error()})
	}
	defer p.Close()

	err = p.Init()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, echo.Map{"status": "error", "error": "Printer reageert niet, check status en papier"})
	}

	p.Size(4, 4)
	p.PrintLn(fmt.Sprintf("%d", data.Speelgoed.VolgNummer))

	p.Size(4, 4)
	p.PrintLn("")
	p.PrintLn(data.Speelgoed.MVMNummer)

	log.Printf("Speelgoed voor %s\n", data.Speelgoed.MVMNummer)

	p.Size(2, 2)
	p.PrintLn(data.Speelgoed.Naam)
	p.PrintLn("")

	p.PrintLn("Sinterklaas")

	for _, entry := range data.Speelgoed.Paketten {
		p.PrintLn("")
		p.PrintLn("----------------")
		p.PrintLn("")

		p.PrintLn(entry.Naam)
		p.PrintLn(entry.Geslacht)
		if entry.Leeftijd < 1 {
			p.PrintLn(fmt.Sprintf("%.1f jaar", entry.Leeftijd))
		} else {
			p.PrintLn(fmt.Sprintf("%.0f jaar", entry.Leeftijd))
		}

		p.PrintLn(entry.Opmerking)

		p.PrintLn("")
		p.PrintLn("----------------")
		p.PrintLn("")

		log.Printf("%s is braaf geweest\n", entry.Naam)
	}

	p.Cut()
	p.End()

	// p.Size(3, 3)
	// p.PrintLn(data.Snoep.MVMNummer)

	// p.Size(2, 2)
	// p.PrintLn(data.Snoep.Naam)
	// p.PrintLn("")

	// p.PrintLn("Sinterklaas Snoep")

	// p.PrintLn(fmt.Sprintf("volwassenen: %d", data.Snoep.Volwassenen))
	// p.PrintLn(fmt.Sprintf("kinderen: %d", data.Snoep.Kinderen))

	// p.Cut()

	// p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func handleMarktPrint(c echo.Context) error {
	data := MarktRequest{}
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

	p.Align(escpos.AlignLeft)
	p.Size(2, 2)
	if data.MVMNummer != "" {
		p.PrintLn("MVM" + data.MVMNummer)
	}
	p.Size(2, 2)
	p.PrintLn("")
	p.PrintLn(data.Naam)
	p.Size(2, 1)
	p.PrintLn(data.Beschrijving)
	p.Feed(2)

	for _, entry := range data.Kinderen {
		p.Size(1, 1)
		p.PrintLn("==========================================")
		p.Size(2, 2)

		if entry.Naam != "" {
			p.Size(1, 1)
			p.PrintLn(entry.Naam)
			p.Size(2, 1)
			p.Print(entry.Geslacht)
			p.Print(" ")
			p.PrintLn(fmt.Sprintf("%d jaar", entry.Leeftijd))
		}
	}

	p.Size(1, 1)
	p.PrintLn("==========================================")
	p.AztecViaImage(data.Barcode, 400, 400)

	p.Feed(2)

	p.Cut()
	p.End()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
