package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-takemikazuchi-api/cmd/injection"
	"go-takemikazuchi-api/configs"
	"go-takemikazuchi-api/pkg/exception"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	viperConfig.ReadInConfig()

	// Database Initialization
	databaseCredentials := &configs.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}

	databaseInstance := configs.NewDatabaseConnection(databaseCredentials)
	databaseConnection := databaseInstance.GetDatabaseConnection()

	// Config Initialization
	validatorInstance, engTranslator := configs.InitializeValidator()
	identityProvider := configs.NewIdentityProvider(viperConfig)
	mailerService := configs.NewMailerService(viperConfig)

	googleMapsInstance := configs.NewGoogleMaps(viperConfig)
	googleMapsClient := googleMapsInstance.InitializeGoogleMaps()
	// Gin Initialization
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())
	rootRouterGroup := ginEngine.Group("/")
	midtransService := configs.NewMidtransService(viperConfig)
	midtransClient := midtransService.InitializeMidtransConfiguration()
	nominatimBaseUrl := viperConfig.GetString("NOMINATIM_BASE_URL")
	nominatimHttpClient := configs.NewHttpClient(nil, &nominatimBaseUrl)
	_, initRoutesError := injection.InitializeRoutes(
		rootRouterGroup,
		databaseConnection,
		validatorInstance,
		engTranslator,
		viperConfig,
		mailerService,
		identityProvider,
		googleMapsClient,
		midtransClient,
		nominatimHttpClient,
	)
	if initRoutesError != nil {
		panic(initRoutesError)
	}
	ginError := ginEngine.Run(":8000")
	if ginError != nil {
		panic(ginError)
	}
}
