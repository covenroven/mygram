package main

import (
	"log"
	"net/http"

	"github.com/covenroven/mygram/api"
	"github.com/covenroven/mygram/api/middleware"
	"github.com/covenroven/mygram/config"
	"github.com/covenroven/mygram/database"
	"github.com/covenroven/mygram/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repos := repository.NewRepository(db)

	app := api.NewApp(db, repos)

	r := router(app)
	r.Run(":" + config.SRV_PORT)
}

// Defines route for the service
func router(app *api.App) *gin.Engine {
	r := gin.Default()

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something wrong occured",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	auth := r.Group("users")
	auth.POST("register", app.RegisterUser)
	auth.POST("login", app.LoginUser)

	users := r.Group("users")
	users.Use(middleware.AuthorizeJWT())
	users.PUT("", app.UpdateUser)
	users.DELETE("", app.DeleteCurrentUser)

	photos := r.Group("photos")
	photos.Use(middleware.AuthorizeJWT())
	photos.POST("", app.CreatePhoto)
	photos.GET("", app.GetPhotos)
	photos.PUT(":photoID", app.UpdatePhoto)
	photos.DELETE(":photoID", app.DeletePhoto)

	socials := r.Group("social_medias")
	socials.Use(middleware.AuthorizeJWT())
	socials.POST("", app.CreateSocialMedia)
	socials.GET("", app.GetSocialMedias)
	socials.PUT(":socialMediaID", app.UpdateSocialMedia)
	socials.DELETE(":socialMediaID", app.DeleteSocialMedia)

	comments := r.Group("comments")
	comments.Use(middleware.AuthorizeJWT())
	comments.POST("", app.CreateComment)
	comments.GET("", app.GetComments)
	comments.PUT(":commentID", app.UpdateComment)
	comments.DELETE(":commentID", app.DeleteComment)

	return r
}
