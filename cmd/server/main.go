package main

import (
	"log"
	"net/http"

	"github.com/MatheusFontesDaGama/api-go/configs"
	_ "github.com/MatheusFontesDaGama/api-go/docs"
	"github.com/MatheusFontesDaGama/api-go/internal/entity"
	"github.com/MatheusFontesDaGama/api-go/internal/infra/database"
	"github.com/MatheusFontesDaGama/api-go/internal/infra/webserver/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Matheus Fontes da Gama
// @contact.url    https://github.com/MatheusFontesDaGama

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handler.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handler.NewUserHandler(userDB, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(LogRequest)

	r.Route("/products", func(routerProducts chi.Router) {
		routerProducts.Use(jwtauth.Verifier(configs.TokenAuth))
		routerProducts.Use(jwtauth.Authenticator)

		routerProducts.Post("/", productHandler.CreateProduct)
		routerProducts.Get("/", productHandler.GetProducts)
		routerProducts.Get("/{id}", productHandler.GetProduct)
		routerProducts.Put("/{id}", productHandler.UpdateProduct)
		routerProducts.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	http.ListenAndServe(":8000", r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		log.Printf("Request: %s %s", request.Method, request.URL.Path)
		next.ServeHTTP(response, request)
	})
}
