package main

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/adrg/sysfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// resolveFontPath resolves a font name to a file path using fontconfig
// It tries the requested font first, then falls back to the provided alternatives
func resolveFontPath(fontName string, fallbacks []string) (string, error) {
	finder := sysfont.NewFinder(nil)

	// Try the requested font first
	font := finder.Match(fontName)
	if font != nil && font.Filename != "" {
		slog.Info("resolved font", "requested", fontName, "path", font.Filename)
		return font.Filename, nil
	}

	// Try fallback fonts
	for _, fallback := range fallbacks {
		font = finder.Match(fallback)
		if font != nil && font.Filename != "" {
			slog.Warn("using fallback font", "requested", fontName, "fallback", fallback, "path", font.Filename)
			return font.Filename, nil
		}
	}

	return "", fmt.Errorf("could not resolve font: %s (tried %d fallbacks)", fontName, len(fallbacks))
}

// formatFrequency formats a frequency string to 3 decimal places
func (a *App) formatFrequency(freqStr string) string {
	freq, err := strconv.ParseFloat(freqStr, 64)
	if err != nil {
		// If parsing fails, return the original string
		return freqStr
	}
	return fmt.Sprintf("%.3f", freq)
}

// App represents the main application
type App struct {
	renderer   *sdl.Renderer
	model      *AppModel
	fontLarge  *ttf.Font
	fontMedium *ttf.Font
	fontSmall  *ttf.Font
	running    bool
}

// NewApp creates a new application instance
func NewApp(renderer *sdl.Renderer, kwctlCmd string) *App {
	return &App{
		renderer: renderer,
		model:    NewAppModel(kwctlCmd),
		running:  true,
	}
}

// Init initializes the application
func (a *App) Init() error {
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("failed to initialize TTF: %w", err)
	}

	// Use fontconfig to resolve the font name to a file path
	fontPath, err := resolveFontPath(
		"DejaVu Sans Mono",
		[]string{
			"Liberation Mono", // Common on RHEL/Fedora/CentOS
			"FreeMono",        // GNU FreeFont alternative
			"Courier New",     // Windows fallback
			"Courier",         // macOS fallback
			"monospace",       // Generic monospace alias
		},
	)
	if err != nil {
		return fmt.Errorf("failed to find monospace font: %w", err)
	}

	slog.Info("using font", "path", fontPath)

	a.fontLarge, err = ttf.OpenFont(fontPath, 72)
	if err != nil {
		return fmt.Errorf("failed to load large font: %w", err)
	}

	a.fontMedium, err = ttf.OpenFont(fontPath, 24)
	if err != nil {
		return fmt.Errorf("failed to load medium font: %w", err)
	}

	a.fontSmall, err = ttf.OpenFont(fontPath, 16)
	if err != nil {
		return fmt.Errorf("failed to load small font: %w", err)
	}

	// Start background status polling (non-blocking)
	a.model.StartPolling()

	return nil
}

// Cleanup releases application resources
func (a *App) Cleanup() {
	if a.fontLarge != nil {
		a.fontLarge.Close()
	}
	if a.fontMedium != nil {
		a.fontMedium.Close()
	}
	if a.fontSmall != nil {
		a.fontSmall.Close()
	}
	if a.model != nil {
		a.model.Cleanup()
	}
	ttf.Quit()
}

// Run starts the main application loop
func (a *App) Run() error {
	for a.running {
		a.handleEvents()

		// Check for status updates (non-blocking)
		select {
		case update := <-a.model.statusChan:
			a.model.HandleStatusUpdate(update)
		default:
		}

		a.render()
		sdl.Delay(16) // ~60 FPS
	}
	return nil
}

// handleEvents processes SDL events
func (a *App) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			a.running = false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				a.handleKeyPress(e.Keysym)
			}
		}
	}
}

// handleKeyPress handles keyboard input
func (a *App) handleKeyPress(keysym sdl.Keysym) {
	switch keysym.Sym {
	case sdl.K_q, sdl.K_ESCAPE:
		a.running = false
	case sdl.K_a:
		// Switch to VFO A (not implemented yet - would call kwctl command)
		slog.Info("VFO A selected")
	case sdl.K_b:
		// Switch to VFO B (not implemented yet - would call kwctl command)
		slog.Info("VFO B selected")
	}
}

// render draws the entire UI
func (a *App) render() {
	// Clear screen with background color
	a.renderer.SetDrawColor(colorBackground.R, colorBackground.G, colorBackground.B, colorBackground.A)
	a.renderer.Clear()

	// Draw VFO panels (title removed to free up vertical space)
	a.drawVfoPanel(0, 10, 40, 380, 360)
	a.drawVfoPanel(1, 410, 40, 380, 360)

	// Draw status bar
	a.drawStatusBar()

	// Present the rendered frame
	a.renderer.Present()
}

// drawVfoPanel renders a single VFO display panel
func (a *App) drawVfoPanel(vfoIdx int, x, y, width, height int32) {
	vfo := a.model.status.Vfos[vfoIdx]

	// Draw panel border
	borderColor := colorBorder
	if vfoIdx == a.model.status.CtlVfo {
		borderColor = colorAmber
	}
	a.drawRect(x, y, width, height, borderColor, false)

	// Draw VFO label
	label := fmt.Sprintf("VFO %c", 'A'+vfoIdx)
	a.drawText(label, a.fontMedium, colorAmberGlow, x+10, y+5)

	// Draw status indicators (PTT and CTL) as buttons
	indicatorX := x + 90
	if vfoIdx == a.model.status.PttVfo {
		a.drawIndicatorButton("PTT", indicatorX, y+5, true)
		indicatorX += 65
	}
	if vfoIdx == a.model.status.CtlVfo {
		a.drawIndicatorButton("CTL", indicatorX, y+5, true)
	}

	// Draw frequency (large display)
	freqY := y + 50
	freqText := a.formatFrequency(vfo.Vfo.RxFreq)
	a.drawText(freqText, a.fontLarge, colorAmber, x+20, freqY)
	// Draw MHz label with smaller font below the frequency
	a.drawText("MHz", a.fontSmall, colorAmberDim, x+25, freqY+80)

	// Draw mode info
	infoY := freqY + 110
	//a.drawText(fmt.Sprintf("Mode: %s", vfo.Mode), a.fontSmall, colorAmberDim, x+20, infoY)

	// Draw channel with optional name
	var channelText string
	if vfo.Mode == "memory" {
		channelText = fmt.Sprintf("Channel: %d", vfo.ChannelNumber)
		if vfo.ChannelName != "" {
			channelText += " " + vfo.ChannelName
		}
	} else {
		channelText = "Channel: ---"
	}
	a.drawText(channelText, a.fontSmall, colorAmberDim, x+20, infoY+25)

	a.drawText(fmt.Sprintf("TX Power: %s", vfo.TxPower), a.fontSmall, colorAmberDim, x+20, infoY+50)
	a.drawText(fmt.Sprintf("Shift: %s", vfo.Vfo.Shift), a.fontSmall, colorAmberDim, x+20, infoY+75)

	if vfo.Vfo.Tone == "true" || vfo.Vfo.CTCSS == "true" {
		toneText := fmt.Sprintf("Tone: %s Hz", vfo.Vfo.ToneFreq)
		if vfo.Vfo.CTCSS == "true" {
			toneText = fmt.Sprintf("%s/%s", toneText, vfo.Vfo.CTCSSFreq)
		}
		a.drawText(toneText, a.fontSmall, colorAmberDim, x+20, infoY+100)
	}

	// Draw mode buttons (moved down to avoid overlap with TONE label)
	a.drawModeButtons(vfo.Mode, x+20, y+height-70)
}

// drawIndicatorButton renders a single indicator button (PTT, CTL)
func (a *App) drawIndicatorButton(label string, x, y int32, active bool) {
	buttonWidth := int32(55)
	buttonHeight := int32(24)

	// Draw button with appropriate style
	color := colorBorder
	if active {
		color = colorAmber
	}
	a.drawRect(x, y, buttonWidth, buttonHeight, color, active)

	// Draw text with appropriate color
	textColor := colorAmberDim
	if active {
		textColor = colorBackground
	}
	a.drawText(label, a.fontSmall, textColor, x+8, y+4)
}

// drawModeButtons renders the VFO mode selection buttons
func (a *App) drawModeButtons(currentMode string, x, y int32) {
	modes := []string{"vfo", "memory", "call", "wx"}
	buttonWidth := int32(70)
	buttonHeight := int32(30)
	spacing := int32(10)

	for i, mode := range modes {
		bx := x + int32(i)*(buttonWidth+spacing)
		color := colorBorder
		if mode == currentMode {
			color = colorAmber
		}
		a.drawRect(bx, y, buttonWidth, buttonHeight, color, mode == currentMode)

		// Center text in button
		textColor := colorAmberDim
		if mode == currentMode {
			textColor = colorBackground
		}
		a.drawText(mode, a.fontSmall, textColor, bx+10, y+7)
	}
}

// drawStatusBar renders the bottom status bar
func (a *App) drawStatusBar() {
	y := int32(windowHeight - 40)

	if a.model.errorMsg != "" {
		a.drawText("ERROR: "+a.model.errorMsg, a.fontSmall, colorAmber, 20, y)
	} else {
		helpText := "[A] VFO A  [B] VFO B  [Q]/[ESC] Exit"
		a.drawText(helpText, a.fontSmall, colorAmberDim, 20, y)

		// Draw last update time
		if !a.model.lastUpdate.IsZero() {
			updateText := fmt.Sprintf("Updated: %s", a.model.lastUpdate.Format("15:04:05"))
			a.drawText(updateText, a.fontSmall, colorAmberDim, 600, y)
		}
	}
}

// drawText renders text at the specified position
func (a *App) drawText(text string, font *ttf.Font, color sdl.Color, x, y int32) {
	if font == nil {
		return
	}

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return
	}
	defer surface.Free()

	texture, err := a.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	a.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H})
}

// drawRect draws a rectangle, optionally filled
func (a *App) drawRect(x, y, width, height int32, color sdl.Color, filled bool) {
	rect := sdl.Rect{X: x, Y: y, W: width, H: height}
	a.renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	if filled {
		a.renderer.FillRect(&rect)
	} else {
		a.renderer.DrawRect(&rect)
	}
}
