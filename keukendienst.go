package main

type day string
type task string
type name string
type tasks map[task]name

// Keukendienst contains all info for a keukendienst slide
type Keukendienst struct {
	Tasks   []string      `json:"tasks"`
	Days    []string      `json:"days"`
	Content map[day]tasks `json:"content"`
}

// NewKeukendienst gives a Keukendienst with the defaults set
func NewKeukendienst() Keukendienst {
	return Keukendienst{
		Days: []string{"Maandag", "Dinsdag", "Woendsdag", "Donderdag", "Vrijdag"},
	}
}
