package main

import (
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
	accessToken_url           = "https://oauth.fatsecret.com/connect/token" // TODO: remove
	foodSearch_url            = "https://platform.fatsecret.com/rest/foods/search/v1"
	fat_info_file_path        = "./fat_info.json"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessTokenResponse := food.AccessTokenJSON{}
	client := &http.Client{}
	ClientID := os.Getenv("FATSECRET_ClientID")
	ClientSecret := os.Getenv("FATSECRET_Client_Secret")

	if _, err := os.Stat(fat_info_file_path); err != nil {
		fmt.Println("fat_info.json does not exists")

		if accessToken_err := food.GetNewAccessToken(
			client, &accessTokenResponse, ClientID, ClientSecret,
		); accessToken_err != nil {
			log.Fatalf("Error: %v", accessToken_err)
		}

	}

	// fmt.Println("fat_info.json exists")
	fileReadError := food.ProcessSecretFile(fat_info_file_path, &accessTokenResponse)
	if fileReadError != nil {
		log.Fatalf("Error loading secret file: %v\n", fileReadError)
	}

	searchOptions := food.NewFatSecretSearchOptions(food.FSS_SearchExpression("raw chicken breast"), food.FSS_MaxResults(2), food.FSS_PageNumber(0))

	search_req, err_search_req := food.NewSearchRequest(searchOptions, foodSearch_url, accessTokenResponse)
	if err_search_req != nil {
		log.Fatalf("Error with NewSearchRequest: %v", err_search_req)
	}

	fatResponse := food.NewFatResponseHandler(food.FRH_DataRequest(client, search_req), food.FRH_NewError(), food.FRH_NewSearch())
	if err := fatResponse.ErrorUnmarshal(); err != nil {
		log.Fatalf("Error umarshall main: %v", err)
	}
	if err := fatResponse.SearchUnmarshal(); err != nil {
		log.Fatalf("Search umarshall main: %v", err)
	}

	if !reflect.DeepEqual(fatResponse.FatError, food.FatSecretError{}) { // FALSE BE GOOD
		if fatResponse.FatError.Error.Code == 13 {
			// TODO: I also need to make a new request.
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

	fmt.Println("V.2")
	fmt.Println("Total Results", fatResponse.FatSearch.Foods.TotalResults)
	fmt.Println("FOOD NUTRITION v") // get serving size
	nutritionForFood := fatResponse.FatSearch.Foods.Food[0].ParseNutritionFromFoodItem()
	fmt.Println(nutritionForFood)
	fmt.Println(fatResponse.FatSearch.Foods.Food[0].FoodDescription)
	fmt.Printf("END -- %v\n\n", 0)
}
