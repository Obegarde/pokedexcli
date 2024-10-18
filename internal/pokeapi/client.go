package pokeapi
import(
	"time"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"sync"
	"math/rand"
)


func NewClient() (*Client, error){
	newcache,err := NewCache(5)
	if err!= nil{
		return nil, err
	}
	c :=&Client{
		baseUrl:"https://pokeapi.co/api/v2/location-area",
		Cache: newcache,
		baseClient: &http.Client{
		Timeout: time.Second * 10,
	},
		CaughtPokemon: make(map[string][]byte),
	
}
	return c, nil
}

func GetBaseLocationAreas(c *Client)error{
	req, err := http.NewRequest("GET",c.baseUrl,nil)
	if err !=nil{
		return err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil{
	return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
	return nil
	}
	err = c.Cache.Add(c.baseUrl,data)
	if err != nil{
	return err
	}
	err = json.Unmarshal(data,&c.LastResponse)
	if err != nil{
	return err
	}
	return nil


}

func GetNextLocationAreas(c *Client)error{
	if c.LastResponse.Next == nil{
	return fmt.Errorf("Error, no Next location found")
	}
	nextURL := *c.LastResponse.Next
	entry, exists := c.Cache.Get(nextURL)
	if exists{	
		err := json.Unmarshal(entry.val,&c.LastResponse)
	if err != nil{
	return err
	}	
	}else{

	req, err := http.NewRequest("GET",nextURL,nil)
	if err != nil{
	return err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil{
	return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
	return nil
	}
	err = json.Unmarshal(data,&c.LastResponse)
	if err != nil{
	return err
	}	
	err = c.Cache.Add(nextURL,data)
	if err != nil{
	return err
	}
	}

	return nil
}

func GetPreviousLocationAreas(c *Client)error{
	if c.LastResponse.Previous == nil{
	return fmt.Errorf("Error, no Previous location found")
	}

	previousURL := *c.LastResponse.Previous
	entry, exists := c.Cache.Get(previousURL)
	if exists{	
		err := json.Unmarshal(entry.val,&c.LastResponse)
	if err != nil{
	return err
	}	
	}else{
	

	req, err := http.NewRequest("GET",previousURL,nil)
	if err != nil{
	return err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil{
	return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
	return err 
	}
	err = json.Unmarshal(data,&c.LastResponse)
	if err != nil{
	return err
	}
	err = c.Cache.Add(previousURL, data)
		if err != nil{
		return err
		}
	}
	return nil
}

func CatchPokemon(c *Client)error{
	lastPokemonName := c.LastPokemonResponse.Name
	pokemonExperience := c.LastPokemonResponse.BaseExperience
	fmt.Println("Throwing a pokeball at " + lastPokemonName + "...")
	if (rand.Float32()* 1000) > float32(pokemonExperience){
	 fmt.Println(lastPokemonName + " has been caught!")
	cachedPokemonEntry, exists := c.Cache.Get("https://pokeapi.co/api/v2/pokemon/"+lastPokemonName)
		if !exists{
		return fmt.Errorf("Error: Pokemon not found in cache")
		}
	c.CaughtPokemon[lastPokemonName] = cachedPokemonEntry.val
	return nil
	}
	fmt.Println(lastPokemonName+" escaped!")
	return nil
}


func GetPokemon(c *Client, pokemon string)error{
	pokemonURL := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	entry, exists := c.Cache.Get(pokemonURL)
	if exists{
		err := json.Unmarshal(entry.val, &c.LastPokemonResponse)
	
	if err != nil{
	return err
	}
	}else{
	req, err := http.NewRequest("GET",pokemonURL,nil)	
	if err != nil{
	return err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil{
	return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
	return err 
	}
	err = json.Unmarshal(data,&c.LastPokemonResponse)
	if err != nil{
	return err
	}
	err = c.Cache.Add(pokemonURL, data)
		if err != nil{
		return err
		}
	}
	return nil
}



func ExploreLocation(c *Client, name string)error{
	ExploreURL := "https://pokeapi.co/api/v2/location-area/" + name
	entry, exists := c.Cache.Get(ExploreURL)
	if exists{	
		err := json.Unmarshal(entry.val,&c.LastExploreResponse)
	if err != nil{
	return err
	}	
	}else{
	

	req, err := http.NewRequest("GET",ExploreURL,nil)
	if err != nil{
	return err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil{
	return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
	return err 
	}
	err = json.Unmarshal(data,&c.LastExploreResponse)
	if err != nil{
	return err
	}
	err = c.Cache.Add(ExploreURL, data)
		if err != nil{
		return err
		}
	}
	return nil
}

type cacheEntry struct {
	createdAt time.Time
	val	[]byte 
}

type PokeCache struct{	
	mu	sync.Mutex
	entries map[string]cacheEntry
	interval int64 
}


func NewCache(inputinterval int64)(*PokeCache,error){
	cache := PokeCache{
		entries: make(map[string]cacheEntry),
		interval: inputinterval,
	}
	go cache.reapLoop()

	return &cache, nil
}

func (c *PokeCache) Add(key string, inputVal []byte) error{
	if key == ""{
	return fmt.Errorf("No key given")
	}
	currentTime := time.Now()
	newCache := cacheEntry{
	createdAt: currentTime,
	val: inputVal,
	}
	c.mu.Lock()
	c.entries[key] = newCache	
	c.mu.Unlock()
	return nil
}

func (c *PokeCache) Get(key string)(cacheEntry, bool){
	if key == ""{
		fmt.Println("Key is zero value. Please enter a valid key.")
	}
	c.mu.Lock()
	entry, exists := c.entries[key]
	c.mu.Unlock()
	if exists{	
		return entry,true
	}
	return cacheEntry{}, false
}

func (c *PokeCache) reapLoop(){
	ticker := time.NewTicker(time.Duration(c.interval) * time.Minute)
	for range ticker.C {		
		c.mu.Lock()
		for key,val := range c.entries{
			if(time.Since(val.createdAt) > time.Duration(c.interval)* time.Minute){
				delete(c.entries, key)
				}
			}
		c.mu.Unlock()
		}
	
	}


