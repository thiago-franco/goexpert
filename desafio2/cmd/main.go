package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CEPData struct {
	cep          string `json:"cep"`
	state        string `json:"state"`
	city         string `json:"city"`
	neighborhood string `json:"neighborhood"`
	street       string `json:"street"`
}

func (c *CEPData) UnmarshalJSON(data []byte) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	c.cep = fmt.Sprintf("%v", tmp["cep"])
	state, ok := tmp["state"]
	if ok {
		c.state = fmt.Sprintf("%v", state)
		c.city = fmt.Sprintf("%v", tmp["city"])
		c.neighborhood = fmt.Sprintf("%v", tmp["neighborhood"])
		c.street = fmt.Sprintf("%v", tmp["street"])
	} else if uf, ok := tmp["uf"]; ok {
		c.state = fmt.Sprintf("%v", uf)
		c.city = fmt.Sprintf("%v", tmp["localidade"])
		c.neighborhood = fmt.Sprintf("%v", tmp["bairro"])
		c.street = fmt.Sprintf("%v", tmp["logradouro"])
	}
	return nil
}

func main() {
	url1 := "https://brasilapi.com.br/api/cep/v1/22631450"
	url2 := "http://viacep.com.br/ws/22631450/json/"

	c1 := make(chan CEPData)
	c2 := make(chan CEPData)

	go makeRequest(url1, c1)
	go makeRequest(url2, c2)

	select {
	case cep1 := <-c1:
		fmt.Printf("API: %v\n", url1)
		fmt.Printf("CEP: %v", cep1)
	case cep2 := <-c2:
		fmt.Printf("API: %v\n", url2)
		fmt.Println(cep2)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}

}

func makeRequest(url string, c chan CEPData) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var cep CEPData
	json.Unmarshal(body, &cep)

	c <- cep
}
