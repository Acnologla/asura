package rinha

import "asura/src/entities"

type Room struct {
	Boss  entities.Rooster `json:"boss"`
	Level int              `json:"level"`
}
