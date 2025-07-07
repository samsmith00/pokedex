package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/samsmith00/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cgf *Config, cache *pokecache.Cache, input []string, u *User) error
	Config
}

type Config struct {
	Offset int
}

type Location struct {
	Name string `json:"name"`
}

type LocationResponse struct {
	Results []Location `json:"results"`
}

type CaughtPokemon struct {
	Name  string
	Color string
}

type User struct {
	Pokemon map[string]CaughtPokemon
}

type InspectData struct {
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}

type Stats struct {
	Hp             int
	Attack         int
	Defence        int
	SpecialAttat   int
	SpecialDefence int
	Speed          int
}

type Types struct {
	PokemonTypes []string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	commands := initCommands()
	cfg := &Config{}
	cache := pokecache.NewCache(500 * time.Second)
	user := &User{
		make(map[string]CaughtPokemon),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			err = fmt.Errorf("%s", err)
			fmt.Println(err.Error())
		}

		input := scanner.Text()

		parsedInput := cleanInput(input)

		if parsedInput[0] == commands["exit"].name {
			commands["exit"].callback(cfg, cache, parsedInput, user)
			break
		}
		if parsedInput[0] == commands["help"].name {
			commands["help"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["map"].name {
			commands["map"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["mapb"].name {
			commands["mapb"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["explore"].name {
			commands["explore"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["catch"].name {
			commands["catch"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["listpk"].name {
			commands["listpk"].callback(cfg, cache, parsedInput, user)
		}
		if parsedInput[0] == commands["inspect"].name {
			commands["inspect"].callback(cfg, cache, parsedInput, user)
		}
	}
}

func cleanInput(text string) []string {
	newText := strings.TrimSpace(text)
	newText = strings.ToLower(newText)

	slice := strings.Split(newText, " ")

	return slice
}

func initCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"usage": {
			name:        "usage",
			description: " ",
			callback:    nil,
		},
		"map": {
			name:        "map",
			description: "Display next locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"listpk": {
			name:        "listpk",
			description: "List my Pokemon",
			callback:    commandList,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect Pokemon",
			callback:    commandInspect,
		},
	}
}

func commandMap(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", cfg.Offset)

	cachedval, ok := cache.Get(pokeURL)
	if ok {
		var locres LocationResponse

		fmt.Println("----------cashed used-----------")

		json.Unmarshal(cachedval, &locres)

		for _, loc := range locres.Results {
			fmt.Println(loc.Name)
		}
		cfg.Offset += 20
		return nil
	}

	res, err := http.Get(pokeURL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	cache.Add(pokeURL, body)

	var locres LocationResponse
	json.Unmarshal(body, &locres)

	for _, loc := range locres.Results {
		fmt.Println(loc.Name)
	}

	cfg.Offset += 20
	return nil
}

func commandMapb(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	if cfg.Offset == 0 {
		return nil
	}

	cfg.Offset -= 20

	pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", cfg.Offset)

	cachedVal, ok := cache.Get(pokeURL)
	fmt.Println(ok)
	if ok {
		var locRes LocationResponse

		fmt.Println("----------CASHED USED-----------")

		json.Unmarshal(cachedVal, &locRes)

		for _, loc := range locRes.Results {
			fmt.Println(loc.Name)
		}
		return nil
	}

	res, err := http.Get(pokeURL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	cache.Add(pokeURL, body)

	var locRes LocationResponse

	json.Unmarshal(body, &locRes)

	for _, loc := range locRes.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExit(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	fmt.Println("Welcome to the Pokedex!")
	commands := initCommands()

	fmt.Println("Usage: " + commands["usage"].description)
	fmt.Println()
	fmt.Println("help: " + commands["help"].description)
	fmt.Println("exit: " + commands["exit"].description)
	return nil
}

func commandExplore(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", pi[1])

	res, err := http.Get(pokeURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data map[string]any

	json.Unmarshal(body, &data)

	encounters, ok := data["pokemon_encounters"].([]any)
	if !ok {
		return err
	}

	for _, e := range encounters {
		encounterMap := e.(map[string]any)
		pokemon := encounterMap["pokemon"].(map[string]any)
		fmt.Println("- " + pokemon["name"].(string))
	}

	return nil
}

func commandCatch(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", pi[1])
	fmt.Printf("Throwing a Pokeball at %s...\n", pi[1])

	res, err := http.Get(pokeURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var data map[string]any

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, &data)

	captureChance, ok := data["capture_rate"].(float64)
	if !ok {
		return err
	}

	capturePercent := (float64(captureChance) / 255) * 100

	chance := rand.Intn(101)
	wasCaught := false

	if chance >= int(capturePercent) {
		fmt.Printf("%s escaped!\n", pi[1])
	} else {
		fmt.Printf("%s was caught!\n", pi[1])
		wasCaught = true
	}

	colorMap, ok := data["color"].(map[string]any)
	if !ok {
		return nil
	}

	pokeColor := colorMap["name"].(string)

	if wasCaught {
		user.AddPokemon(pi[1], pokeColor)
	}

	return nil
}

func (u *User) AddPokemon(name string, color string) {
	u.Pokemon[name] = CaughtPokemon{
		name,
		color,
	}
}

func commandList(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	for name, poke := range user.Pokemon {
		fmt.Printf("%s: {%s, %s}\n", name, poke.Name, poke.Color)
	}
	return nil
}

func commandInspect(cfg *Config, cache *pokecache.Cache, pi []string, user *User) error {
	pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pi[1])

	res, err := http.Get(pokeURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data map[string]any

	json.Unmarshal(body, &data)

	inspectData := InspectData{
		Stats: make(map[string]int),
	}

	inspectData.Height = int(data["height"].(float64))
	inspectData.Weight = int(data["weight"].(float64))

	statsArr := data["stats"].([]any)
	for _, s := range statsArr {
		stat := s.(map[string]any)
		name := stat["stat"].(map[string]any)["name"].(string)
		value := int(stat["base_stat"].(float64))
		inspectData.Stats[name] = value
	}

	typesArr := data["types"].([]any)
	for _, t := range typesArr {
		typeName := t.(map[string]any)["type"].(map[string]any)["name"].(string)
		inspectData.Types = append(inspectData.Types, typeName)
	}

	fmt.Printf("\nHeight: %d\n", inspectData.Height)
	fmt.Printf("Weight: %d\n", inspectData.Weight)
	fmt.Println("Stats:")
	for k, v := range inspectData.Stats {
		fmt.Printf("  - %s: %d\n", k, v)
	}
	fmt.Println("Types:")
	for _, typ := range inspectData.Types {
		fmt.Printf("  - %s\n", typ)
	}

	return nil
}
