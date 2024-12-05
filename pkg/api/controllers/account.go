package controllers

import (
	"net/http"
	"os"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GET_AccountPFP(c *gin.Context) {
	path := database.GetUserPFPPath(c.Param("account_id"))
	if path == "bad" {
		utils.RouteError(c, "something went wrong grabbing account pfp")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/"+path)
}

func GET_Account(c *gin.Context) {
	c.HTML(http.StatusOK, "account.tmpl", gin.H{
		"IsLoggedIn":         true,
		"Username":           c.GetString("username"),
		"ImageUploadMessage": c.Query("imageUploadMessage"),
	})
}

func POST_AccountPFPUpload(c *gin.Context) {
	uploadPath := "./public/uploads"
	uploadPathFull := uploadPath + "/" + uuid.New().String()

	// Getting image from POST form
	img, imgHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=There+was+a+problem+uploading+the+image.")
		return
	}

	// Making sure image isn't too big
	if imgHeader.Size >= 1000000 {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=Image+too+big.")
		return
	}

	// Checking to see if it is a JPG or PNG
	imageBuffer := make([]byte, 512)
	_, err = img.Read(imageBuffer)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=There+was+a+problem+reading+the+image.")
		return
	}
	imageType := http.DetectContentType(imageBuffer)
	if imageType == "image/jpeg" {
		uploadPathFull += ".jpg"
	} else if imageType == "image/png" {
		uploadPathFull += ".png"
	} else {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=Bad+image+format.")
		return
	}

	// Try to make the public uploads directory
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=There+was+a+problem+uploading+the+image.")
		return
	}

	// Save input image
	if err := c.SaveUploadedFile(imgHeader, uploadPathFull); err != nil {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=There+was+a+problem+saving+the+image.")
		return
	}

	// Deleting old profile picture from disk if they had one
	if err = database.UserHasNoPFP(c.GetString("username")); err != nil {
		if err = database.DeleteUserPFPFromDisk(c.GetString("username")); err != nil {
			c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=Try+again")
			return
		}
	}

	// Adding image to database then updating user PFP ID
	pfpId, err := database.AddMedia(uploadPathFull[2:], c.GetString("username"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/account?imageUploadMessage=There+was+a+problem+saving+the+image.")
		return
	}
	database.UpdateUserPicture(c.GetString("username"), pfpId)

	c.Redirect(http.StatusSeeOther, "/account")
}
