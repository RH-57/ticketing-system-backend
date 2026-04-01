package routes

import (
	"backend-golang-api/controllers"
	"backend-golang-api/middlewares"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.185:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Service is up and running!",
		})
	})

	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/refresh", controllers.RefreshToken)

		admin := api.Group("", middlewares.AuthMiddleware())
		{
			dashboard := admin.Group("/dashboard")
			{
				dashboard.GET("/stats", controllers.GetDashboardStats)
			}

			users := admin.Group("/users")
			{
				users.GET("", controllers.ShowUsers)
				users.POST("", controllers.CreateUser)
				users.GET("/:id", controllers.FindUserById)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
				users.PUT("/change-password", controllers.ChangePassword)
			}

			departments := admin.Group("/departments")
			{
				departments.GET("", controllers.GetAllDepartment)
			}

			branches := admin.Group("/branches")
			{
				branches.GET("", controllers.ShowBranches)
				branches.POST("", controllers.CreateBranch)
				branches.GET("/:id", controllers.FindBranchById)
				branches.PUT("/:id", controllers.UpdateBranch)
				branches.DELETE("/:id", controllers.DeleteBranch)

				branches.GET("/:id/divisions", controllers.ShowDivisionsByBranch)
				branches.POST("/:id/divisions", controllers.CreateDivision)
				branches.PUT("/:id/divisions/:divisionId", controllers.UpdateDivision)
				branches.DELETE("/:id/divisions/:divisionId", controllers.DeleteDivision)

				branches.GET("/:id/divisions/:divisionId/departments", controllers.ShowDepartmentByDivision)
				branches.POST("/:id/divisions/:divisionId/departments", controllers.CreateDepartment)
				branches.PUT("/:id/divisions/:divisionId/departments/:departmentId", controllers.UpdateDepartment)
				branches.DELETE("/:id/divisions/:divisionId/departments/:departmentId", controllers.DeleteDepartment)
			}

			employees := admin.Group("/employees")
			{
				employees.GET("", controllers.ShowEmployees)
				employees.POST("", controllers.CreateEmployee)
				employees.GET("/:id", controllers.ShowEmployeeById)
				employees.PUT("/:id", controllers.UpdateEmployee)
				employees.DELETE("/:id", controllers.DeleteEmployee)
			}

			categories := admin.Group("/categories")
			{
				categories.GET("", controllers.ShowCategories)
				categories.POST("", controllers.CreateCategory)
				categories.GET("/:id", controllers.ShowCategoryById)
				categories.PUT("/:id", controllers.UpdateCategory)
				categories.DELETE("/:id", controllers.DeleteCategory)

				categories.GET("/:id/sub-categories", controllers.ShowSubCategoriesByCategory)
				categories.POST("/:id/sub-categories", controllers.CreateSubCategory)
				categories.GET("/:id/sub-categories/:subCategoryId", controllers.ShowSubCategoryByID)
				categories.PUT("/:id/sub-categories/:subCategoryId", controllers.UpdateSubCategory)
				categories.DELETE("/:id/sub-categories/:subCategoryId", controllers.DeleteSubCategory)

				categories.GET("/:id/sub-categories/:subCategoryId/items", controllers.ShowItemBySubCategory)
				categories.POST("/:id/sub-categories/:subCategoryId/items", controllers.CreateItem)
				categories.GET("/:id/sub-categories/:subCategoryId/items/:itemId", controllers.GetItemByID)
				categories.PUT("/:id/sub-categories/:subCategoryId/items/:itemId", controllers.UpdateItem)
				categories.DELETE("/:id/sub-categories/:subCategoryId/items/:itemId", controllers.DeleteItem)
			}

			tickets := admin.Group("/tickets")
			{
				tickets.GET("", controllers.ShowTickets)
				tickets.POST("", controllers.CreateTicket)
				tickets.GET("/:ticket_number", controllers.ShowTicketDetail)
				tickets.PUT("/:ticket_number", controllers.UpdateTicket)
				tickets.PUT("/:ticket_number/status", controllers.UpdateTicketStatus)
				tickets.DELETE("/:ticket_number", controllers.DeleteTicket)
				tickets.GET("/actives", controllers.ShowActiveTickets)
			}

			comments := admin.Group("/comments")
			{
				// Karena TicketID dikirim via JSON body (sesuai struct), kita cukup gunakan POST /api/comments
				comments.POST("", controllers.CreateCommentAndCloseTicket)
			}
		}
	}
	return router
}
