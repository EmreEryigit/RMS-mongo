package controller

import (
	"context"
	"net/http"
	"rms/model"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMenus() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := menuCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			return c.JSON(echo.ErrInternalServerError.Code, "error while fetching menus")
		}
		var allMenus []bson.M
		defer cancel()
		err = result.All(ctx, &allMenus)
		if err != nil {
			return c.JSON(echo.ErrInternalServerError.Code, "error while fetching menus2")
		}
		return c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		menuId := c.Param("menu_id")
		var menu model.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error while fetching the menu")
		}
		return c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu model.Menu
		if err := c.Bind(&menu); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID() // try
		menu.Menu_id = menu.ID.Hex()
		validationError := validate.Struct(menu)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, insertErr := menuCollection.InsertOne(ctx, &menu)
		if insertErr != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "menu could not be saved")
		}
		defer cancel()
		return c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu model.Menu

		menuId := c.Param("menu_id")
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this menu does not exist")
		}
		if err := c.Bind(&menu); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date) {
				defer cancel()
				return c.JSON(http.StatusBadRequest, "must be a valid interval")
			}
		}
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		validationError := validate.Struct(menu)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, updateErr := menuCollection.UpdateOne(ctx, bson.M{"menu_id": menuId}, bson.D{{"$set", &menu}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)
	}
}

func inTimeSpan(start time.Time, end time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
