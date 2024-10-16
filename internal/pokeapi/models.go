package pokeapi
import (
	"net/http"
	
)


type PokeResponse struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	LocationArea []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Client struct{
	baseUrl string
	baseClient *http.Client
	LastResponse PokeResponse 	
	Cache	*PokeCache 

}
