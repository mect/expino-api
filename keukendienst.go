package main

import "time"

type day string
type task string
type name string
type tasks map[task]name

// KeukendienstItem contains for a a period of keukendienst slide
type KeukendienstItem struct {
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Names []string  `json:"names"`
}
