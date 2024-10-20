package food

// import "encoding/json"

type FoodSearchCriteria_USDA struct {
	DataType               []string `json:"dataType"`
	Query                  string   `json:"query"`
	GeneralSearchInput     string   `json:"generalSearchInput"`
	BrandOwner             string   `json:"brandOwner"`
	PageNumber             int      `json:"pageNumber"`
	SortBy                 string   `json:"sortBy"`
	SortOrder              string   `json:"sortOrder"`
	NumberOfResultsPerPage int      `json:"numberOfResultsPerPage"`
	PageSize               int      `json:"pageSize"`
	RequireAllWords        bool     `json:"requireAllWords"`
	FoodTypes              []string `json:"foodTypes"`
}

type FoodNutrients_USDA struct {
	NutrientID                    int     `json:"nutrientId"`
	NutrientName                  string  `json:"nutrientName"`
	NutrientNumber                string  `json:"nutrientNumber"`
	UnitName                      string  `json:"unitName"`
	DerivationCode                string  `json:"derivationCode"`
	DerivationDescription         string  `json:"derivationDescription"`
	DerivationID                  int     `json:"derivationId"`
	Value                         float64 `json:"value"`
	FoodNutrientSourceID          int     `json:"foodNutrientSourceId"`
	FoodNutrientSourceCode        string  `json:"foodNutrientSourceCode"`
	FoodNutrientSourceDescription string  `json:"foodNutrientSourceDescription"`
	Rank                          int     `json:"rank"`
	IndentLevel                   int     `json:"indentLevel"`
	FoodNutrientID                int     `json:"foodNutrientId"`
	PercentDailyValue             int     `json:"percentDailyValue,omitempty"`
}

type FoodAttributes_USDA struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

type FoodAttributeTypes_USDA struct {
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	ID             int                   `json:"id"`
	FoodAttributes []FoodAttributes_USDA `json:"foodAttributes"`
}

type Foods_USDA struct {
	FdcID                    int                       `json:"fdcId"`
	Description              string                    `json:"description"`
	DataType                 string                    `json:"dataType"`
	GtinUpc                  string                    `json:"gtinUpc"`
	PublishedDate            string                    `json:"publishedDate"`
	BrandOwner               string                    `json:"brandOwner"`
	BrandName                string                    `json:"brandName,omitempty"`
	Ingredients              string                    `json:"ingredients"`
	MarketCountry            string                    `json:"marketCountry"`
	FoodCategory             string                    `json:"foodCategory"`
	ModifiedDate             string                    `json:"modifiedDate"`
	DataSource               string                    `json:"dataSource"`
	PackageWeight            string                    `json:"packageWeight,omitempty"`
	ServingSizeUnit          string                    `json:"servingSizeUnit"`
	ServingSize              float64                   `json:"servingSize"`
	HouseholdServingFullText string                    `json:"householdServingFullText"`
	TradeChannels            []string                  `json:"tradeChannels"`
	AllHighlightFields       string                    `json:"allHighlightFields"`
	Score                    float64                   `json:"score"`
	Microbes                 []interface{}             `json:"microbes"`
	FoodNutrients            []FoodNutrients_USDA      `json:"foodNutrients"`
	FinalFoodInputFoods      []interface{}             `json:"finalFoodInputFoods"`
	FoodMeasures             []interface{}             `json:"foodMeasures"`
	FoodAttributes           []FoodAttributes_USDA     `json:"foodAttributes"`
	FoodAttributeTypes       []FoodAttributeTypes_USDA `json:"foodAttributeTypes"`
	FoodVersionIds           []interface{}             `json:"foodVersionIds"`
	SubbrandName             string                    `json:"subbrandName,omitempty"`
	ShortDescription         string                    `json:"shortDescription,omitempty"`
	PreparationStateCode     string                    `json:"preparationStateCode,omitempty"`
	Footnote                 string                    `json:"footnote,omitempty"`
}

type Search_USDA struct {
	TotalHits          int                     `json:"totalHits"`
	CurrentPage        int                     `json:"currentPage"`
	TotalPages         int                     `json:"totalPages"`
	PageList           []int                   `json:"pageList"`
	FoodSearchCriteria FoodSearchCriteria_USDA `json:"foodSearchCriteria"`
	Foods              []Foods_USDA            `json:"foods"`
	Aggregations       Aggregations_USDA       `json:"aggregations"`
}

type Aggregations_USDA struct {
	DataType  DataType_USDA `json:"dataType"`
	Nutrients interface{}   `json:"nutrients"`
}

type DataType_USDA struct {
	Branded     int `json:"Branded"`
	SurveyFNDDS int `json:"Survey (FNDDS)"`
	SRLegacy    int `json:"SR Legacy"`
	Foundation  int `json:"Foundation"`
}

func (f *Foods_USDA) GetDescription() string {
	return f.Description
}

func (f *Foods_USDA) GetEnergy() float64 {
	for _, fn := range f.FoodNutrients {
		if fn.NutrientID == 2047 {
			return fn.Value
		}
	}
	return 0.0
}
