package resources

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type NyxTheme struct{}

const (
	IconNameCreditCard      fyne.ThemeIconName = "creditCard"
	IconNameWifi            fyne.ThemeIconName = "wifi"
	IconNameBell            fyne.ThemeIconName = "bell"
	IconNameBellExclamation fyne.ThemeIconName = "bellExclamation"
)

// It is a good idea to assert that we implement an interface
// so that compile errors are closer to the defining type.
var _ fyne.Theme = (*NyxTheme)(nil)

func (m NyxTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (m NyxTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	switch n {
	case IconNameWifi:
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			return resourceWifiSvg
		} else {
			return resourceWifiWhiteSvg
		}
	case IconNameCreditCard:
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			return resourceCreditCardSvg
		} else {
			return resourceCreditCardWhiteSvg
		}
	case IconNameBell:
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			return resourceBellSvg
		} else {
			return resourceBellWhiteSvg
		}
	case IconNameBellExclamation:
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			return resourceBellExclamationSvg
		} else {
			return resourceBellExclamationWhiteSvg
		}
	}
	return theme.DefaultTheme().Icon(n)
}

func (m NyxTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m NyxTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
