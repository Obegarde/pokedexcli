package pokeapi
import(
	"time"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
)


func NewClient() (*Client, error){
	c :=&Client{
		baseUrl:"https://pokeapi.co/api/v2/location-area",
		baseClient: &http.Client{
		Timeout: time.Second * 10,
	},
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
	return nil
}

func GetPreviousLocationAreas(c *Client)error{
	if c.LastResponse.Previous == nil{
	return fmt.Errorf("Error, no Previous location found")
	}
	previousURL := *c.LastResponse.Previous


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
	return nil
	}
	err = json.Unmarshal(data,&c.LastResponse)
	if err != nil{
	return err
	}
	return nil
}
