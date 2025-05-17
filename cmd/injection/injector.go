//go:build wireinject
// +build wireinject

package injection

import (
	"github.com/gin-gonic/gin"
	universalTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/spf13/viper"
	"go-takemikazuchi-api/configs"
	categoryFeature "go-takemikazuchi-api/internal/category"
	jobFeature "go-takemikazuchi-api/internal/job"
	jobApplicationFeature "go-takemikazuchi-api/internal/job_application"
	jobResourceFeature "go-takemikazuchi-api/internal/job_resource"
	reviewFeature "go-takemikazuchi-api/internal/review"
	"go-takemikazuchi-api/internal/routes"
	"go-takemikazuchi-api/internal/storage"
	transactionFeature "go-takemikazuchi-api/internal/transaction"
	userFeature "go-takemikazuchi-api/internal/user"
	userAddressFeature "go-takemikazuchi-api/internal/user_address"
	validatorFeature "go-takemikazuchi-api/internal/validator"
	withdrawalFeature "go-takemikazuchi-api/internal/withdrawal"
	workerFeature "go-takemikazuchi-api/internal/worker"
	workerResourceFeature "go-takemikazuchi-api/internal/worker_resource"
	workerWalletFeature "go-takemikazuchi-api/internal/worker_wallet"
	"googlemaps.github.io/maps"
	"gorm.io/gorm"
)

var routeSet = wire.NewSet(
	ProvidePublicRoutes,
	ProvideAuthenticationRoutes,
	ProvideProtectedRoutes,
)

func ProvidePublicRoutes(routerGroup *gin.RouterGroup, transactionController transactionFeature.Controller, jobController jobFeature.Controller) *routes.PublicRoutes {
	publicRoutes := routes.NewPublicRoutes(routerGroup, transactionController, jobController)
	publicRoutes.Setup()
	return publicRoutes
}

func ProvideAuthenticationRoutes(routerGroup *gin.RouterGroup, userController userFeature.Controller) *routes.AuthenticationRoutes {
	authenticationRoutes := routes.NewAuthenticationRoutes(routerGroup, userController)
	authenticationRoutes.Setup()
	return authenticationRoutes
}

func ProvideProtectedRoutes(routerGroup *gin.RouterGroup,
	categoryController categoryFeature.Controller,
	jobController jobFeature.Controller,
	jobApplicationController jobApplicationFeature.Controller,
	workerController workerFeature.Controller,
	transactionController transactionFeature.Controller,
	reviewController reviewFeature.Controller,
	withdrawalController withdrawalFeature.Controller,
	viperConfig *viper.Viper) *routes.ProtectedRoutes {
	protectedRoutes := routes.NewProtectedRoutes(routerGroup, categoryController, jobController, viperConfig, workerController, transactionController, jobApplicationController, reviewController, withdrawalController)
	protectedRoutes.Setup()
	return protectedRoutes
}

var userSet = wire.NewSet(
	userFeature.NewRepository,
	wire.Bind(new(userFeature.Repository), new(*userFeature.RepositoryImpl)),
	userFeature.NewService,
	wire.Bind(new(userFeature.Service), new(*userFeature.ServiceImpl)),
	userFeature.NewHandler,
	wire.Bind(new(userFeature.Controller), new(*userFeature.Handler)),
)

var userAddressSet = wire.NewSet(
	userAddressFeature.NewUserAddressRepository,
	wire.Bind(new(userAddressFeature.Repository), new(*userAddressFeature.RepositoryImpl)),
)

var categorySet = wire.NewSet(
	categoryFeature.NewRepository,
	wire.Bind(new(categoryFeature.Repository), new(*categoryFeature.RepositoryImpl)),
	categoryFeature.NewService,
	wire.Bind(new(categoryFeature.Service), new(*categoryFeature.ServiceImpl)),
	categoryFeature.NewHandler,
	wire.Bind(new(categoryFeature.Controller), new(*categoryFeature.Handler)),
)

var jobSet = wire.NewSet(
	jobFeature.NewRepository,
	wire.Bind(new(jobFeature.Repository), new(*jobFeature.RepositoryImpl)),
	jobFeature.NewService,
	wire.Bind(new(jobFeature.Service), new(*jobFeature.ServiceImpl)),
	jobFeature.NewHandler,
	wire.Bind(new(jobFeature.Controller), new(*jobFeature.Handler)),
)

var workerSet = wire.NewSet(
	workerFeature.NewRepository,
	wire.Bind(new(workerFeature.Repository), new(*workerFeature.RepositoryImpl)),
	workerFeature.NewService,
	wire.Bind(new(workerFeature.Service), new(*workerFeature.ServiceImpl)),
	workerFeature.NewHandler,
	wire.Bind(new(workerFeature.Controller), new(*workerFeature.Handler)),
)

var workerResourceSet = wire.NewSet(
	workerResourceFeature.NewRepository,
	wire.Bind(new(workerResourceFeature.Repository), new(*workerResourceFeature.RepositoryImpl)),
)

var workerWalletSet = wire.NewSet(
	workerWalletFeature.NewRepository,
	wire.Bind(new(workerWalletFeature.Repository), new(*workerWalletFeature.RepositoryImpl)),
)

var jobApplicationSet = wire.NewSet(
	jobApplicationFeature.NewRepository,
	wire.Bind(new(jobApplicationFeature.Repository), new(*jobApplicationFeature.RepositoryImpl)),
	jobApplicationFeature.NewService,
	wire.Bind(new(jobApplicationFeature.Service), new(*jobApplicationFeature.ServiceImpl)),
	jobApplicationFeature.NewHandler,
	wire.Bind(new(jobApplicationFeature.Controller), new(*jobApplicationFeature.Handler)),
)

var transactionSet = wire.NewSet(
	transactionFeature.NewRepository,
	wire.Bind(new(transactionFeature.Repository), new(*transactionFeature.RepositoryImpl)),
	transactionFeature.NewService,
	wire.Bind(new(transactionFeature.Service), new(*transactionFeature.ServiceImpl)),
	transactionFeature.NewHandler,
	wire.Bind(new(transactionFeature.Controller), new(*transactionFeature.Handler)),
)

var reviewSet = wire.NewSet(
	reviewFeature.NewRepository,
	wire.Bind(new(reviewFeature.Repository), new(*reviewFeature.RepositoryImpl)),
	reviewFeature.NewService,
	wire.Bind(new(reviewFeature.Service), new(*reviewFeature.ServiceImpl)),
	reviewFeature.NewHandler,
	wire.Bind(new(reviewFeature.Controller), new(*reviewFeature.Handler)),
)

var withdrawalSet = wire.NewSet(
	withdrawalFeature.NewRepository,
	wire.Bind(new(withdrawalFeature.Repository), new(*withdrawalFeature.RepositoryImpl)),
	withdrawalFeature.NewService,
	wire.Bind(new(withdrawalFeature.Service), new(*withdrawalFeature.ServiceImpl)),
	withdrawalFeature.NewHandler,
	wire.Bind(new(withdrawalFeature.Controller), new(*withdrawalFeature.Handler)),
)

var jobResourceSet = wire.NewSet(
	jobResourceFeature.NewRepository,
	wire.Bind(new(jobResourceFeature.Repository), new(*jobResourceFeature.RepositoryImpl)),
)

var validatorServiceSet = wire.NewSet(
	validatorFeature.NewService,
	wire.Bind(new(validatorFeature.Service), new(*validatorFeature.ServiceImpl)),
)

// wire.go
func InitializeRoutes(
	ginRouterGroup *gin.RouterGroup,
	dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator universalTranslator.Translator,
	viperConfig *viper.Viper,
	mailerService *configs.MailerService,
	identityProvider *configs.IdentityProvider,
	googleMapsClient *maps.Client,
	midtransClient *snap.Client,
) (*routes.ApplicationRoutes, error) {
	wire.Build(
		wire.Struct(new(routes.ApplicationRoutes), "*"),
		routeSet,
		userSet,
		jobSet,
		categorySet,
		jobResourceSet,
		jobApplicationSet,
		workerSet,
		workerWalletSet,
		userAddressSet,
		transactionSet,
		workerResourceSet,
		reviewSet,
		withdrawalSet,
		validatorServiceSet,
		storage.ProvideFileStorage, // Fungsi untuk memilih implementasi yang sesuai
	)
	return nil, nil
}
