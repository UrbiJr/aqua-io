package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func MakeCirclePNG(source, dest string) error {
	// Apri il file immagine.
	imgFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	// Decodifica il file immagine.
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	// Crea una nuova immagine vuota di dimensioni quadrate e uguali all'immagine originale.
	bounds := img.Bounds()
	size := bounds.Max.X - bounds.Min.X
	circleImg := image.NewRGBA(image.Rect(0, 0, size, size))

	// Calcola il raggio del cerchio.
	radius := float64(size) / 2

	// Scorri tutti i pixel dell'immagine originale e verifica se si trovano all'interno del cerchio.
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			// Calcola la distanza del pixel dal centro del cerchio.
			dist := math.Sqrt(math.Pow(float64(x)-radius, 2) + math.Pow(float64(y)-radius, 2))

			// Se il pixel si trova all'interno del cerchio, copialo nella nuova immagine.
			if dist <= radius {
				circleImg.Set(x, y, img.At(x, y))
			}
		}
	}

	// Salva la nuova immagine in un file.
	circleImgFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer circleImgFile.Close()

	err = png.Encode(circleImgFile, circleImg)
	if err != nil {
		return err
	}

	return nil
}

func MakeCircleJPG(source, dest string) error {
	// Apri il file immagine.
	inFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Decodifica l'immagine JPG
	img, err := jpeg.Decode(inFile)
	if err != nil {
		return err
	}

	// Crea un'immagine trasparente circolare dello stesso size dell'immagine originale
	mask := image.NewAlpha(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.DrawMask(mask, mask.Bounds(), image.White, image.Point{}, &circle{img.Bounds().Dx() / 2}, image.Point{}, draw.Src)

	// Combina l'immagine originale e la maschera circolare
	out := image.NewRGBA(img.Bounds())
	draw.DrawMask(out, img.Bounds(), img, image.Point{}, mask, image.Point{}, draw.Over)

	// Salva l'immagine risultante su un file JPG
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err = jpeg.Encode(outFile, out, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return nil
}

type circle struct {
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	r := c.r
	return image.Rect(-r, -r, r, r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x), float64(y), float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{A: 255}
	}
	return color.Alpha{}
}
