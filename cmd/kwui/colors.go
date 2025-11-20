package main

import "github.com/veandco/go-sdl2/sdl"

var (
	// Amber retro color scheme
	colorBackground = sdl.Color{R: 26, G: 15, B: 0, A: 255}    // Dark brown-black
	colorAmber      = sdl.Color{R: 255, G: 191, B: 0, A: 255}  // Primary amber
	colorAmberDim   = sdl.Color{R: 180, G: 135, B: 0, A: 255}  // Dimmed amber
	colorAmberGlow  = sdl.Color{R: 255, G: 215, B: 64, A: 255} // Bright amber highlight
	colorBorder     = sdl.Color{R: 100, G: 75, B: 0, A: 255}   // Dark amber for borders
)
