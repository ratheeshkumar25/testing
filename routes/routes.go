package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant/controllers"
	"restaurant/middleware"
)

// Routes sets up the routes for the application.
func Routes() *gin.Engine {
	//creates a new Gin engine instance with default configurations
	r := gin.Default()

	//Define user router
	r.GET("/users", controllers.GetHome)
	//r.POST("l/login",controllers.Login)
	r.POST("/users/login", controllers.PostLogin)
	r.POST("/users/login/verify", controllers.SignupVerify)
	r.POST("/logout", controllers.UserLogout)

	//Define the Admin Routes
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/admin/logout", controllers.AdminLogout)

	//admin middleware authentication to add edit and delete
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/menuList", controllers.GetMenuList)
		admin.POST("/menu/add", controllers.CreateMenu)
		admin.PUT("/menu/:id", controllers.UpdateMenu)
		admin.DELETE("menu/:id", controllers.DeleteMenu)

		//table control
		admin.GET("/table", controllers.GetTables)
		admin.GET("/table/:id", controllers.GetTable)
		admin.POST("table/add", controllers.CreateTable)
		admin.PUT("table/:id", controllers.UpdateTable)
		admin.DELETE("table/:id", controllers.RemoveTable)

		//staff control
		admin.GET("/staff", controllers.GetStaff)
		admin.GET("/staff/:id", controllers.GetStaffByIDs)
		admin.POST("/staff/add", controllers.AddStaff)
		admin.PUT("/staff/:id", controllers.UpdateStaff)
		admin.DELETE("/staff/:id", controllers.RemoveStaff)

		//order and invoice controller
		admin.GET("invoice", controllers.GetInvoice)

		//Sales Report
		admin.GET("/totalorder", controllers.TotalOrder)
		admin.GET("/sales", controllers.TotalSales)
		admin.GET("/employeeperformance", controllers.EmployeePerfomance)
		admin.GET("/revenue", controllers.RevenueReport)
		admin.GET("/mostorderitem", controllers.MostOrderedItems)
		admin.GET("/invoices/:id/pdf", controllers.GetPDFInvoice)
	}

	//Users middleware authentication view menulist , specified menu
	users := r.Group("/users")
	users.Use(middleware.UserauthMiddleware())
	{

		users.GET("menu/:id", controllers.GetMenu)
		users.GET("/menulist", controllers.GetMenuList)
		users.GET("/table", controllers.GetTables)
		// users.GET("searchtable/:id", controllers.GetTable)
		users.GET("/searchreservation", controllers.SearchAvailableTables)
		users.POST("/reservation", controllers.CreateReservartion)
		users.PUT("/movereservation/:id", controllers.UpdateReservation)
		users.GET("/cancelreservation/:id", controllers.CancelReservation)
		users.POST("/placeorder/invoice", controllers.PlaceOrder)
		users.POST("/payinvoice/:id", controllers.PayInvoice)
		users.PUT("/updateorder/:id", controllers.UpdatePlaceOrder)
		users.GET("cancelorder/:id", controllers.CancelOrder)
		users.POST("/rating", controllers.Rating)
		users.GET("rating", controllers.ViewReview)

	}
	r.GET("/online/pay/", controllers.MakePayment)
	r.GET("/payment/success", controllers.SuccessPage)
	r.GET("/failed", controllers.FailurePage)

	return r
}
