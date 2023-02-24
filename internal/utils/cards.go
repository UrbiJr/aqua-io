package utils

import creditcard "github.com/durango/go-credit-card"

func CardMonths() []string {
	months := []string{
		"01",
		"02",
		"03",
		"04",
		"05",
		"06",
		"07",
		"08",
		"09",
		"10",
		"11",
		"12",
	}
	return months
}

func CardYears() []string {
	years := []string{
		"2023",
		"2024",
		"2025",
		"2026",
		"2027",
		"2028",
	}
	return years
}

func GetCardCompany(number string, month string, year string, cvv string) (string, error) {
	card := creditcard.Card{Number: number, Cvv: cvv, Month: month, Year: year}
	company, err := card.MethodValidate()
	if err != nil {
		return "", err
	}
	return company.Short, nil
}
