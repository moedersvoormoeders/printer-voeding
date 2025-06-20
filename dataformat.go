package main

type VoedingRequest struct {
	TicketCount      int    `json:"ticketCount"`
	PrintType        string `json:"printType"`
	Doelgroepnummer  string `json:"doelgroepnummer"`
	EenmaligenNummer string `json:"eenmaligenNummer"`
	Naam             string `json:"naam"`
	Voornaam         string `json:"voornaam"`
	TypeVoeding      string `json:"typeVoeding"`
	Code             string `json:"code"`
	Volwassenen      int    `json:"volwassenen"`
	Kinderen         int    `json:"kinderen"`
	SpecialeVoeding  string `json:"specialeVoeding"`
	NeedsVerjaardag  bool   `json:"needsVerjaardag"`
	NeedsMelkpoeder  bool   `json:"needsMelkpoeder"`
}

type VoedingRequestEenmaligen struct {
	EenmaligenNummer string `json:"eenmaligenNummer"`
	Naam             string `json:"naam"`
	TypeVoeding      string `json:"typeVoeding"`
	Grootte          string `json:"grootte"`
}

type SinterklaasRequest struct {
	// fun fact: we designed these first to be seperate calls
	Speelgoed struct {
		MVMNummer  string `json:"mvmNummer"`
		VolgNummer int    `json:"volgNummer"`
		Naam       string `json:"naam"`
		Paketten   []struct {
			Naam      string  `json:"naam"`
			Geslacht  string  `json:"geslacht"`
			Leeftijd  float64 `json:"leeftijd"` // float64 really!
			Opmerking string  `json:"opmerking"`
		} `json:"paketten"`
	} `json:speelgoed`
	Snoep struct {
		MVMNummer string `json:"mvmNummer"`
		Naam      string `json:"naam"`
		Personen  int    `json:"personen"`
	} `json:snoep`
}

type MarktTicketKind struct {
	Naam     string `json:"naam"`
	Leeftijd int    `json:"leeftijd"`
	Geslacht string `json:"geslacht"`
}

type MarktRequest struct {
	Naam         string `json:"naam"`
	Beschrijving string `json:"beschrijving"`

	MVMNummer string `json:"mvmnummer"`

	Kinderen []MarktTicketKind `json:"kinderen"`

	MarktID uint `json:"markt_id"`

	Barcode string `json:"barcode"`
}
