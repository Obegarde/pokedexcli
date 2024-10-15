package main
import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/obegarde/pokedexcli/internal/pokeapi"
)

type cliCommand struct{
	name string
	description string
	callback func(*pokeapi.Client) error
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

func commandHelp(_ *pokeapi.Client) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage: \n\n help: Displays a help message\n exit: Exit the Pokedex")
	return nil
}

func commandExit(_ *pokeapi.Client) error{	
	return fmt.Errorf("exit") 
}
func commandMap(c *pokeapi.Client) error{
	
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


func commandMapb(c *pokeapi.Client) error{
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
	command, err := getInput()
	if err != nil{
		fmt.Println(err)
		continue
	}
	cmd, exists := cliCommands[command]	
	if !exists{
		fmt.Println("Unknown command")
		continue
		}
	err = cmd.callback(c)
	if err !=nil{
		if err.Error() == "exit"{
			break
			}
			fmt.Println("Command error: ", err)

		}
	}
}

