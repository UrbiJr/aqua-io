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
	// TODO: fetch the actual releases from stockx.go
	var releases []Release
	url, _ := url.Parse("https://images.stockx.com/360/MSCHF-Big-Red-Boot/Images/MSCHF-Big-Red-Boot/Lv2/img01.jpg?fit=fill&bg=FFFFFF&w=300&h=214&auto=compress&trim=color&dpr=2&updated_at=1675338467&q=60")
	stockxLink, _ := url.Parse("https://stockx.com/mschf-big-red-boot")
	r := Release{
		Title:      "MSCHF Red Boot",
		Date:       "16 Feb",
		StockXLink: stockxLink,
		ImageURL:   url,
	}
	releases = append(releases, r)

	return releases
}
