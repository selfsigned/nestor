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

	. "github.com/logrusorgru/aurora"
)

const (
	apiURL    = "https://api-nestor.com/"
	menuRoute = "menu/"
)

var (
	// API
	zipCode = flag.Int("zip", 75017, "The postal code for the APi query")
	// Time
	flagToday    = flag.Bool("today", false, "Show today's menu")
	flagWeek     = flag.Bool("week", false, "Show all the menus for the week")
	flagTomorrow = flag.Bool("tomorrow", false, "Show tomorrow's menu")
	// Display
	flagSeen        = flag.Bool("seen", false, "Show the number of times the item has been seen")
	flagIngredients = flag.Bool("ingredients", false, "Show the ingredients (ugly)")
)

type food struct {
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

type weekMenus struct {
	Days []struct {
		Menus []struct {
			Label               string    `json:"label"`
			Price               int       `json:"price"`
			Desc                string    `json:"desc"`
			ID                  string    `json:"_id"`
			Entree              food      `json:"entree"`
			Dish                food      `json:"dish"`
			Dessert             food      `json:"dessert"`
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

func buildFoodStr(f food) string {
	str := fmt.Sprintf("%s\n\t", Black(f.Name).BgGray())
	if *flagIngredients {
		str = fmt.Sprintf("%s%s\n\t", str, f.Ingredients)
	}
	note := f.Review.Note
	switch {
	case note >= 4.1:
		str = fmt.Sprintf("%sNote: %.2f (%d votes)", str,
			Green(note), f.Review.NbVote)
	case note > 3.7:
		// Brown is Yellow, nani the fuck??!
		str = fmt.Sprintf("%sNote: %.2f (%d votes)", str,
			Brown(note), f.Review.NbVote)
	case note != 0:
		str = fmt.Sprintf("%sNote: %.2f (%d votes)", str,
			Red(note), f.Review.NbVote)
	}
	if len(f.Releases) <= 1 {
		str = fmt.Sprintf("%s %s", str, Brown("New!"))
	} else if *flagSeen {
		str = fmt.Sprintf("%s\n\tSeen %d times", str, len(f.Releases))
	}
	if f.Vegan {
		str = fmt.Sprintf("%s %s", str, Green("Vegan"))
	} else if f.Vegetarian {
		str = fmt.Sprintf("%s %s", str, Green("Vegetarian"))
	}
	if f.Cold {
		str = fmt.Sprintf("%s %s", str, Blue("Cold"))
	}
	if f.Spicy {
		str = fmt.Sprintf("%s %s", str, Red("Spicy"))
	}
	return str
}

func dispDay(day int, menu weekMenus) {
	t := menu.Days[day].Menus[0]
	if t.Soldout {
		fmt.Printf("%s\n", Red("Soldout!"))
	} else {
		println()
	}
	fmt.Printf("Entrée\t%s\nDish\t%s\nDessert\t%s\n\n",
		buildFoodStr(t.Entree),
		buildFoodStr(t.Dish),
		buildFoodStr(t.Dessert))
}

func getDayIndex(day time.Time, menu weekMenus) int {
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

func showDailyMenu(day time.Time, dayStr string, menu weekMenus) {
	d := getDayIndex(day, menu)
	if d == -1 {
		println("Error: No menu for " + dayStr)
		return
	}
	print(dayStr + "'s menu: ")
	dispDay(d, menu)
}

func showWeeklyMenu(menu weekMenus) {
	if menu.NextWeek {
		println("Showing next week's menus")
	}
	for _, v := range menu.Days {
		showDailyMenu(v.Date, v.Date.Weekday().String(), menu)
	}
}

func main() {
	flag.Parse()

	// JSON GET
	url := apiURL + menuRoute + strconv.Itoa(*zipCode)
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

	var menu weekMenus
	err = json.Unmarshal(body, &menu)
	if err != nil {
		panic(err)
	}

	// Arg handling
	if *flagWeek {
		showWeeklyMenu(menu)
	} else if *flagTomorrow {
		showDailyMenu(time.Now().AddDate(0, 0, 1), "Tomorrow", menu)
	} else if *flagToday {
		showDailyMenu(time.Now(), "Today", menu)
	} else {
		if time.Now().Hour() >= 14 {
			showDailyMenu(time.Now().AddDate(0, 0, 1), "Tomorrow", menu)
		} else {
			showDailyMenu(time.Now(), "Today", menu)
		}
	}
}
