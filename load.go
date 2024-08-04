package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func loadCheckInitialised() (bool, []string) {
	initialised := true
	var missing []string

	if hasCountryDatabase() {
		initialised = dbInitialised("COUNTRY")
		if !initialised {
			missing = append(missing, "COUNTRY")
		}
	}

	if initialised && hasASNDatabase() {
		initialised = dbInitialised("ASN")
		if !initialised {
			missing = append(missing, "ASN")
		}
	}

	if initialised && hasCityDatabase() {
		initialised = dbInitialised("CITY")
		if !initialised {
			missing = append(missing, "CITY")
		}
	}

	return initialised, missing
}

func loadData(dataToLoad []DataToLoad) {
	for _, item := range dataToLoad {
		switch item.Download.Type {
			case "CITY":	loadCities(item)
			case "ASN":		loadASNs(item)
			case "COUNTRY":	loadCountries(item)
		}
	}
}

func loadCities(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version := dbQueryMaxVersion("ip_city", dataToLoad.Version) + 1
	cities	:= []IpCity{}
	count	:= 0;

	fmt.Println("rebuilding: ip_city ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if err != nil {
			break
		}
		lat, _ := strconv.ParseFloat(record[7], 64)
		lon, _ := strconv.ParseFloat(record[8], 64)

		cities = append(cities, IpCity{ record[0], record[1], record[2], record[3], record[4], record[5], record[6], lat, lon, record[9], dataToLoad.Version, version })

		if len(cities) == 1000 {
			count += len(cities);
			fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
			fmt.Printf("Saved: %d entries\n", count)
			dbSaveCities(cities)
			cities = []IpCity{}
		}
	}

	if len(cities) > 0 {
		count += len(cities);
		fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
		fmt.Printf("Saved: %d entries\n", count)
		dbSaveCities(cities)
	}

	dbDropOld("ip_city", dataToLoad.Version, version)
}

func loadASNs(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version := dbQueryMaxVersion("ip_asn", dataToLoad.Version) + 1
	ASNs	:= []IpASN{}
	count	:= 0;

	fmt.Println("rebuilding: ip_asn ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if err != nil {
			break
		}
		asn, _ := strconv.Atoi(record[2])
		ASNs = append(ASNs, IpASN{ record[0], record[1], asn, record[3], dataToLoad.Version, version })

		if len(ASNs) == 1000 {
			count += len(ASNs);
			fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
			fmt.Printf("Saved: %d entries\n", count)
			dbSaveASNs(ASNs)
			ASNs = []IpASN{}
		}
	}

	if len(ASNs) > 0 {
		count += len(ASNs);
		fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
		fmt.Printf("Saved: %d entries\n", count)
		dbSaveASNs(ASNs)
	}

	dbDropOld("ip_asn", dataToLoad.Version, version)
}

func loadCountries(dataToLoad DataToLoad) {
	csvFile, err := os.Open(dataToLoad.Path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFileReader := csv.NewReader(csvFile)

	version		:= dbQueryMaxVersion("ip_country", dataToLoad.Version) + 1
	countries	:= []IpCountry{}
	count		:= 0;

	fmt.Println("rebuilding: ip_country ipv", dataToLoad.Version)
	fmt.Print("\033[s") // Save the cursor position
	for {
		record, err := csvFileReader.Read()
		if err != nil {
			break
		}
		countries = append(countries, IpCountry{ record[0], record[1], record[2], dataToLoad.Version, version })

		if len(countries) == 1000 {
			count += len(countries);
			fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
			fmt.Printf("Saved: %d entries\n", count)
			dbSaveCountries(countries)
			countries = []IpCountry{}
		}
	}

	if len(countries) > 0 {
		count += len(countries);
		fmt.Print("\033[u\033[K") // Restore the cursor position and clear the line
		fmt.Printf("Saved: %d entries\n", count)
		dbSaveCountries(countries)
	}

	dbDropOld("ip_country", dataToLoad.Version, version)
}

func loadDbStructure() {
	dbFile()
}