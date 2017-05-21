package model_test

import (
	"strings"
	"testing"
	"time"

	"github.com/gobonoid/svc-recipes/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestCSVString = `id,created_at,updated_at,box_type,title,slug,short_title,marketing_description,calories_kcal,protein_grams,fat_grams,carbs_grams,bulletpoint1,bulletpoint2,bulletpoint3,recipe_diet_type_id,season,base,protein_source,preparation_time_minutes,shelf_life_days,equipment_needed,origin_country,recipe_cuisine,in_your_box,gousto_reference
1,30/06/2015 17:58:00,30/06/2015 17:58:00,vegetarian,test_title,test_slug,test_short_title,"very long marketing description",401,12,35,0,a,b,c,meat,all,noodles,beef,35,4,Appetite,Great Britain,asian,"lots, of, stuff",59
2,30/06/2015 17:58:00,30/06/2015 17:58:00,gourmet,Tamil Nadu Prawn Masala,tamil-nadu-prawn-masala,,"Tamil Nadu is a state on the eastern coast of the southern tip of India. Curry from there is particularly famous and it's easy to see why. This one is brimming with exciting contrasting tastes from ingredients like chilli powder, coriander and fennel seed",524,12,22,0,Vibrant & Fresh,"Warming, not spicy",Curry From Scratch,fish,all,pasta,seafood,40,4,Appetite,Great Britain,italian,"king prawns, basmati rice, onion, tomatoes, garlic, ginger, ground tumeric, red chilli powder, ground cumin, fresh coriander, curry leaves, fennel seeds",58
3,30/06/2015 17:58:00,30/06/2015 17:58:00,vegetarian,Umbrian Wild Boar Salami Ragu with Linguine,umbrian-wild-boar-salami-ragu-with-linguine,,"This delicious pasta dish comes from the Italian region of Umbria. It has a smoky and intense wild boar flavour which combines the earthy garlic, leek and onion flavours, while the chilli flakes add a nice deep aroma. Enjoy within 5-6 days of delivery.",609,17,29,0,,,,meat,all,pasta,pork,35,4,Appetite,Great Britain,british,,1
4,30/06/2015 17:58:00,30/06/2015 17:58:00,gourmet,Tenderstem and Portobello Mushrooms with Corn Polenta,tenderstem-and-portobello-mushrooms-with-corn-polenta,,"One for those who like their veggies with a slightly spicy kick. However, those short on time, be warned ' this is a time-consuming dish, but if you're willing to spend a few extra minutes in the kitchen, the fresh corn mash is extraordinary and worth a t",508,28,20,0,,,,vegetarian,all,,cheese,50,4,None,Great Britain,british,,56
5,30/06/2015 17:58:00,30/06/2015 17:58:00,vegetarian,Fennel Crusted Pork with Italian Butter Beans,fennel-crusted-pork-with-italian-butter-beans,,"A classic roast with a twist. The pork loin is marinated in rosemary, fennel seeds and chilli flakes then teamed with baked potato wedges and butter beans in tomato sauce. Enjoy within 5-6 days of delivery.",511,11,62,0,A roast with a twist,Low fat & high protein,With roast potatoes,meat,all,beans/lentils,pork,45,4,Pestle & Mortar (optional),Great Britain,british,"pork tenderloin, potatoes, butter beans, garlic, fennel seeds, medium onion, chilli flakes, fresh rosemary, tomatoes, vegetable stock cube",55
6,01/07/2015 17:58:00,01/07/2015 17:58:00,gourmet,Pork Chilli,pork-chilli,,"Succulent pork tenderloin and feathery white bean and parsnip mash mingle with feisty cumin seeds and tangy leek in this lighter, less conventional take on a British classic. Welcome to the outer limits of food!",401,12,35,0,,,,meat,all,,pork,35,4,Appetite,Great Britain,asian,,60
7,02/07/2015 17:58:00,02/07/2015 17:58:00,vegetarian,Courgette Pasta Rags,courgette-pasta-rags,,"Kick-start the new year with some get-up and go with this lean green vitality machine. Protein-packed chicken and mineral-rich kale are blended into a smooth, nut-free version of pesto; creating the ultimate composition of nutrition and taste",524,12,22,0,,,,meat,all,,chicken,40,4,Appetite,Great Britain,british,,59
8,03/07/2015 17:58:00,03/07/2015 17:58:00,vegetarian,Homemade Eggs & Beans,homemade-egg-beans,,"A Goustofied British institution, learn how to make beautifully golden breaded chicken escalopes drizzled in homemade garlic butter and served atop fluffy potato and broccoli mash.",609,17,29,0,,,,meat,all,,eggs,35,3,Appetite,Great Britain,italian,,2
9,04/07/2015 17:58:00,04/07/2015 17:58:00,gourmet,Grilled Jerusalem Fish,grilled-jerusalem-fish,,"I love this super healthy fish dish, it contains a punch from zingy ginger, a kick from chili and a salty sweet balance from soy sauce and mirim. A cleansing and restorative meal, great for body and soul.",508,28,20,0,,,,meat,all,,fish,50,4,Appetite,Great Britain,mediterranean,,57
10,05/07/2015 17:58:00,05/07/2015 17:58:00,gourmet,Pork Katsu Curry,pork-katsu-curry,,"Comprising all the best bits of the classic American number and none of the mayo, this is a warm & tasty chicken and bulgur salad with just a hint of Scandi influence. A beautifully summery medley of flavours and textures",511,11,62,0,,,,meat,all,,pork,45,4,Appetite,Great Britain,mexican,,56`

func Test_LoadFromCSV_ValidCSV(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	err := recipesModel.LoadFromCSV(strings.NewReader(TestCSVString))
	require.NoError(t, err)
	testRecipe, err := recipesModel.FetchOneByID(1)
	require.NoError(t, err)
	assert := assert.New(t)
	assert.Equal(1, testRecipe.Id)

	expectedCreatedAt, _ := time.Parse("02/01/2006 15:04:05", "30/06/2015 17:58:00")
	assert.True(testRecipe.CreatedAt.Equal(expectedCreatedAt), "CreatedAt is not what was expected")

	assert.Equal("vegetarian", testRecipe.BoxType)
	assert.Equal("test_title", testRecipe.Title)
	assert.Equal("test_slug", testRecipe.Slug)
	assert.Equal("test_short_title", testRecipe.ShortTitle)
	assert.Equal("very long marketing description", testRecipe.MarketingDescription)
	assert.Equal(401, testRecipe.CaloriesKCal)
	assert.Equal(12, testRecipe.ProteinGrams)
	assert.Equal(35, testRecipe.FatGrams)
	assert.Equal(0, testRecipe.CarbsGrams)
	assert.Equal("a", testRecipe.Bulletpoint1)
	assert.Equal("b", testRecipe.Bulletpoint2)
	assert.Equal("c", testRecipe.Bulletpoint3)
	assert.Equal("meat", testRecipe.RecipeDietTypeId)
	assert.Equal("all", testRecipe.Season)
	assert.Equal("noodles", testRecipe.Base)
	assert.Equal("beef", testRecipe.ProteinSource)
	assert.Equal(35, testRecipe.PreparationTimeMinutes)
	assert.Equal(4, testRecipe.ShelfLifeDays)
	assert.Equal("Appetite", testRecipe.EquipmentNeeded)
	assert.Equal("Great Britain", testRecipe.OriginCountry)
	assert.Equal("asian", testRecipe.RecipeCuisine)
	assert.Equal("lots, of, stuff", testRecipe.InYourBox)
	assert.Equal(59, testRecipe.GoustoReference)
}

func TestRecipesModel_LoadFromCSV_InvalidCSV(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	err := recipesModel.LoadFromCSV(strings.NewReader(`this, is crap, '""""""""""'`))
	require.Error(t, err)
}

func TestRecipesModel_FetchRecipeByID_NotFound(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	testRecipe, err := recipesModel.FetchOneByID(11234123)
	assert.Nil(t, testRecipe)
	assert.Equal(t, model.NotFoundError, err)
}

func TestRecipesModel_FetchRecipes(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	err := recipesModel.LoadFromCSV(strings.NewReader(TestCSVString))
	require.NoError(t, err)
	recipes := recipesModel.FetchRecipes(&model.Limiter{Limit: 2, Page: 1})
	assert.Equal(t, 2, len(recipes))

	recipes = recipesModel.FetchRecipes(&model.Limiter{})
	assert.Equal(t, 10, len(recipes))
}

func TestRecipesModel_CreateRecipe(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	require.NoError(t, recipesModel.CreateRecipe(&model.Recipe{Id: 1}))
	recipe, err := recipesModel.FetchOneByID(1)
	require.NoError(t, err)
	assert.Equal(t, 1, recipe.Id)
	assert.Equal(t, model.DuplicateError, recipesModel.CreateRecipe(&model.Recipe{Id: 1}))
}

func TestRecipesModel_RateRecipe(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	require.NoError(t, recipesModel.CreateRecipe(&model.Recipe{Id: 1}))
	require.NoError(t, recipesModel.RateRecipe(1, &model.RecipeRate{Rate: 5}))
	recipe, err := recipesModel.FetchOneByID(1)
	require.NoError(t, err)
	assert.Equal(t, float32(5), recipe.AverageRate)

	require.NoError(t, recipesModel.RateRecipe(1, &model.RecipeRate{Rate: 6}))
	recipe, err = recipesModel.FetchOneByID(1)
	require.NoError(t, err)
	assert.Equal(t, float32(5.5), recipe.AverageRate)
}

func TestRecipesModel_UpdateRecipe(t *testing.T) {
	recipesModel := model.NewRecipesModel()
	require.NoError(t, recipesModel.CreateRecipe(&model.Recipe{Id: 1}))
	require.NoError(t, recipesModel.UpdateRecipe(1, &model.Recipe{Id: 1, CaloriesKCal: 5}))
	assert.Equal(t, model.NotFoundError, recipesModel.UpdateRecipe(2, &model.Recipe{Id: 1, CaloriesKCal: 5}))

}
