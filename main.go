package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {

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

	for _, repo := range repositories {
		fmt.Println(repo.Name)

		languagesURL := repo.Languages
		languagesBody, err := getURL(languagesURL)
		if err != nil {
			log.Fatal(err)
		}

		var languages map[string]interface{}
		err = json.Unmarshal(languagesBody, &languages)
		if err != nil {
			log.Fatal("Error decoding languages:", err)
		}

		for lang, count := range languages {
			fmt.Printf("Language: %s, Count: %.f\n", lang, count)
		}

	}

	// store, err := NewPostgresStore()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v\n", store)
	// server := NewServer(":3000", store)
	// server.Run()
}
