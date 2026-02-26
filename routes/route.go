package routes

import (
	"backend-golang-api/controllers"
	"backend-golang-api/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/refresh", controllers.RefreshToken)

		admin := api.Group("", middlewares.AuthMiddleware())
		{
			users := admin.Group("/users")
			{
				users.GET("", controllers.ShowUsers)
				users.POST("", controllers.CreateUser)
				users.GET("/:id", controllers.FindUserById)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
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
				tickets.DELETE("/:ticket_number", controllers.DeleteTicket)
				tickets.GET("/actives", controllers.ShowActiveTickets)
			}
		}
	}
	return router
}
