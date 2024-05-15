package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	wikipediaAPIURL = "https://en.wikipedia.org/w/api.php"
)

type CountryInfo struct {
	Country          string
	Population       string
	Capital          string
	CapitalPopulation string
}

func fetchWikipediaData(query string) (map[string]interface{}, error) {
	resp, err := http.Get(wikipediaAPIURL + "?action=query&format=json&prop=extracts&exintro&titles=" + query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func extractCountryInfo(data map[string]interface{}) CountryInfo {
	pages := data["query"].(map[string]interface{})["pages"].(map[string]interface{})
	for _, page := range pages {
		extract := page.(map[string]interface{})["extract"].(string)
		lines := strings.Split(extract, "\n")
		country := lines[0]
		population := lines[1]
		capital := lines[2]
		capitalPopulation := lines[3]
		return CountryInfo{
			Country:          country,
			Population:       population,
			Capital:          capital,
			CapitalPopulation: capitalPopulation,
		}
	}
	return CountryInfo{}
}

func writeCSV(filename string, data []CountryInfo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Country", "Population", "Capital", "Capital Population"})

	// Write CSV rows
	for _, info := range data {
		writer.Write([]string{info.Country, info.Population, info.Capital, info.CapitalPopulation})
	}

	return nil
}

func main() {
	// Define the --csv flag
	csvFile := flag.String("csv", "countries.csv", "The name of the CSV file to output the results to")
	flag.Parse()

	countries := []string{"Germany", "France", "Italy"}
	var countryInfos []CountryInfo

	for _, country := range countries {
		data, err := fetchWikipediaData(country)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", country, err)
			continue
		}
		countryInfo := extractCountryInfo(data)
		countryInfos = append(countryInfos, countryInfo)
	}

	if err := writeCSV(*csvFile, countryInfos); err != nil {
		fmt.Printf("Error writing to CSV file: %v\n", err)
	} else {
		fmt.Printf("Data successfully written to %s\n", *csvFile)
	}
}
