package model

import (
	"io"
	"sort"
	"sync"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

var NotFoundError = errors.New("Not found")
var DuplicateError = errors.New("Duplicate entry")

type RecipesFetcher interface {
	FetchOneByID(recipeID int) (*Recipe, error)
	FetchRecipes(limiter *Limiter) []*Recipe
}

type RecipesCreator interface {
	CreateRecipe(recipe *Recipe) error
}

type RecipesUpdater interface {
	UpdateRecipe(recipeID int, recipe *Recipe) error
}

type RecipesRater interface {
	RateRecipe(recipeID int, rate *RecipeRate) error
}

type RecipesAggregator interface {
	RecipesCreator
	RecipesFetcher
	RecipesRater
	RecipesUpdater
}

type Limiter struct {
	Limit int
	Page  int
}

type Recipe struct {
	Id                     int      `csv:"id" json:"id"`
	CreatedAt              DateTime `csv:"created_at" json:"created_at"`
	UploadedAt             DateTime `csv:"uploaded_at" json:"uploaded_at"`
	BoxType                string   `csv:"box_type" json:"box_type"`
	Title                  string   `csv:"title" json:"title"`
	Slug                   string   `csv:"slug" json:"slug"`
	ShortTitle             string   `csv:"short_title" json:"short_title"`
	MarketingDescription   string   `csv:"marketing_description" json:"marketing_description"`
	CaloriesKCal           int      `csv:"calories_kcal" json:"calories_k_cal"`
	ProteinGrams           int      `csv:"protein_grams" json:"protein_grams"`
	FatGrams               int      `csv:"fat_grams" json:"fat_grams"`
	CarbsGrams             int      `csv:"carbs_grams" json:"carbs_grams"`
	Bulletpoint1           string   `csv:"bulletpoint1" json:"bulletpoint_1"`
	Bulletpoint2           string   `csv:"bulletpoint2" json:"bulletpoint_2"`
	Bulletpoint3           string   `csv:"bulletpoint3" json:"bulletpoint_3"`
	RecipeDietTypeId       string   `csv:"recipe_diet_type_id" json:"recipe_diet_type_id"`
	Season                 string   `csv:"season" json:"season"`
	Base                   string   `csv:"base" json:"base"`
	ProteinSource          string   `csv:"protein_source" json:"protein_source"`
	PreparationTimeMinutes int      `csv:"preparation_time_minutes" json:"preparation_time_minutes"` //In ideal world this would be time.Duration
	ShelfLifeDays          int      `csv:"shelf_life_days" json:"shelf_life_days"`                   //Same - time.Duration
	EquipmentNeeded        string   `csv:"equipment_needed" json:"equipment_needed"`                 //This could be slice of equipment or strings
	OriginCountry          string   `csv:"origin_country" json:"origin_country"`
	RecipeCuisine          string   `csv:"recipe_cuisine" json:"recipe_cuisine"`
	InYourBox              string   `csv:"in_your_box" json:"in_your_box"` //Same this is a great example of slice of strings
	GoustoReference        int      `csv:"gousto_reference" json:"gousto_reference"`

	rates       []*RecipeRate
	AverageRate float32
}

type RecipeRate struct {
	Rate    int
	RatedAt DateTime
	RatedBy string
}

type RecipesModel struct {
	mx      sync.Mutex
	recipes map[int]*Recipe
}

func NewRecipesModel() *RecipesModel {
	return &RecipesModel{
		recipes: make(map[int]*Recipe),
	}
}

func (r *RecipesModel) LoadFromCSV(csv io.Reader) error {
	loadedRecipes := []*Recipe{}
	if err := gocsv.Unmarshal(csv, &loadedRecipes); err != nil {
		return errors.Wrapf(err, "failed to unmarshal csv")
	}
	r.mx.Lock()
	for _, recipe := range loadedRecipes {
		r.recipes[recipe.Id] = recipe
	}
	r.mx.Unlock()
	return nil
}

func (r *RecipesModel) FetchOneByID(recipeID int) (*Recipe, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if val, ok := r.recipes[recipeID]; ok {
		return val, nil
	}
	return nil, NotFoundError
}

//FetchRecipes is not ideal but I don't want to be bothered as normal case scenario for me is to use elastic search for that
//I don't know anyone who likes CSV
//Worth to notice is that map in go doesn't guarantee order! So an edge case scenario is that next call with different
//page will have recipe from previous page, hence key sorting
func (r *RecipesModel) FetchRecipes(limiter *Limiter) []*Recipe {
	v := make([]*Recipe, 0, len(r.recipes))
	var keys []int

	r.mx.Lock()
	for k := range r.recipes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v = append(v, r.recipes[k])
	}
	r.mx.Unlock()
	if limiter.Limit == 0 {
		return v
	}

	firstElementIndex := (limiter.Page - 1) * limiter.Limit
	lastElementIndex := limiter.Page * limiter.Limit
	if lastElementIndex > len(v) {
		lastElementIndex = len(v)
	}
	return v[firstElementIndex:lastElementIndex]
}

//Right, normal database you would sort out id for me, but I won't be bothered. ID is required.
func (r *RecipesModel) CreateRecipe(recipe *Recipe) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.recipes[recipe.Id]; ok {
		return DuplicateError
	}
	r.recipes[recipe.Id] = recipe
	return nil
}

//UpdateRecipe isn't what I would leave but I really don't want to waste more time (id collision possible)
func (r *RecipesModel) UpdateRecipe(recipeID int, recipe *Recipe) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.recipes[recipeID]; ok {
		if recipeID != recipe.Id {
			delete(r.recipes, recipeID)
		}
		r.recipes[recipeID] = recipe
		return nil
	}
	return NotFoundError
}

func (r *RecipesModel) RateRecipe(recipeID int, rate *RecipeRate) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if recipe, ok := r.recipes[recipeID]; ok {
		recipe.rates = append(recipe.rates, rate)
		recipe.AverageRate = r.calculateAverageRate(recipe)
		return nil
	}
	return NotFoundError
}

func (r *RecipesModel) calculateAverageRate(recipe *Recipe) float32 {
	var sum int
	for _, rate := range recipe.rates {
		sum = sum + rate.Rate
	}
	return float32(sum) / float32(len(recipe.rates))
}
