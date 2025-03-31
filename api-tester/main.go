package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

type User struct {
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
}

type ResponseData struct {
	Source string `json:"source"`
	Users  []User `json:"users"`
}

func generateRandomName(r *rand.Rand) string {
	names := []string{"John", "Alice", "Bob", "Charlie", "Eve"}
	return names[r.Intn(len(names))]
}

func generateRandomSecondName(r *rand.Rand) string {
	secondNames := []string{"Doe", "Smith", "Johnson", "Williams", "Brown"}
	return secondNames[r.Intn(len(secondNames))]
}

func createUser(url string, r *rand.Rand) {
	defer wg.Done()

	user := User{
		Name:       generateRandomName(r),
		SecondName: generateRandomSecondName(r),
	}

	userData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter o usuário para JSON: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userData))
	if err != nil {
		log.Printf("Erro ao fazer requisição POST: %v\n", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Usuário criado: %s | Status: %s\n", user.Name, resp.Status)
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	url := "http://localhost:8020/users"

	users := make([]int, 5000)

	for range users {
		time.Sleep(time.Millisecond * 20)
		wg.Add(1)
		go createUser(url, r)
	}
	wg.Wait()

}
