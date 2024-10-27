package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	food "kalmod.github.com/fat_go/internal"
)

const (
	search_url         string = "https://api.nal.usda.gov/fdc/v1/foods/search"
	foodSearch_url            = "https://platform.fatsecret.com/rest/foods/search/v1"
	fat_info_file_path        = "./fat_info.json"
)

var (
	entered_search_expression string
	entered_weight            int
	entered_display_amount    int
)

func init() {
	flag.StringVar(&entered_search_expression, "s", "broccoli", "text that will be used to search the FATSECRET database")
	flag.IntVar(&entered_weight, "w", -1, "Adjust the weight of the item. Updates macros")
	flag.IntVar(&entered_display_amount, "d", 1, "Adjust the number of food displayed")
	flag.Parse()
}

// if file exists, return true
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func CreateAndWriteClientSecret() error {
	clientID_input, clientSecret_input := ClientSecretPrompt()
	fatsecret_client_info := fmt.Sprintf("FATSECRET_ClientID=%v\nFATSECRET_Client_Secret=%v", clientID_input, clientSecret_input)
	file, err := os.OpenFile("./.env", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fatsecret_client_info)
	if err != nil {
		return err
	}
	return nil
}

func ClientSecretPrompt() (string, string) {
	var clientID_input string
	var clientSecret_input string
	fmt.Print("Please enter ClientID: ")
	fmt.Scan(&clientID_input)
	fmt.Print("\nEnter ClientSecret: ")
	fmt.Scan(&clientSecret_input)
	return clientID_input, clientSecret_input
}

func main() {

	if !fileExists("./.env") { // if file doesn't exists, create and ask for client info
		err := CreateAndWriteClientSecret()
		if err != nil {
			log.Fatalf("Creating .env file error: %v", err)
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessTokenResponse := food.AccessTokenJSON{}
	client := &http.Client{}
	ClientID := os.Getenv("FATSECRET_ClientID")
	ClientSecret := os.Getenv("FATSECRET_Client_Secret")
	if ClientID == "" && ClientSecret == "" { // if doesn't exists in env file, prompt again
		CreateAndWriteClientSecret()
		ClientID = os.Getenv("FATSECRET_ClientID")
		ClientSecret = os.Getenv("FATSECRET_Client_Secret")
	}

	if _, err := os.Stat(fat_info_file_path); err != nil {
		fmt.Println("fat_info.json does not exists")

		if accessToken_err := food.GetNewAccessToken(
			client, &accessTokenResponse, ClientID, ClientSecret,
		); accessToken_err != nil {
			log.Fatalf("Error: %v", accessToken_err)
		}

	}

	fileReadError := food.ProcessSecretFile(fat_info_file_path, &accessTokenResponse)
	if fileReadError != nil {
		log.Fatalf("Error loading secret file: %v\n", fileReadError)
	}

	// Create Search options
	searchOptions := food.NewFatSecretSearchOptions(food.FSS_SearchExpression(entered_search_expression), food.FSS_MaxResults(entered_display_amount), food.FSS_PageNumber(0))

	// Create Request using our search options
	search_req, err_search_req := food.NewSearchRequest(searchOptions, foodSearch_url, accessTokenResponse)
	if err_search_req != nil {
		log.Fatalf("Error with NewSearchRequest: %v", err_search_req)
	}

	// Send a request and receive our response handler.
	// We unmarshall the data into our search and error structs.
	fatResponse := food.NewFatResponseHandler(food.FRH_DataRequest(client, search_req), food.FRH_NewError(), food.FRH_NewSearch())
	if err := fatResponse.ErrorUnmarshal(); err != nil {
		log.Fatalf("Error umarshall main: %v", err)
	}
	if err := fatResponse.SearchUnmarshal(); err != nil {
		log.Fatalf("Search umarshall main: %v", err)
	}

	if !reflect.DeepEqual(fatResponse.FatError, food.FatSecretError{}) { // FALSE BE GOOD
		if fatResponse.FatError.Error.Code == 2 {
			fmt.Println("Bad clientSecret info")
			os.Exit(1)
		}
		if fatResponse.FatError.Error.Code == 13 {
			if accessToken_err := food.GetNewAccessToken(
				client, &accessTokenResponse, ClientID, ClientSecret,
			); accessToken_err != nil {
				log.Fatalf("Error: %v", accessToken_err)
			}
			search_req, err_search_req = food.NewSearchRequest(searchOptions, foodSearch_url, accessTokenResponse)
			if err_search_req != nil {
				log.Fatalf("Error with NewSearchRequest: %v", err_search_req)
			}
			if err := fatResponse.NewSearch(client, search_req); err != nil {
				log.Fatalf("Request after error 13: %v", err)
			}
			if err := fatResponse.ErrorUnmarshal(); err != nil {
				log.Fatalf("Error umarshall main: %v", err)
			}
			if err := fatResponse.SearchUnmarshal(); err != nil {
				log.Fatalf("Search umarshall main: %v", err)
			}
		}
	}

	if reflect.DeepEqual(fatResponse.FatSearch, food.FatSecretSearchResult{}) { // IF TRUE IT IS BAD
		fmt.Println("FatSearchResult no match")
	}

	for i := 0; i < entered_display_amount && i < len(fatResponse.FatSearch.Foods.Food); i++ {
		var nutritionForFood food.FoodNutrition
		if entered_weight > -1 {
			nutritionForFood = food.GetNewNutritionByWeight(fatResponse.FatSearch.Foods.Food[i].ParseNutritionFromFoodItem(), entered_weight)
		} else {
			nutritionForFood = fatResponse.FatSearch.Foods.Food[i].ParseNutritionFromFoodItem()
		}
		fmt.Printf("%v\n\n", nutritionForFood.PrettyPrintNutrition())
		// fmt.Println(fatResponse.FatSearch.Foods.Food[0].FoodDescription)

	}
}
