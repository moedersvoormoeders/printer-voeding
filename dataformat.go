package main

type VoedingRequest struct {
	TicketCount     int    `json:"ticketCount"`
	PrintType       string `json:"printType"`
	Doelgroepnummer string `json:"doelgroepnummer"`
	Naam            string `json:"naam"`
	Voornaam        string `json:"voornaam"`
	TypeVoeding     string `json:"typeVoeding"`
	Code            string `json:"code"`
	Volwassenen     int    `json:"volwassenen"`
	Kinderen        int    `json:"kinderen"`
	SpecialeVoeding string `json:"specialeVoeding"`
	NeedsVerjaardag bool   `json:"needsVerjaardag"`
	NeedsMelkpoeder bool   `json:"needsMelkpoeder"`
}

type VoedingRequestEenmaligen struct {
	EenmaligenNummer string `json:"eenmaligenNummer"`
	Naam             string `json:"naam"`
	TypeVoeding      string `json:"typeVoeding"`
	Grootte          string `json:"grootte"`
}
