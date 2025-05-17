package routes

import (
	"github.com/gin-gonic/gin"
	"go-takemikazuchi-api/internal/job"
	"go-takemikazuchi-api/internal/transaction"
)

type PublicRoutes struct {
	routerGroup           *gin.RouterGroup
	transactionController transaction.Controller
	jobController         job.Controller
}

func NewPublicRoutes(routerGroup *gin.RouterGroup, transactionController transaction.Controller, jobController job.Controller) *PublicRoutes {
	return &PublicRoutes{routerGroup: routerGroup.Group("public"), transactionController: transactionController, jobController: jobController}
}

func (publicRoutes *PublicRoutes) Setup() {
	publicRoutes.routerGroup.GET("jobs", publicRoutes.jobController.FindAll)
	publicRoutes.routerGroup.POST("/transactions/notifications", publicRoutes.transactionController.Notification)
}
