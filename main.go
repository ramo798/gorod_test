package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
)

type Pokecard struct {
	Name    string `json:"name"`
	Modelno string `json:"modelno"`
	Cardno  string `json:"cardno"`
	Pics    int    `json:"pics"`
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "living2",
		})
	})
	r.GET("/test",
		sctest)

	r.Run(":8090")

	// sctest()
}

func sctest(c *gin.Context) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://www.pokemon-card.com/deck/deck.html?deckID=kvvFfk-zYe5rz-kFFwVF").MustWaitLoad()
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
	c.JSON(200, all)
}
