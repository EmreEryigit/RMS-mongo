package controller

/* func GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

	}
}

func GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {

	}
}

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// first initialize private user for password validation
		var userPrivate model.UserPrivate
		if err := c.Bind(&userPrivate); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid request")
		}
		// validate
		validationError := validate.Struct(userPrivate)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		var count int64
		repo.Model(&model.User{}).Where("email = ?", userPrivate.Email).Count(&count)

		if count > 0 {
			defer cancel()
			return c.JSON(http.StatusConflict, "email already taken")
		}
		userPrivate.HashPassword()
		user := userPrivate.User
		repo.Save(&user)
		jwtToken, err := helper.GenerateJWT(fmt.Sprint(userPrivate.ID), *userPrivate.Name, *userPrivate.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating jwt token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = jwtToken
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			defer cancel()
			c.JSON(http.StatusInternalServerError, "error while generating jwt token")
			return err
		}
		c.JSON(http.StatusOK, user)
		defer cancel()
		return err
	}

}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var user model.UserPrivate
		var foundUser model.User
		if err := c.Bind(&user); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		result := repo.Where("email = ?", user.Email).First(&foundUser)
		if result.Error != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "user does not exist")
		}
		if foundUser.Email == nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "user not found")
		}
		isValid := foundUser.VerifyPassword(*user.Password)
		if !isValid {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid email or password")
		}
		token, err := helper.GenerateJWT(fmt.Sprint(foundUser.ID), *foundUser.Name, *foundUser.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = token
		err1 := session.Save(c.Request(), c.Response())
		if err1 != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "could not save the cookie")
		}
		c.JSON(http.StatusOK, foundUser)
		defer cancel()
		return err
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := Store.Get(c.Request(), "auth-session")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error sessions")
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error saving session")
		}
		return err
	}
}

func WhoAmI() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("current-user").(*helper.SignedDetails)
		var user model.User
		repo.Model(&model.User{}).Preload("Products").First(&user, claims.UserID)
		return c.JSON(http.StatusOK, user)
	}
} */
