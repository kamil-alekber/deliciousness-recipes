package main

import (
	"context"
	"database/sql"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/kamil-alekber/deliciousness-recipes/internal/models/recipes"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	recipesList, err := app.recipes.ListRecipes(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Recipes = recipesList

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) recipeView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	recipe, err := app.recipes.GetRecipe(context.Background(), int64(id))

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}
	data := app.newTemplateData(r)
	data.Recipe = recipe

	app.render(w, http.StatusOK, "recipe.html", data)

}

func (app *application) recipeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	newRecipe, err := app.recipes.CreateRecipe(context.Background(), recipes.CreateRecipeParams{
		ID:          rand.Int64(),
		Name:        "New Recipe",
		Description: "Description of the new recipe",
		// Ingredients:  []string{"Ingredient 1", "Ingredient 2"},
		// Instructions: []string{"Instruction 1", "Instruction 2"},
	})

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Write([]byte{byte(newRecipe.ID)})
}
