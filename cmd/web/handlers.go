package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/kamil-alekber/deliciousness-recipes/internal/models/recipes"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.render(w, http.StatusNotFound, "404.html", app.newTemplateData(r))
		return
	}

	recipesList, err := app.recipes.ListRecipes(context.Background())
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Recipes = recipesList

	// migrate to middleware or put inside the context?
	userId, _ := r.Cookie("userid")

	if userId != nil {
		user, err := app.users.GetUser(context.Background(), userId.Value)

		if err != nil {
			app.serverError(w, err)
			return
		}

		data.User = user
	}

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) recipeView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.render(w, http.StatusNotFound, "404.html", app.newTemplateData(r))
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
	if r.Method == http.MethodGet {
		app.render(w, http.StatusOK, "create.html", app.newTemplateData(r))
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		name := r.Form.Get("name")
		description := r.Form.Get("description")
		ingredients := r.Form.Get("ingredients")
		instructions := r.Form.Get("instructions")

		cookingTime, err := strconv.Atoi(r.PostForm.Get("cook_time"))

		if err != nil {
			app.infoLog.Printf("parse form error in here: %v", err.Error())
			app.clientError(w, http.StatusBadRequest)
			return
		}

		newRecipe, err := app.recipes.CreateRecipe(context.Background(), recipes.CreateRecipeParams{
			ID:           rand.Int64(),
			Name:         name,
			Description:  description,
			Ingredients:  ingredients,
			Instructions: instructions,
			CookingTime:  int64(cookingTime),
		})

		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/recipe/view?id=%d", newRecipe.ID), http.StatusSeeOther)
		return
	}

	w.Header().Set("Allow", http.MethodPost)
	w.Header().Set("Allow", http.MethodGet)

	app.clientError(w, http.StatusMethodNotAllowed)

}
