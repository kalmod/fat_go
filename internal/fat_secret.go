package food

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type fatResponseJSON interface {
	Print() string
}

type FatSecretError struct { // JSON
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (f *FatSecretError) Print() string {
	prettyJson, err := json.MarshalIndent(f, "", " ")
	if err != nil {
		return fmt.Sprintf("Error trying to pretty print FatSecretError: %v", err)
	}
	return string(prettyJson)
}

type AccessTokenJSON struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type FatSecretFoodItem struct {
	FoodDescription string `json:"food_description"`
	FoodID          string `json:"food_id"`
	FoodName        string `json:"food_name"`
	FoodType        string `json:"food_type"`
	FoodURL         string `json:"food_url"`
	BrandName       string `json:"brand_name,omitempty"`
}

func (f *FatSecretFoodItem) ParseNutritionFromFoodItem() FoodNutrition {

	result := strings.Split(f.FoodDescription, " - ")
	macros := strings.Split(result[1], " | ")
	servingInfo := regexp.MustCompile(`(\d+)\s*(\w+)`)
	parsedServingInfo := servingInfo.FindStringSubmatch(result[0])

	servingAmount, parseIntErr := strconv.ParseInt(parsedServingInfo[1], 10, 32)
	if parseIntErr != nil {
		return FoodNutrition{}
	}

	serving := ServingSize{Amount: int(servingAmount), Metric: GetUnitType(parsedServingInfo[2])}

	macroName := regexp.MustCompile(`(\w+):`)
	macroAmount := regexp.MustCompile(`([\d\.]+)`)
	macroMetric := regexp.MustCompile(`\d([^\d\s\W]+)`)
	allMacros := []Macro{}
	for _, n := range macros {
		f, err := strconv.ParseFloat(macroAmount.FindStringSubmatch(n)[1], 64)
		if err != nil {
			fmt.Println("Could not parse macros")
			continue
		}
		allMacros = append(allMacros,
			Macro{
				Name:   macroName.FindStringSubmatch(n)[1],
				Amount: f,
				Metric: GetUnitType(macroMetric.FindStringSubmatch(n)[1]),
			})
	}
	return FoodNutrition{f.FoodName, f.FoodID, f.BrandName, serving, allMacros}
}

type FatSecretSearchResult struct { // JSON
	Foods struct {
		Food         []FatSecretFoodItem `json:"food"`
		MaxResults   string              `json:"max_results"`
		PageNumber   string              `json:"page_number"`
		TotalResults string              `json:"total_results"`
	} `json:"foods"`
}

func (f *FatSecretSearchResult) Print() string {
	prettyJson, err := json.MarshalIndent(f, "", " ")
	if err != nil {
		return fmt.Sprintf("Error pretty printing FatSecretSearchResult: %v", err)
	}
	return string(prettyJson)
}

const accessToken_url string = "https://oauth.fatsecret.com/connect/token"

func GetAccessTokenUsingClientSecret(client *http.Client, accessToken *AccessTokenJSON, clientID, clientSecret string) {
	data := "grant_type=client_credentials&scope=basic"
	getAccessToken, err := http.NewRequest(
		http.MethodPost, accessToken_url, bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	getAccessToken.Header.Add("content-type", "application/x-www-form-urlencoded")
	getAccessToken.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(getAccessToken)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	decodeErr := dec.Decode(&accessToken)
	if decodeErr != nil {
		log.Fatalf("Error: %v", decodeErr)
	}
}

func GetNewAccessToken(client *http.Client, accessToken *AccessTokenJSON, clientID, clientSecret string) error {
	GetAccessTokenUsingClientSecret(client, accessToken, clientID, clientSecret)
	jsonString, err := json.Marshal(accessToken)
	if err != nil {
		return errors.New("Failed to marshall json when getting AccessToken.")
	}
	os.WriteFile("./fat_info.json", jsonString, 0644) //TODO: What is 0644?
	return nil
}

// TODO: Return an error number - DO NOT USER
func CheckForError(body io.Reader) int {
	var genericMap map[string]interface{}
	data, err := io.ReadAll(body)
	if err != nil {
		return -1 // non json error
	}
	if err := json.Unmarshal(data, &genericMap); err != nil {
		return -1 // non json error
	}

	// I have to type assert the contents to be able to accessthe contents
	errorMap, ok := genericMap["error"]
	if !ok {
		return 1
	}
	errRespCode, ok := errorMap.(map[string]interface{})["code"]
	if !ok {
		return 1
	}

	return int(errRespCode.(float64))
}

// TODO: We'll create a master struct that will hold all my json.
// Create a parser function for each struct that returns an error.
// Pass all of these as functional arguments.
func CheckResponse(body io.Reader, jsonStructs ...fatResponseJSON) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	for i, js := range jsonStructs {
		unmarshalErr := json.Unmarshal(data, js)
		fmt.Printf("%v -- %v\n", i, unmarshalErr)
	}
	return nil
}

// TODO: Should this return an error?
func ProcessSecretFile(filepath string, accessTokenResponse *AccessTokenJSON) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteData, &accessTokenResponse)
	if err != nil {
		return err
	}
	return nil
}
