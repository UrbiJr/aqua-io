package resources

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type NyxDarkTheme struct{}
type NyxLightTheme struct{}

const (
	IconNameCreditCard      fyne.ThemeIconName = "creditCard"
	IconNameWifi            fyne.ThemeIconName = "wifi"
	IconNameBell            fyne.ThemeIconName = "bell"
	IconNameBellExclamation fyne.ThemeIconName = "bellExclamation"
)

// It is a good idea to assert that we implement an interface
// so that compile errors are closer to the defining type.
var _ fyne.Theme = (*NyxDarkTheme)(nil)
var _ fyne.Theme = (*NyxLightTheme)(nil)

func (m NyxDarkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		// rgb(13, 17, 23)
		return &color.RGBA{R: 13, G: 17, B: 23, A: 255}
	case theme.ColorNameForeground:
		// rgb(201, 209, 217);
		return &color.RGBA{R: 201, G: 209, B: 217, A: 255}
	case theme.ColorNamePrimary:
		// rgb(88, 166, 255);
		return &color.RGBA{R: 88, G: 166, B: 255, A: 255}
	case theme.ColorNameHover:
		return &color.RGBA{A: 42}
	case theme.ColorNameFocus:
		return &color.RGBA{R: 13, G: 17, B: 23, A: 255}
	case theme.ColorNameInputBackground:
		// removes background from entry
		return color.Transparent
		//return &color.RGBA{R: 198, G: 210, B: 16, A: 75}
	}

	return theme.DefaultTheme().Color(n, v)
}

func (m NyxDarkTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	switch n {
	case IconNameWifi:
		return resourceWifiWhiteSvg
	case IconNameCreditCard:
		return resourceCreditCardWhiteSvg
	case IconNameBell:
		return resourceBellWhiteSvg
	case IconNameBellExclamation:
		return resourceBellExclamationWhiteSvg
	}
	return theme.DefaultTheme().Icon(n)
}

func (m NyxDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m NyxDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m NyxLightTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		// rgb(13, 17, 23)
		return &color.RGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameForeground:
		// rgb(36, 41, 47);
		return &color.RGBA{R: 36, G: 41, B: 47, A: 255}
	case theme.ColorNameButton:
		//rgb(246 248 250)
		return &color.RGBA{R: 200, G: 200, B: 200, A: 255}
	case theme.ColorNameHover:
		return &color.RGBA{A: 42}
	case theme.ColorNamePrimary:
		// rgb(9, 105, 218);
		return &color.RGBA{R: 9, G: 105, B: 218, A: 255}
	case theme.ColorNameFocus:
		return &color.RGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameInputBackground:
		// removes background from entry
		return color.Transparent
		//return &color.RGBA{R: 198, G: 210, B: 16, A: 75}
	}

	return theme.DefaultTheme().Color(n, v)
}

func (m NyxLightTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	switch n {
	case IconNameWifi:
		return resourceWifiSvg
	case IconNameCreditCard:
		return resourceCreditCardSvg
	case IconNameBell:
		return resourceBellSvg
	case IconNameBellExclamation:
		return resourceBellExclamationSvg
	}
	return theme.DefaultTheme().Icon(n)
}

func (m NyxLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m NyxLightTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
