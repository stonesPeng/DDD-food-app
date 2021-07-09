/**
  @author: honor
  @since: 2021/7/9
  @desc: //TODO
**/
package main

import (
	"DDD-food-app/infra/auth"
	"DDD-food-app/infra/repository"
	"DDD-food-app/intefaces"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no env got")
	}
}

func main() {
	//dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	if repositories, err := repository.NewRepositories(user, password, port, host, dbname); err != nil {
		panic(err)
	} else {
		repositories.AutoMigrate()
		if redisService, er := auth.NewRedisDB(redis_host, redis_port, redis_password); er != nil {
			log.Fatal(er)
		} else {
			token := auth.NewToken()
			users := intefaces.NewUsers(repositories.User, redisService.Auth, token)
			authenticate := intefaces.NewAuthenticate(repositories.User, redisService.Auth, token)

			engine := gin.Default()

			//TODO middleware add for CORS

			engine.POST("/users", users.SaveUser)
			engine.GET("/users", users.GetUsers)
			engine.GET("/users/:user_id", users.GetUser)

			engine.POST("/login", authenticate.Login)
			engine.POST("/logout", authenticate.Logout)

			engine.POST("/refresh", authenticate.Refresh)

			app_port := os.Getenv("PORT")
			if app_port == "" {
				app_port = "8080"
			}
			log.Fatal(engine.Run(":" + app_port))
		}
	}

}
