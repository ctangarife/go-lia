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

type Child struct {
	Children []Children
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
	Children        []Children
}

type Children struct {
	Category        Category
	Name            string
	Description     string
	Sku             string
	Price           string
	Stock           string
	Type            string
	SortingPosition int
	MaxLimit        int
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

		//fmt.Printf("%+v\n", result.Children)
		nameLen := len(result.Name)
		DescLen := len(result.Description)
		CatNameLen := len(result.Category.Name)
		CatDescLen := len(result.Category.Description)

		if nameLen > 39 || DescLen > 179 || CatNameLen > 30 || CatDescLen > 179 {
			long := fmt.Sprintf("NameLen: %d, DescLen: %d, CatNameLen : %d, CatDescLen: %d, SKU: %s ", nameLen, DescLen, CatNameLen, CatDescLen, result.Sku)
			c <- long
		}

		for x := range result.Children {

			//fmt.Println(result.Children[x].Name)
			//fmt.Println(val)
			nameChildLen := len(result.Children[x].Name)
			CatNameChildLen := len(result.Children[x].Category.Name)
			if nameChildLen > 39 || CatNameChildLen > 39 {
				skuC := "SKUP:" + result.Sku + " SKUC:" + result.Children[x].Sku
				longC := fmt.Sprintf("NameLen: %d,CatNameLen : %d, SKU: %s ", nameChildLen, CatNameChildLen, skuC)
				c <- longC
			}
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

		if nameLen > 39 || DescLen > 179 || CatNameLen > 30 || CatDescLen > 179 {
			long := fmt.Sprintf("NameLen: %d, DescLen: %d, CatNameLen : %d, CatDescLen: %d, SKU: %s ", nameLen, DescLen, CatNameLen, CatDescLen, result.Sku)
			ch <- long
		}

		for x := range result.Children {

			//fmt.Println(result.Children[x].Name)
			//fmt.Println(val)
			nameChildLen := len(result.Children[x].Name)
			CatNameChildLen := len(result.Children[x].Category.Name)
			if nameChildLen > 39 || CatNameChildLen > 39 {
				skuC := "SKUP:" + result.Sku + " SKUC:" + result.Children[x].Sku
				longC := fmt.Sprintf("NameLen: %d, CatNameLen : %d, SKU: %s ", nameChildLen, CatNameChildLen, skuC)
				ch <- longC
			}
		}

	}

}

func showChannel(ch chan string, name string) {
	namefile := fmt.Sprintf("skuFaild%s", name)
	namefileext := "./tmp/s" + namefile + ".txt"
	jsonErr, _ := os.Create(namefileext)
	for message := range ch {
		sku := message
		//fmt.Println(sku)
		writeToFile(jsonErr, sku)

	}
}

func main() {
	fmt.Println("Inicio de ejecución")
	jsonFile, err := os.Open("./tmp/menu.json")
	jsonFileFull, err := os.Open("./tmp/menu_text.json")

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
	wg.Add(2)
	go liaSearch(result, c, &wg)
	go liaTail(resultFull, ch, &wg)
	go showChannel(c, "menu")
	go showChannel(ch, "menuFull")
	defer func() {
		wg.Wait()
		close(c)
		close(ch)
		wg.Done()
		jsonFile.Close()
		jsonFileFull.Close()
	}()

}
