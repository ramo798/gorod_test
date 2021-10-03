package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type Pokecard struct {
	Name    string `json:"name"`
	Modelno string `json:"modelno"`
	Cardno  string `json:"cardno"`
	Pics    int    `json:"pics"`
}

type DeckInfo struct {
	DeckName      string `json:"deck_name"`
	DeckExp       string `json:"deck_exp"`
	MainPokemon   string `json:"main_pokemon"`
	KeyPokemon1   string `json:"key_pokemon1"`
	KeyPokemon2   string `json:"key_pokemon2"`
	KeyPokemon3   string `json:"key_pokemon3"`
	MainSymbol    string `json:"main_symbol"`
	MainCardNo    string `json:"main_card_no"`
	MainCardModel string `json:"main_crad_model"`
}

type PostData struct {
	CardList []Pokecard `json:"card_list"`
	DeckInfo DeckInfo   `json:"deck_info"`
}

func main() {
	url := "FV1FFV-KxIEoS-k5VdVf"

	cardlist := getDeckList(url)

	info := DeckInfo{
		DeckName:      "ビクティニブラッキーミュウミュウ",
		DeckExp:       "",
		MainPokemon:   "victini",
		KeyPokemon1:   "mewtwo",
		KeyPokemon2:   "umbreon",
		KeyPokemon3:   "vikavolt",
		MainSymbol:    "fire",
		MainCardNo:    "013/070",
		MainCardModel: "S5R",
	}

	Postdata := PostData{
		CardList: cardlist,
		DeckInfo: info,
	}

	fmt.Println(Postdata)

	json, err := json.Marshal(Postdata)
	if err != nil {
		fmt.Println(err)
	}

	os.Stdout.Write(json)

	content := []byte(json)
	ioutil.WriteFile("14.json", content, os.ModePerm)
}

func getDeckList(deckcode string) []Pokecard {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://www.pokemon-card.com/deck/deck.html?deckID=" + deckcode).MustWaitLoad()
	page.Timeout(5 * time.Second)

	var all []Pokecard

	if page.MustElement("#cardListView").MustHas("div.Grid_item") {
		outer := page.MustElements("div.Grid_item")

		for k := 0; k < len(outer); k++ {
			innner := outer[k].MustElements("tr")
			for i := 0; i < len(innner); i++ {
				// カードナンバーなどの取得処理
				cardnoP, err := innner[i].Attribute("id")
				if err != nil {
					continue
				}
				if cardnoP == nil {
					continue
				}
				cardnoId := *cardnoP
				cardno := strings.Replace(cardnoId, "txtView", "#cardName", -1)
				var cardname string
				cardnameRaW, err := outer[k].MustElement(cardno).HTML()
				if err != nil {
					cardname = ""
				} else {
					cardname = cardnameRaW[strings.Index(cardnameRaW, ">")+1:]
				}
				cardname = strings.Replace(cardname, "</a>", "", -1)
				cardname = strings.Replace(cardname, "<br>", " ", -1)
				// fmt.Println(cardname)

				// 枚数の取得処理
				picId := strings.Replace(cardnoId, "txtView", "#txtNumView", -1)
				picRaW := outer[k].MustElement(picId)
				pic := picRaW.MustText()
				// fmt.Println(pic)

				// 構造体に入れる処理
				var cardinfo Pokecard
				poke := strings.Split(cardname, " ")
				// fmt.Println(poke)
				if len(poke) > 1 {
					cardinfo.Name = poke[0]
					cardinfo.Modelno = poke[1]
					cardinfo.Cardno = poke[2]
				} else {
					cardinfo.Name = poke[0]
					cardinfo.Modelno = "other"
					cardinfo.Cardno = strings.Replace(cardnoId, "txtView_", "", -1)
				}
				picInt, err := strconv.Atoi(pic)
				if err != nil {
					i = 0
				}
				cardinfo.Pics = picInt
				// fmt.Println(cardinfo)
				all = append(all, cardinfo)
			}

		}

		// fmt.Println(all)

	}

	return all

}
