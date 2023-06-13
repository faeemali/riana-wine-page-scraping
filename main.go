package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func canIgnore(data string) bool {
	if strings.Trim(data, " \t\r\n") == "" {
		return true
	}

	if strings.Contains(data, "function(") || strings.HasPrefix(data, "var ") {
		return true
	}

	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s base_dir\n", os.Args[0])
		os.Exit(1)
	}

	baseDir := os.Args[1]
	file, err := os.Open(fmt.Sprintf("%s/index.html", baseDir))
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing directory: %s\n", baseDir)

	const (
		WaitForDescription    int = 0
		ProcessDescription        = 1
		TastingRoomDetails        = 2
		RedCultivars              = 3
		Specialities              = 4
		RestaurantDescription     = 5
		Outdoor                   = 6
		WaitForBusinessInfo       = 7
		BusinessInfo              = 8
		Times                     = 9
		RestaurantsAndFood        = 10
		WhiteCultivars            = 11
	)
	state := WaitForDescription

	firstRun := 0

	var writeFile *os.File
	tokenizer := html.NewTokenizer(file)
	for {
		tok := tokenizer.Next()
		if tok == html.ErrorToken {
			fmt.Printf("Got error token. Exiting\n")
			break
		}

		if tok == html.TextToken {
			token := tokenizer.Token()
			data := token.Data

			if canIgnore(data) {
				continue
			}

			dataLower := strings.ToLower(data)
			switch state {
			case WaitForDescription:
				if dataLower == "description" || strings.Contains(dataLower, "white cultivars") {
					if dataLower == "description" {
						writeFile, err = os.Create(fmt.Sprintf("%s/description.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating description file: %v\n", err)
							os.Exit(1)
						}
						state = ProcessDescription
					} else {
						writeFile, err = os.Create(fmt.Sprintf("%s/white_cultivars.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating white_cultivars file: %v\n", err)
							os.Exit(1)
						}

						state = WhiteCultivars
					}
				}
			case ProcessDescription:
				if strings.Contains(dataLower, "tasting room details") ||
					strings.Contains(dataLower, "white cultivars") {
					writeFile.Close()

					if strings.Contains(dataLower, "tasting room details") {
						writeFile, err = os.Create(fmt.Sprintf("%s/tasting_room_details.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating tasting room details file: %v\n", err)
							os.Exit(1)
						}

						state = TastingRoomDetails
					} else if strings.Contains(dataLower, "white cultivars") {
						writeFile, err = os.Create(fmt.Sprintf("%s/white_cultivars.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating white_cultivars file: %v\n", err)
							os.Exit(1)
						}

						state = WhiteCultivars
					}
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				//fmt.Printf("Writing: %s\n", data)
				size, err := writeFile.Write([]byte(formatted))
				if err != nil {
					fmt.Printf("Error writing description: (size=%d) %v\n", size, err)
					os.Exit(1)
				}

			case TastingRoomDetails:
				if strings.Contains(dataLower, "red cultivars") || strings.Contains(dataLower, "white cultivars") {
					writeFile.Close()

					if strings.Contains(dataLower, "red cultivars") {
						writeFile, err = os.Create(fmt.Sprintf("%s/red_cultivars.txt", baseDir))
						if err != nil {
							fmt.Printf("Error opening red_cultivars file: %v\n", err)
							os.Exit(1)
						}
						state = RedCultivars
					} else {
						writeFile, err = os.Create(fmt.Sprintf("%s/white_cultivars.txt", baseDir))
						if err != nil {
							fmt.Printf("Error opening white_cultivars file: %v\n", err)
							os.Exit(1)
						}
						state = WhiteCultivars
					}
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case WhiteCultivars:
				if strings.Contains(dataLower, "red cultivars") {
					writeFile.Close()
					writeFile, err = os.Create(fmt.Sprintf("%s/red_cultivars.txt", baseDir))
					if err != nil {
						fmt.Printf("Error opening red_cultivars file: %v\n", err)
						os.Exit(1)
					}
					state = RedCultivars
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case RedCultivars:
				if strings.Contains(dataLower, "specialities") ||
					strings.Contains(dataLower, "restaurant description") ||
					strings.Contains(dataLower, "restaurants and food") ||
					strings.Contains(dataLower, "business info") ||
					strings.Contains(dataLower, "photos") {
					writeFile.Close()

					if strings.Contains(dataLower, "specialities") {
						writeFile, err = os.Create(fmt.Sprintf("%s/specialities.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating specialities file: %v\n", err)
							os.Exit(1)
						}
						state = Specialities
					} else if strings.Contains(dataLower, "restaurants and food") {
						writeFile, err = os.Create(fmt.Sprintf("%s/restaurants_and_food.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating restaurants_and_food file: %v\n", err)
							os.Exit(1)
						}
						state = RestaurantsAndFood
					} else if strings.Contains(dataLower, "business info") {
						writeFile, err = os.Create(fmt.Sprintf("%s/business_info.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating business_info file: %v\n", err)
							os.Exit(1)
						}
						state = BusinessInfo
					} else if strings.Contains(dataLower, "photos") {
						state = WaitForBusinessInfo
					} else {
						writeFile, err = os.Create(fmt.Sprintf("%s/restaurant_description.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating restaurant description file: %v\n", err)
							os.Exit(1)
						}
						state = RestaurantDescription
					}

					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case Specialities:
				if strings.Contains(dataLower, "restaurant description") || strings.Contains(dataLower, "photos") {
					writeFile.Close()

					if strings.Contains(dataLower, "restaurant description") {
						writeFile, err = os.Create(fmt.Sprintf("%s/restaurant_description.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating restaurant description file: %v\n", err)
							os.Exit(1)
						}

						state = RestaurantDescription
					} else {
						state = WaitForBusinessInfo
					}
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case RestaurantDescription:
				if strings.Contains(dataLower, "outdoor") || strings.Contains(dataLower, "restaurants and food") {
					writeFile.Close()

					if strings.Contains(dataLower, "outdoor") {
						writeFile, err = os.Create(fmt.Sprintf("%s/outdoor.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating outdoor file: %v\n", err)
							os.Exit(1)
						}

						state = Outdoor
					} else {
						writeFile, err = os.Create(fmt.Sprintf("%s/restaurants_and_food.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating restaurants_and_food file: %v\n", err)
							os.Exit(1)
						}
						state = RestaurantsAndFood
					}
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case Outdoor:
				if strings.Contains(dataLower, "photos") {
					writeFile.Close()
					state = WaitForBusinessInfo
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case WaitForBusinessInfo:
				if strings.Contains(dataLower, "business info") {
					writeFile, err = os.Create(fmt.Sprintf("%s/business_info.txt", baseDir))
					if err != nil {
						fmt.Printf("Error creating business_info file: %v\n", baseDir)
						os.Exit(1)
					}

					state = BusinessInfo
					continue
				} else if strings.Contains(dataLower, "map") {
					os.Exit(1)
				}

			case BusinessInfo:
				if strings.Contains(dataLower, "today") || strings.Contains(dataLower, "map") {
					writeFile.Close()

					if strings.Contains(dataLower, "today") {
						writeFile, err = os.Create(fmt.Sprintf("%s/times.txt", baseDir))
						if err != nil {
							fmt.Printf("Error creating times file: %v\n", err)
							os.Exit(1)
						}

						firstRun = 1
						state = Times
					} else {
						os.Exit(1)
					}
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case Times:
				if strings.Contains(dataLower, "map") {
					writeFile.Close()
					os.Exit(1)
				}

				if firstRun != 0 {
					writeFile.Write([]byte("Today\n"))
					firstRun = 0
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			case RestaurantsAndFood:
				if strings.Contains(dataLower, "photos") {
					writeFile.Close()
					state = WaitForBusinessInfo
					break
				}

				formatted := fmt.Sprintf("%s\n", data)
				writeFile.Write([]byte(formatted))

			default:
				continue
			}
		}
	}
}
