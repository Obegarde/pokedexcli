package main
import(
	"fmt"
	"os"
	"bufio"
	"strings"
)

type cliCommand struct{
	name string
	description string
	callback func() error
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

func commandHelp() error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage: \n\n help: Displays a help message\n exit: Exit the Pokedex")
	return nil
}

func commandExit() error{	
	return fmt.Errorf("exit") 
}

func main(){
	fmt.Println("Welcome To the Pokedex!")
	err := initializeCommands()
	if err != nil{
	fmt.Println("Error: initializing, exiting")
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
	err = cmd.callback()
	if err !=nil{
		if err.Error() == "exit"{
			break
			}
			fmt.Println("Command error: ", err)

		}
	}
}

