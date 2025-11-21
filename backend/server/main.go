package main

import (
	"fmt"
	"log"
	"path/filepath"

	"mypage-backend/internal/config"
	"mypage-backend/internal/handler"
	"mypage-backend/internal/repo"
	"mypage-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Println("GO: Loaded configuration:", *cfg)

	// Initialize database
	dbPath := cfg.DB_Path
	if err := repo.InitDatabase(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.CloseDatabase()
	fmt.Println("GO: Database initialized successfully")

	// Get database instance
	db := repo.GetDB()

	// Initialize repositories
	projectRepo := repo.NewProjectRepository(db)
	donorRepo := repo.NewDonorRepository(db)
	donationRepo := repo.NewDonationRepository(db)
	volunteerRepo := repo.NewVolunteerRepository(db)
	employeeRepo := repo.NewEmployeeRepository(db)
	locationRepo := repo.NewLocationRepository(db)
	fundRepo := repo.NewFundRepository(db)
	expenseRepo := repo.NewExpenseRepository(db)
	transactionRepo := repo.NewTransactionRepository(db)
	purchaseRepo := repo.NewPurchaseRepository(db)
	payrollRepo := repo.NewPayrollRepository(db)
	inventoryRepo := repo.NewInventoryRepository(db)
	giftTypeRepo := repo.NewGiftTypeRepository(db)
	giftRepo := repo.NewGiftRepository(db)
	inventoryTransactionRepo := repo.NewInventoryTransactionRepository(db)
	deliveryRepo := repo.NewDeliveryRepository(db)
	
	// 关系表 repositories
	volunteerProjectRepo := repo.NewVolunteerProjectRepository(db)
	employeeProjectRepo := repo.NewEmployeeProjectRepository(db)
	fundProjectRepo := repo.NewFundProjectRepository(db)
	donationInventoryRepo := repo.NewDonationInventoryRepository(db)
	scheduleRepo := repo.NewScheduleRepository(db)

	// Initialize services
	projectService := services.NewProjectService(projectRepo, locationRepo)
	donorService := services.NewDonorService(donorRepo)
	donationService := services.NewDonationService(donationRepo, donorRepo)
	volunteerService := services.NewVolunteerService(volunteerRepo)
	employeeService := services.NewEmployeeService(employeeRepo)
	locationService := services.NewLocationService(locationRepo)
	fundService := services.NewFundService(fundRepo)
	expenseService := services.NewExpenseService(expenseRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo)
	payrollService := services.NewPayrollService(payrollRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)
	giftTypeService := services.NewGiftTypeService(giftTypeRepo)
	giftService := services.NewGiftService(giftRepo)
	inventoryTransactionService := services.NewInventoryTransactionService(inventoryTransactionRepo)
	deliveryService := services.NewDeliveryService(deliveryRepo)
	
	// 关系表 services
	volunteerProjectService := services.NewVolunteerProjectService(volunteerProjectRepo)
	employeeProjectService := services.NewEmployeeProjectService(employeeProjectRepo)
	fundProjectService := services.NewFundProjectService(fundProjectRepo)
	donationInventoryService := services.NewDonationInventoryService(donationInventoryRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)

	// Initialize handlers
	erpHandler := handler.NewERPHandler(
		projectService,
		donorService,
		donationService,
		volunteerService,
		employeeService,
		locationService,
		fundService,
		expenseService,
		transactionService,
		purchaseService,
		payrollService,
		inventoryService,
		giftTypeService,
		giftService,
		inventoryTransactionService,
		deliveryService,
		volunteerProjectService,
		employeeProjectService,
		fundProjectService,
		donationInventoryService,
		scheduleService,
	)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Link")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Static file serving
	r.Static("/static", cfg.Static_Path)

	// Public routes
	r.GET("/", handler.HomeHandler())
	r.GET("/erp-management.html", func(c *gin.Context) {
		c.File(filepath.Join(cfg.Html_Path, "erp-management.html"))
	})
	r.GET("/index.html", func(c *gin.Context) {
		c.File(filepath.Join(cfg.Html_Path, "index.html"))
	})
	r.GET("/blog.html", func(c *gin.Context) {
		c.File(filepath.Join(cfg.Html_Path, "blog.html"))
	})
	r.GET("/register.html", func(c *gin.Context) {
		c.File(filepath.Join(cfg.Html_Path, "register.html"))
	})

	// ERP API routes - v1
	v1 := r.Group("/api/v1")
	{
		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", erpHandler.GetAllProjects)
			projects.POST("", erpHandler.CreateProject)
			projects.GET("/:id", erpHandler.GetProject)
			projects.PUT("/:id", erpHandler.UpdateProject)
			projects.DELETE("/:id", erpHandler.DeleteProject)
		}

		// Donors
		donors := v1.Group("/donors")
		{
			donors.GET("", erpHandler.GetAllDonors)
			donors.POST("", erpHandler.CreateDonor)
			donors.GET("/:id", erpHandler.GetDonor)
			donors.PUT("/:id", erpHandler.UpdateDonor)
			donors.DELETE("/:id", erpHandler.DeleteDonor)
		}

		// Donations
		donations := v1.Group("/donations")
		{
			donations.GET("", erpHandler.GetAllDonations)
			donations.POST("", erpHandler.CreateDonation)
			donations.PUT("/:id", erpHandler.UpdateDonation)
			donations.DELETE("/:id", erpHandler.DeleteDonation)
		}

		// Volunteers
		volunteers := v1.Group("/volunteers")
		{
			volunteers.GET("", erpHandler.GetAllVolunteers)
			volunteers.POST("", erpHandler.CreateVolunteer)
			volunteers.PUT("/:id", erpHandler.UpdateVolunteer)
			volunteers.DELETE("/:id", erpHandler.DeleteVolunteer)
		}

		// Employees
		employees := v1.Group("/employees")
		{
			employees.GET("", erpHandler.GetAllEmployees)
			employees.POST("", erpHandler.CreateEmployee)
			employees.PUT("/:id", erpHandler.UpdateEmployee)
			employees.DELETE("/:id", erpHandler.DeleteEmployee)
		}

		// Locations
		locations := v1.Group("/locations")
		{
			locations.GET("", erpHandler.GetAllLocations)
			locations.POST("", erpHandler.CreateLocation)
			locations.PUT("/:id", erpHandler.UpdateLocation)
			locations.DELETE("/:id", erpHandler.DeleteLocation)
		}

		// Funds
		funds := v1.Group("/funds")
		{
			funds.GET("", erpHandler.GetAllFunds)
			funds.POST("", erpHandler.CreateFund)
			funds.GET("/:id", erpHandler.GetFund)
			funds.PUT("/:id", erpHandler.UpdateFund)
			funds.DELETE("/:id", erpHandler.DeleteFund)
		}

		// Expenses
		expenses := v1.Group("/expenses")
		{
			expenses.GET("", erpHandler.GetAllExpenses)
			expenses.POST("", erpHandler.CreateExpense)
			expenses.PUT("/:id", erpHandler.UpdateExpense)
			expenses.DELETE("/:id", erpHandler.DeleteExpense)
		}

		// Transactions
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", erpHandler.GetAllTransactions)
			transactions.POST("", erpHandler.CreateTransaction)
			transactions.PUT("/:id", erpHandler.UpdateTransaction)
			transactions.DELETE("/:id", erpHandler.DeleteTransaction)
		}

		// Purchases
		purchases := v1.Group("/purchases")
		{
			purchases.GET("", erpHandler.GetAllPurchases)
			purchases.POST("", erpHandler.CreatePurchase)
			purchases.PUT("/:id", erpHandler.UpdatePurchase)
			purchases.DELETE("/:id", erpHandler.DeletePurchase)
		}

		// Payrolls
		payrolls := v1.Group("/payrolls")
		{
			payrolls.GET("", erpHandler.GetAllPayrolls)
			payrolls.POST("", erpHandler.CreatePayroll)
			payrolls.PUT("/:id", erpHandler.UpdatePayroll)
			payrolls.DELETE("/:id", erpHandler.DeletePayroll)
		}

		// Inventory
		inventory := v1.Group("/inventory")
		{
			inventory.GET("", erpHandler.GetAllInventories)
			inventory.POST("", erpHandler.CreateInventory)
			inventory.PUT("/:id", erpHandler.UpdateInventory)
			inventory.DELETE("/:id", erpHandler.DeleteInventory)
		}

		// Gift Types
		giftTypes := v1.Group("/gift-types")
		{
			giftTypes.GET("", erpHandler.GetAllGiftTypes)
			giftTypes.POST("", erpHandler.CreateGiftType)
			giftTypes.PUT("/:id", erpHandler.UpdateGiftType)
			giftTypes.DELETE("/:id", erpHandler.DeleteGiftType)
		}

		// Gifts
		gifts := v1.Group("/gifts")
		{
			gifts.GET("", erpHandler.GetAllGifts)
			gifts.POST("", erpHandler.CreateGift)
			gifts.PUT("/:id", erpHandler.UpdateGift)
			gifts.DELETE("/:id", erpHandler.DeleteGift)
		}

		// Inventory Transactions
		inventoryTransactions := v1.Group("/inventory-transactions")
		{
			inventoryTransactions.GET("", erpHandler.GetAllInventoryTransactions)
			inventoryTransactions.POST("", erpHandler.CreateInventoryTransaction)
			inventoryTransactions.PUT("/:id", erpHandler.UpdateInventoryTransaction)
			inventoryTransactions.DELETE("/:id", erpHandler.DeleteInventoryTransaction)
		}

		// Deliveries
		deliveries := v1.Group("/deliveries")
		{
			deliveries.GET("", erpHandler.GetAllDeliveries)
			deliveries.POST("", erpHandler.CreateDelivery)
			deliveries.PUT("/:id", erpHandler.UpdateDelivery)
			deliveries.DELETE("/:id", erpHandler.DeleteDelivery)
			}
		
		// 关系表路由
		// Volunteer-Projects
		volunteerProjects := v1.Group("/volunteer-projects")
		{
			volunteerProjects.GET("", erpHandler.GetAllVolunteerProjects)
			volunteerProjects.POST("", erpHandler.CreateVolunteerProject)
			volunteerProjects.PUT("/:id", erpHandler.UpdateVolunteerProject)
			volunteerProjects.DELETE("/:id", erpHandler.DeleteVolunteerProject)
		}

		// Employee-Projects
		employeeProjects := v1.Group("/employee-projects")
		{
			employeeProjects.GET("", erpHandler.GetAllEmployeeProjects)
			employeeProjects.POST("", erpHandler.CreateEmployeeProject)
			employeeProjects.PUT("/:id", erpHandler.UpdateEmployeeProject)
			employeeProjects.DELETE("/:id", erpHandler.DeleteEmployeeProject)
		}

		// Fund-Projects
		fundProjects := v1.Group("/fund-projects")
		{
			fundProjects.GET("", erpHandler.GetAllFundProjects)
			fundProjects.POST("", erpHandler.CreateFundProject)
			fundProjects.PUT("/:id", erpHandler.UpdateFundProject)
			fundProjects.DELETE("/:id", erpHandler.DeleteFundProject)
		}

		// Donation-Inventory
		donationInventory := v1.Group("/donation-inventory")
		{
			donationInventory.GET("", erpHandler.GetAllDonationInventories)
			donationInventory.POST("", erpHandler.CreateDonationInventory)
			donationInventory.PUT("/:id", erpHandler.UpdateDonationInventory)
			donationInventory.DELETE("/:id", erpHandler.DeleteDonationInventory)
		}

		// Schedules
		schedules := v1.Group("/schedules")
		{
			schedules.GET("", erpHandler.GetAllSchedules)
			schedules.POST("", erpHandler.CreateSchedule)
			schedules.PUT("/:id", erpHandler.UpdateSchedule)
			schedules.DELETE("/:id", erpHandler.DeleteSchedule)
		}
	}

	// Start server
	fmt.Println("GO: Server starting on port", "http://localhost"+cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
