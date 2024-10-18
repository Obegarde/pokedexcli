package main
import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/obegarde/pokedexcli/internal/pokeapi"
	"encoding/json"
)

type cliCommand struct{
	name string
	description string
	callback func(c *pokeapi.Client, name string) error
}
var cliCommands map[string]cliCommand



func initializeCommands() error{
cliCommands = map[string]cliCommand{
	"help":{
		name: "help",
		description: "Displays a help message\n",
		callback: commandHelp,
	},
	"exit":{
		name: "exit",
		description: "Exit the Pokedex\n",
		callback: commandExit,
	},
	"map":{
		name:"map",
		description:"Displays the name of 20 locations in Pokemon or the next 20 locations\n",
		callback:commandMap,
		},
	"mapb":{
		name:"mapb",
		description:"Displays the previous 20 locations\n",
		callback:commandMapb,
		},
	"explore":{
		name:"explore",
		description:"Displays the Pokemon found in the area",
		callback:commandExplore,
		},
	"catch":{
		name:"catch",
		description:"Try and catch a pokemon by name",
		callback:commandCatch,
		},
	"inspect":{
		name:"inspect",
		description:"Inspect a pokemon you have caught",
		callback:commandInspect,
		},
	"pokedex":{
		name:"pokedex",
		description:"Displays all the pokemon you have caught",
		callback:commandPokedex,
		},

}
	
	return nil
}

func getInput() (string, error){	
	var reader = bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil{
		return "",err
	}
	command = strings.TrimSpace(command)
	return command, nil
	
}

func commandExplore(c *pokeapi.Client, name string) error{
	err := pokeapi.ExploreLocation(c, name)
	if err != nil{
		return err
	}
	printExplorePokemon(c)
	return nil
}


func printExplorePokemon(c *pokeapi.Client)error{	
	for _, pokemonEncounter := range c.LastExploreResponse.PokemonEncounters{
		fmt.Println(" - " + pokemonEncounter.Pokemon.Name)
	
		
	}
	return nil
}
func commandCatch(c *pokeapi.Client, name string)error{
	err := pokeapi.GetPokemon(c,name)
	if err != nil {
	return err
	}
	err = pokeapi.CatchPokemon(c)
	if err != nil{
	return err
	}
	return nil
}

func commandPokedex(c *pokeapi.Client, pokemon string)error{
	fmt.Println("Your Pokedex")
	if len(c.CaughtPokemon) < 1{
		return fmt.Errorf("No pokemon caught")
	}
	for key,_ := range c.CaughtPokemon{
		fmt.Println(" - " + key)
	}
	return nil
}

func commandInspect(c *pokeapi.Client, pokemon string)error{
	if data, exists := c.CaughtPokemon[pokemon]; exists{
		pokemonToInspect := pokeapi.PokemonResponse{}		
		err := json.Unmarshal(data,&pokemonToInspect)
	if err != nil{
	return err
	}
		fmt.Println("Name: " + pokemonToInspect.Name)
		fmt.Printf("Height: %v\n",pokemonToInspect.Height)
		fmt.Printf("Weight: %v\n",pokemonToInspect.Weight)
		fmt.Println("Stats:")
		for _, statEntry := range pokemonToInspect.Stats{
			fmt.Printf(" -%v: %v\n", statEntry.Stat.Name, statEntry.BaseStat)
		}
		fmt.Println("Types: ")
		for _, slot := range pokemonToInspect.Types{
			
			fmt.Println("-"+ slot.Type.Name)
			
		}
		return nil
	
	}else{
		fmt.Println("you have not caught that pokemon")
		return nil
	}	
}

func commandHelp(c *pokeapi.Client,name string) error{
	fmt.Println("Welcome to PokeHelp!")
	for key, value := range cliCommands{
		fmt.Printf("%v: %v\n",key, value.description)
	}
	return nil
}

func commandExit(c *pokeapi.Client,name string) error{	
	return fmt.Errorf("exit") 
}
func commandMap(c *pokeapi.Client, name string) error{
	
	if (c.LastResponse.Next == nil ){
		err := pokeapi.GetBaseLocationAreas(c)

		if err != nil{
		return err 
		}
	}else{
		err := pokeapi.GetNextLocationAreas(c)
		if err != nil{
		return err
		}
	}
	
	printLocations(c)	

	return nil
}


func commandMapb(c *pokeapi.Client, name string) error{
	if c.LastResponse.Previous == nil{
		return fmt.Errorf("No Previous Areas found. Please use Map") 
		}else{
		err := pokeapi.GetPreviousLocationAreas(c)
		if err != nil{
		return err
		}
		printLocations(c)
	}
	return nil
}

func printLocations(c *pokeapi.Client)error{	
	for _, Area := range c.LastResponse.LocationArea{
	fmt.Println(Area.Name)	
	}
	return nil
}


func main(){
	fmt.Println("Welcome To the Pokedex!")
	err := initializeCommands()
	if err != nil{
	fmt.Println("Error: initializing, exiting")
	return
	}
	c, err := pokeapi.NewClient()
	if err != nil{
	fmt.Printf("New Client Error: %s", err)
	return
	}
	for{
	input, err := getInput()
	if err != nil{
		fmt.Println(err)
		continue
	}
	inputs := strings.Split(input," ")
	if len(inputs)== 1{
		inputs = append(inputs,"pallet-town-area")
		}
	cmd, exists := cliCommands[inputs[0]]	
	if !exists{
		fmt.Println("Unknown command")
		continue
		}
	err = cmd.callback(c,inputs[1])
	if err !=nil{
		if err.Error() == "exit"{
			break
			}
			fmt.Println("Command error: ", err)

		}
	}
}

