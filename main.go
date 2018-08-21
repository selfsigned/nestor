package main

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	api_url = "https://api-nestor.com/"
	menu_route = "menu/"
)

type Food struct {
	ID          string `json:"_id"`
	Type        int    `json:"type"`
	Desc        string `json:"desc"`
	User        string `json:"user"`
	Name        string `json:"name"`
	V           int    `json:"__v"`
	Ingredients string `json:"ingredients"`
	Review      struct {
		NbVote int     `json:"nb_vote"`
		Note   float64 `json:"note"`
	} `json:"review"`
	NbRelease    int           `json:"nb_release"`
	Releases     []time.Time   `json:"releases"`
	Price        int           `json:"price"`
	Cost         int           `json:"cost"`
	Cold         bool          `json:"cold"`
	Spicy        bool          `json:"spicy"`
	Vegan        bool          `json:"vegan"`
	Vegetarian   bool          `json:"vegetarian"`
	NoEgg        bool          `json:"no_egg"`
	NoNuts       bool          `json:"no_nuts"`
	NoMilk       bool          `json:"no_milk"`
	NoGluten     bool          `json:"no_gluten"`
	Informations []interface{} `json:"informations"`
	ImageURL     string        `json:"image_url"`
	Ranking      int           `json:"ranking"`
	Updated      time.Time     `json:"updated"`
}

type Week_menus struct {
	Menus []struct {
		Menus []struct {
			Label  string `json:"label"`
			Price  int    `json:"price"`
			Desc   string `json:"desc"`
			ID     string `json:"_id"`
			Entree Food   `json:"entree"`
			Dish   Food   `json:"dish"`
			Dessert Food  `json:"dessert"`
			Date                time.Time `json:"date"`
			Type                int       `json:"type"`
			Soldout             bool      `json:"soldout"`
			PushMessage         string    `json:"pushMessage"`
			RestrictedAddresses []string  `json:"restrictedAddresses"`
			NoGluten            bool      `json:"no_gluten"`
			NoMilk              bool      `json:"no_milk"`
			NoNuts              bool      `json:"no_nuts"`
			NoEgg               bool      `json:"no_egg"`
			Vegetarian          bool      `json:"vegetarian"`
			Vegan               bool      `json:"vegan"`
			Spicy               bool      `json:"spicy"`
			Cold                bool      `json:"cold"`
			Quantity            int       `json:"quantity"`
		} `json:"menus"`
		Date       time.Time `json:"date"`
		Battlement []string  `json:"battlement"`
	} `json:"menus"`
	Selected  int    `json:"selected"`
	NextWeek  bool   `json:"next_week"`
	KitchenID string `json:"kitchen_id"`
}

func disp_day (day int, menu Week_menus) {
	t := menu.Menus[day].Menus[0]
	fmt.Printf("Entrée: %s Note:%.2f\nDish:%s Note:%.2f\nDessert:%s Note:%.2f\n",
		t.Entree.Name, t.Entree.Review.Note, t.Dish.Name, t.Dish.Review.Note,
		t.Dessert.Name, t.Dessert.Review.Note)
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	url := api_url + menu_route + os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "%s Error:%d\n",
			url, resp.StatusCode)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var menu Week_menus
	err = json.Unmarshal(body, &menu)
	if err != nil {
		panic(err)
	}

	println("Monday")
	disp_day(0, menu)
	println("\nTuesday")
	disp_day(1, menu)
	println("\nWednesday")
	disp_day(2, menu)
	println("\nThursday")
	disp_day(3, menu)
	println("\nFriday")
	disp_day(4, menu)
}
