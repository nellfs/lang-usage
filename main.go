package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Total struct {
}

type Repository struct {
	Name      string `json:"name"`
	Languages string `json:"languages_url"`
}

type Languages struct {
	Lang any
}

type LanguagePercentage struct {
	Language   string
	Percentage float64
}

func main() {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	lastRequest, err := store.getLastRequest()
	if err != nil {
		log.Fatal(err)
	}

	newRequestID := lastRequest + 1

	// jsonRepos := `[
	// 	{
	// 	"name": "hello",
	// 	"languages_url": {
	// 		"Go": 10,
	// 		"Javascript": 20,
	// 		"Python": 4
	// 		}
	// 	},
	// 	{
	// 	"name": "world",
	// 	"languages_url": {
	// 		"Go": 15,
	// 		"Typescript": 10,
	// 		"Javascript": 3
	// 		}
	// 	}
	// 	]
	// 	`

	// var repositoriesJson []Repository

	// final := make(map[string]int)

	// if err := json.Unmarshal([]byte(jsonRepos), &repositoriesJson); err != nil {
	// 	fmt.Println("Erro ao decodificar o JSON:", err)
	// 	return
	// }

	// for _, repo := range repositoriesJson {
	// 	for lang, count := range repo.Languages {
	// 		if val, ok := final[lang]; ok {
	// 			final[lang] = val + count
	// 		} else {
	// 			final[lang] = count
	// 		}
	// 	}
	// }

	// for lang, count := range final {
	// 	fmt.Println(lang, count)
	// }

	//real code

	body, err := getURL("https://api.github.com/users/nellfs/repos")
	if err != nil {
		log.Fatal(err)
	}

	var repositories []Repository

	err = json.Unmarshal(body, &repositories)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	langMap := make(map[string]int)

	for _, repo := range repositories {

		languagesURL := repo.Languages
		languagesBody, err := getURL(languagesURL)
		if err != nil {
			log.Fatal(err)
		}

		var languages map[string]int
		err = json.Unmarshal(languagesBody, &languages)
		if err != nil {
			log.Fatal("Error decoding languages:", err)
		}

		for lang, score := range languages {
			if val, ok := langMap[lang]; ok {
				langMap[lang] = val + score
			} else {
				langMap[lang] = score
			}
		}

	}

	total := 0
	for _, score := range langMap {
		total += score
	}

	for _, value := range langMap {
		percentage := float64(value) / float64(total) * 100
		roundedPercentage := fmt.Sprintf("%.2f", percentage)

		roundedFinal, err := strconv.ParseFloat(roundedPercentage, 64)
		if err != nil {
		}

		err = store.CreateCodeReport(&CodeReport{newRequestID, 0, value, roundedFinal, time.Now()})

		if err != nil {
			log.Fatal(err)
		}
	}

	// ordered := orderByValue(percentages)

	// for _, kv := range ordered {
	// 	fmt.Printf("%s: %.2f%%\n", kv.Language, kv.Percentage)
	// }

	// err = store.CreateCodeReport(&CodeReport{0, 0, 0, 0.0, time.Now()})

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func getURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s", resp.Status)
	}

	return body, nil
}

func readJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func orderByValue(m map[string]float64) []LanguagePercentage {
	// Create a slice of key-value pairs
	var keyValuePairs []LanguagePercentage
	for key, value := range m {
		keyValuePairs = append(keyValuePairs, LanguagePercentage{key, value})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Percentage > keyValuePairs[j].Percentage
	})

	return keyValuePairs
}
