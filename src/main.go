package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
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

type Children struct {
	Product  Product
	Category Category
}

func writeToFile(f *os.File, word string) {
	datawriter := bufio.NewWriter(f)

	datawriter.WriteString(word + "\n")

	//datawriter.WriteString("\n")
	datawriter.Flush()
}

func transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func liaSearch(lia Menu, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	collar := lia.Items
	for i := range collar {
		wg.Add(1)
		var result Product
		transcode(collar[i], &result)
		//fmt.Printf("%+v\n", result)
		nameLen := len(result.Name)
		DescLen := len(result.Description)
		CatNameLen := len(result.Category.Name)
		CatDescLen := len(result.Category.Description)

		if nameLen > 1 || DescLen > 1 || CatNameLen > 1 || CatDescLen > 1 {
			c <- result.Sku
		}

	}

}

func liaTail(lia Menu, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	collar := lia.Items
	for i := range collar {

		var result Product
		transcode(collar[i], &result)
		//fmt.Printf("%+v\n", result)
		nameLen := len(result.Name)
		DescLen := len(result.Description)
		CatNameLen := len(result.Category.Name)
		CatDescLen := len(result.Category.Description)

		if nameLen > 1 || DescLen > 1 || CatNameLen > 1 || CatDescLen > 1 {
			ch <- result.Sku
		}

	}

}

func showChannel(ch chan string, name string) {
	namefile := fmt.Sprintf("skuFaild%s", name)
	namefileext := "./tmp/s" + namefile + ".txt"
	jsonErr, _ := os.Create(namefileext)
	for message := range ch {
		sku := fmt.Sprintf("Sku donde se infringe la regla de longitud: %s", message)
		fmt.Println(sku)
		writeToFile(jsonErr, sku)

	}
}

func main() {
	fmt.Println("Inicio de ejecución")
	jsonFile, err := os.Open("./tmp/menu.json")
	jsonFileFull, err := os.Open("./tmp/menu_full.json")

	isError(err)
	fmt.Println("Ingresando al archivo de Menú para revisión")

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValueFull, _ := ioutil.ReadAll(jsonFileFull)
	var result Menu
	var resultFull Menu
	json.Unmarshal([]byte(byteValue), &result)
	json.Unmarshal([]byte(byteValueFull), &resultFull)
	c := make(chan string)
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(15)
	go liaSearch(result, c, &wg)
	go liaTail(resultFull, ch, &wg)
	go showChannel(c, "menu")
	go showChannel(ch, "menuFull")
	defer func() {
		wg.Wait()
		close(c)
		wg.Done()
		jsonFile.Close()
		jsonFileFull.Close()
	}()

}
