package main

import (
	"fmt"
	"os"
	"time"

	"flag"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	//api
	API       = "https://api-nestor.com/"
	MenuRoute = "menu/"
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
	Days []struct {
		Menus []struct {
			Label               string    `json:"label"`
			Price               int       `json:"price"`
			Desc                string    `json:"desc"`
			ID                  string    `json:"_id"`
			Entree              Food      `json:"entree"`
			Dish                Food      `json:"dish"`
			Dessert             Food      `json:"dessert"`
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

func disp_day(day int, menu Week_menus) {
	t := menu.Days[day].Menus[0]
	fmt.Printf("Entrée: %s Note:%.2f\nDish:%s Note:%.2f\nDessert:%s Note:%.2f\n\n",
		t.Entree.Name, t.Entree.Review.Note, t.Dish.Name, t.Dish.Review.Note,
		t.Dessert.Name, t.Dessert.Review.Note)
}

func get_day_index(day time.Time, menu Week_menus) int {
	if day.Weekday() == 0 || day.Weekday() > 5 {
		return -1
	}
	for i, v := range menu.Days {
		if v.Date.YearDay() == day.YearDay() {
			return i
		}
	}
	return -1
}

func show_daily_menu(day time.Time, day_str string, menu Week_menus) {
	d := get_day_index(day, menu)
	if d == -1 {
		println("Error: No menu for" + day_str)
		return
	}
	println(day_str + "'s menu: ")
	disp_day(d, menu)
}

func show_weekly_menu(menu Week_menus) {
	if menu.NextWeek {
		println("Showing next week's menus")
	}
	for _, v := range menu.Days {
		show_daily_menu(v.Date, v.Date.Weekday().String(), menu)
	}
}

func main() {
	// Parsing
	var (
		zip_code   = flag.Int("zip", 75017, "The postal code for the APi query")
		p_today    = flag.Bool("today", true, "Show today's menu")
		p_week     = flag.Bool("week", false, "Show all the menus for the week")
		p_tomorrow = flag.Bool("tomorrow", false, "Show tomorrow's menu")
	)
	flag.Parse()

	// JSON GET
	url := API + MenuRoute + strconv.Itoa(*zip_code)
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

	// Arg handling
	if *p_week {
		show_weekly_menu(menu)
	} else if *p_tomorrow {
		show_daily_menu(time.Now().AddDate(0, 0, 1), "Tomorrow", menu)
	} else if *p_today {
		show_daily_menu(time.Now(), "Today", menu)
	}
}
