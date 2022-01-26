package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func isError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type Menu struct {
	Items []interface{}
}
type Category struct {
	Name        string
	Description string
	MaxQty      int
	MinQty      int
}

type Product struct {
	Category        Category
	Name            string
	Description     string
	Sku             string
	Price           string
	Stock           string
	Type            string
	SortingPosition int
	ImageUrl        string
}

func transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

// func countCat(c chan string) {
// 	for i := 0; i < 5; i++ {
// 		c <- "Hola"
// 	}
// 	fmt.Println("done")
// }

func liaSearch(lia Menu, c chan string) {
	fmt.Println(lia.Items)

	fmt.Println("Entra aca")
	collar := lia.Items

	//var sku []string
	for i := range collar {
		var result Product
		transcode(collar[i], &result)
		//fmt.Printf("%+v\n", result)
		//if result.Name
		nameLen := len(result.Name)
		DescLen := len(result.Description)
		CatNameLen := len(result.Category.Name)
		CatDescLen := len(result.Category.Description)
		fmt.Println(nameLen)
		fmt.Println(DescLen)
		fmt.Println(CatNameLen)
		fmt.Println(CatDescLen)
		//fmt.Println(CatDescLen)
		c <- result.Sku
		// if nameLen > 1 {
		// 	//sku = append(sku, )

		// }

	}

}

func main() {
	fmt.Println("Bienvenido")
	jsonFile, err := os.Open("./tmp/menu.json")
	isError(err)
	fmt.Println("Ingresando al Menu para revisión")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result Menu
	json.Unmarshal([]byte(byteValue), &result)
	c := make(chan string)

	//go liaSearch(result, c)
	go countCat(c)
	// closer

	close(c)

	for message := range c {
		fmt.Println(message)
	}

}
