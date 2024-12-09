package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GET_AllUsers_DEV(c *gin.Context) {
	users, _ := database.GetUsers()
	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}
	utils.RouteError(c, strings.Join(usernames, "\n"))
}

func POST_UserRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	password_confirm := c.PostForm("password_confirm")
	recaptchaResponse := c.PostForm("g-recaptcha-response")

	// Verify reCAPTCHA first
	success, err := utils.VerifyRecaptcha(recaptchaResponse)
	if err != nil || !success {
		utils.RouteError(c, "please complete the reCAPTCHA verification")
		return
	}

	// Making sure passwords match
	if password != password_confirm {
		utils.RouteError(c, "passwords do not match")
		return
	}

	// No blank username/password
	if username == "" || password == "" {
		utils.RouteError(c, "blank username/password not allowed")
		return
	}

	// Attempting to add user to the database
	err = database.AddUser(username, password)
	if err != nil {
		utils.RouteError(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/?register=true")
}

func POST_UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	recaptchaResponse := c.PostForm("g-recaptcha-response")

	// Verify reCAPTCHA first
	success, err := utils.VerifyRecaptcha(recaptchaResponse)
	if err != nil || !success {
		utils.RouteError(c, "please complete the reCAPTCHA verification")
		return
	}

	// If the password matches, generate an auth token
	err = database.VerifyPassword(username, password)
	if err != nil {
		utils.RouteError(c, "there was a problem logging in, please try again and verify your password")
		return
	}
	token, err := database.GenerateAuthToken(username)
	if err != nil {
		utils.RouteError(c, "there was a problem logging in, please try again and verify your password")
		return
	}

	c.SetCookie("jwt_token", token, int(time.Hour.Seconds()), "/", "/", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func POST_UserLogout(c *gin.Context) {
	_, err := database.GetUser(c.GetString("username"))
	if err != nil {
		utils.RouteError(c, "there was a problem logging out, please try again")
		return
	}
	err = database.UpdateUserAuthToken(c.GetString("username"), "")
	if err != nil {
		utils.RouteError(c, "there was a problem logging out, please try again")
		return
	}

	c.SetCookie("jwt_token", "", 1, "/", "/", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
