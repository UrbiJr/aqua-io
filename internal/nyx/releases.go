package nyx

import (
	"net/url"
)

type Release struct {
	Title      string
	Date       string
	StockXLink *url.URL
	ImageURL   *url.URL
}

func (app *Config) fetchReleases() []Release {
	// TODO: fetch the actual releases from db
	var releases []Release
	url, _ := url.Parse("https://images.stockx.com/360/MSCHF-Big-Red-Boot/Images/MSCHF-Big-Red-Boot/Lv2/img01.jpg")
	stockxLink, _ := url.Parse("https://stockx.com/mschf-big-red-boot")
	r := Release{
		Title:      "MSCHF Red Boot",
		Date:       "16 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	url, _ = url.Parse("https://images.stockx.com/images/Air-Jordan-13-Retro-Playoffs-2023-Product.jpg")
	stockxLink, _ = url.Parse("https://stockx.com/air-jordan-13-retro-playoffs-2023")
	r = Release{
		Title:      "Jordan 13 Retro",
		Date:       "18 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	return releases
}
