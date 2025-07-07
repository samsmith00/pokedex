package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cgf *Config) error
	Config
}

type Config struct {
	Next     *string
	Previous *string
}

type Location struct {
	Name string `json:"name"`
}

type LocationResponse struct {
	Results []Location `json:"results"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	commands := initCommands()
	cfg := &Config{}

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
			commands["exit"].callback(cfg)
			break
		}
		if parsedInput[0] == commands["help"].name {
			commands["help"].callback(cfg)
		}
		if parsedInput[0] == commands["map"].name {
			commands["map"].callback(cfg)
		}
		if parsedInput[0] == commands["mapb"].name {
			commands["mapb"].callback(cfg)
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
	}
}

func commandMap(cfg *Config) error {
	// Check if Next url exists
	nextExists := cfg.Next

	var res *http.Response
	var err error

	if nextExists != nil {
		resp, err := http.Get(*nextExists)
		if err != nil {
			return err
		}
		res = resp
	} else {
		res, err = http.Get("https://pokeapi.co/api/v2/location-area")
		if err != nil {
			return err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var data map[string]any
	json.Unmarshal(body, &data)

	if val, ok := data["next"].(string); ok {
		cfg.Next = &val
	}

	if val, ok := data["previous"].(string); ok {
		cfg.Previous = &val
	}

	var locRes LocationResponse

	json.Unmarshal(body, &locRes)

	for _, loc := range locRes.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *Config) error {
	previousExists := cfg.Previous

	var res *http.Response
	var err error

	if previousExists != nil {
		resp, err := http.Get(*previousExists)
		if err != nil {
			return err
		}
		res = resp
	} else {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var data map[string]any
	json.Unmarshal(body, &data)

	if val, ok := data["previous"].(string); ok {
		cfg.Previous = &val
	}

	if val, ok := data["previous"].(string); ok {
		cfg.Previous = &val
	}

	var locRes LocationResponse

	json.Unmarshal(body, &locRes)

	for _, loc := range locRes.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	commands := initCommands()

	fmt.Println("Usage: " + commands["usage"].description)
	fmt.Println()
	fmt.Println("help: " + commands["help"].description)
	fmt.Println("exit: " + commands["exit"].description)

	return nil
}

func displayLocations(locations []interface{}) {
	for _, itm := range locations {
		loc, ok := itm.(map[string]any)
		if !ok {
			continue
		}

		if name, ok := loc["name"].(string); ok {
			fmt.Println(name)
		}
	}
}
