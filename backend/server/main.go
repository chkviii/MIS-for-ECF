package main

import (
	"log"
	"net/http"
	"os"

	"erp-backend/internal/config"
	"erp-backend/internal/handlers"
	"erp-backend/internal/middleware"
	"erp-backend/internal/repo"
	"erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	//config log path
	logFile, err := os.OpenFile("./server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)

	// 加载配置
	cfg := config.Load()
	log.Printf("Loaded configuration: Port=%s, DB_Path=%s", cfg.Port, cfg.DB_Path)

	// 初始化数据库
	if err := repo.InitDatabase(cfg.DB_Path); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 获取 GORM DB 实例
	db := repo.GetDB()

	// 初始化 Repositories (使用 GORM)
	userRepo := repo.NewUserRepository(db)
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
	volunteerProjectRepo := repo.NewVolunteerProjectRepository(db)
	employeeProjectRepo := repo.NewEmployeeProjectRepository(db)
	fundProjectRepo := repo.NewFundProjectRepository(db)
	donationInventoryRepo := repo.NewDonationInventoryRepository(db)
	deliveryInventoryRepo := repo.NewDeliveryInventoryRepository(db)
	scheduleRepo := repo.NewScheduleRepository(db)

	// 初始化 Services
	// AuthService 依赖多个 Repository (userRepo, employeeRepo, volunteerRepo, donorRepo)
	authService := services.NewAuthService(userRepo, employeeRepo, volunteerRepo, donorRepo)

	// 其他 Services 现在都依赖各自的 Repository
	projectService := services.NewProjectService(projectRepo)
	donorService := services.NewDonorService(donorRepo)
	donationService := services.NewDonationService(donationRepo)
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
	volunteerProjectService := services.NewVolunteerProjectService(volunteerProjectRepo)
	employeeProjectService := services.NewEmployeeProjectService(employeeProjectRepo)
	fundProjectService := services.NewFundProjectService(fundProjectRepo)
	donationInventoryService := services.NewDonationInventoryService(donationInventoryRepo)
	deliveryInventoryService := services.NewDeliveryInventoryService(deliveryInventoryRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)

	// 初始化 Handlers
	authHandler := handlers.NewAuthHandler(authService)
	erpHandler := handlers.NewERPHandler(
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
		deliveryInventoryService,
		scheduleService,
	)

	// Repositories 现在被使用了（通过 authService）
	// 其他 Repositories 可以在重构其他 Services 时使用
	_ = projectRepo
	_ = donorRepo
	_ = donationRepo
	_ = locationRepo
	_ = fundRepo
	_ = expenseRepo
	_ = transactionRepo
	_ = purchaseRepo
	_ = payrollRepo
	_ = inventoryRepo
	_ = giftTypeRepo
	_ = giftRepo
	_ = inventoryTransactionRepo
	_ = deliveryRepo
	_ = volunteerProjectRepo
	_ = employeeProjectRepo
	_ = fundProjectRepo
	_ = donationInventoryRepo
	_ = deliveryInventoryRepo
	_ = scheduleRepo

	// 初始化 Gin 路由
	r := gin.Default()

	// 加载模板（使用配置中的路径）
	r.LoadHTMLGlob(cfg.Html_Path + "/*")

	// 静态文件（使用配置中的路径）
	r.Static("/static", cfg.Static_Path)

	// 公开路由
	public := r.Group("/")
	{
		// 首页 - Welcome 页面
		public.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "welcome.html", gin.H{
				"title": "Elderly Care Foundation",
			})
		})

		// 认证页面
		public.GET("/login", authHandler.ShowLoginPage)
		public.GET("/register", authHandler.ShowRegisterPage)

		// ERP 管理页面（公开访问，前端 JS 会验证 token）
		public.GET("/erp-management", func(c *gin.Context) {
			c.HTML(http.StatusOK, "erp-management.html", gin.H{
				"title": "ERP Management System",
			})
		})

		// 兼容旧路径
		public.GET("/erp", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/erp-management")
		})

		// 认证 API
		public.POST("/api/auth/register", authHandler.Register)
		public.POST("/api/auth/login", authHandler.Login)

		// 兼容旧版 API 路径
		public.POST("/api/v1/auth/register", authHandler.Register)
		public.POST("/api/v1/auth/login", authHandler.Login)
	}

	// 需要认证的路由
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddlewareGin())
	{
		// 需要认证的路由已移除（页面路由改为公开）
	}

	// 需要认证的 API 路由
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddlewareGin())
	{
		// 认证相关
		api.POST("/auth/logout", authHandler.Logout)

		// 项目管理
		api.POST("/projects", erpHandler.CreateProject)
		api.GET("/projects", erpHandler.GetAllProjects)
		api.GET("/projects/search", erpHandler.FilterProjects)
		api.PUT("/projects/:id", erpHandler.UpdateProject)
		api.DELETE("/projects/:id", erpHandler.DeleteProject)

		// 捐赠者管理
		api.POST("/donors", erpHandler.CreateDonor)
		api.GET("/donors", erpHandler.GetAllDonors)
		api.GET("/donors/search", erpHandler.FilterDonors)
		api.PUT("/donors/:id", erpHandler.UpdateDonor)
		api.DELETE("/donors/:id", erpHandler.DeleteDonor)

		// 捐赠管理
		api.POST("/donations", erpHandler.CreateDonation)
		api.GET("/donations", erpHandler.GetAllDonations)
		api.GET("/donations/search", erpHandler.FilterDonations)
		api.PUT("/donations/:id", erpHandler.UpdateDonation)
		api.DELETE("/donations/:id", erpHandler.DeleteDonation)

		// 志愿者管理
		api.POST("/volunteers", erpHandler.CreateVolunteer)
		api.GET("/volunteers", erpHandler.GetAllVolunteers)
		api.GET("/volunteers/search", erpHandler.FilterVolunteers)
		api.PUT("/volunteers/:id", erpHandler.UpdateVolunteer)
		api.DELETE("/volunteers/:id", erpHandler.DeleteVolunteer)

		// 员工管理
		api.POST("/employees", erpHandler.CreateEmployee)
		api.GET("/employees", erpHandler.GetAllEmployees)
		api.GET("/employees/search", erpHandler.FilterEmployees)
		api.PUT("/employees/:id", erpHandler.UpdateEmployee)
		api.DELETE("/employees/:id", erpHandler.DeleteEmployee)

		// 地点管理
		api.POST("/locations", erpHandler.CreateLocation)
		api.GET("/locations", erpHandler.GetAllLocations)
		api.GET("/locations/search", erpHandler.FilterLocations)
		api.PUT("/locations/:id", erpHandler.UpdateLocation)
		api.DELETE("/locations/:id", erpHandler.DeleteLocation)

		// 基金管理
		api.POST("/funds", erpHandler.CreateFund)
		api.GET("/funds", erpHandler.GetAllFunds)
		api.GET("/funds/search", erpHandler.FilterFunds)
		api.PUT("/funds/:id", erpHandler.UpdateFund)
		api.DELETE("/funds/:id", erpHandler.DeleteFund)

		// 支出管理
		api.POST("/expenses", erpHandler.CreateExpense)
		api.GET("/expenses", erpHandler.GetAllExpenses)
		api.GET("/expenses/search", erpHandler.FilterExpenses)
		api.PUT("/expenses/:id", erpHandler.UpdateExpense)
		api.DELETE("/expenses/:id", erpHandler.DeleteExpense)

		// 交易管理
		api.POST("/transactions", erpHandler.CreateTransaction)
		api.GET("/transactions", erpHandler.GetAllTransactions)
		api.GET("/transactions/search", erpHandler.FilterTransactions)
		api.PUT("/transactions/:id", erpHandler.UpdateTransaction)
		api.DELETE("/transactions/:id", erpHandler.DeleteTransaction)

		// 采购管理
		api.POST("/purchases", erpHandler.CreatePurchase)
		api.GET("/purchases", erpHandler.GetAllPurchases)
		api.GET("/purchases/search", erpHandler.FilterPurchases)
		api.PUT("/purchases/:id", erpHandler.UpdatePurchase)
		api.DELETE("/purchases/:id", erpHandler.DeletePurchase)

		// 薪资管理
		api.POST("/payrolls", erpHandler.CreatePayroll)
		api.GET("/payrolls", erpHandler.GetAllPayrolls)
		api.GET("/payrolls/search", erpHandler.FilterPayrolls)
		api.PUT("/payrolls/:id", erpHandler.UpdatePayroll)
		api.DELETE("/payrolls/:id", erpHandler.DeletePayroll)

		// 库存管理
		api.POST("/inventories", erpHandler.CreateInventory)
		api.GET("/inventories", erpHandler.GetAllInventories)
		api.GET("/inventories/search", erpHandler.FilterInventories)
		api.PUT("/inventories/:id", erpHandler.UpdateInventory)
		api.DELETE("/inventories/:id", erpHandler.DeleteInventory)

		// 礼品类型管理
		api.POST("/gift-types", erpHandler.CreateGiftType)
		api.GET("/gift-types", erpHandler.GetAllGiftTypes)
		api.GET("/gift-types/search", erpHandler.FilterGiftTypes)
		api.PUT("/gift-types/:id", erpHandler.UpdateGiftType)
		api.DELETE("/gift-types/:id", erpHandler.DeleteGiftType)

		// 礼品管理
		api.POST("/gifts", erpHandler.CreateGift)
		api.GET("/gifts", erpHandler.GetAllGifts)
		api.GET("/gifts/search", erpHandler.FilterGifts)
		api.PUT("/gifts/:id", erpHandler.UpdateGift)
		api.DELETE("/gifts/:id", erpHandler.DeleteGift)

		// 库存交易管理
		api.POST("/inventory-transactions", erpHandler.CreateInventoryTransaction)
		api.GET("/inventory-transactions", erpHandler.GetAllInventoryTransactions)
		api.GET("/inventory-transactions/search", erpHandler.FilterInventoryTransactions)
		api.PUT("/inventory-transactions/:id", erpHandler.UpdateInventoryTransaction)
		api.DELETE("/inventory-transactions/:id", erpHandler.DeleteInventoryTransaction)

		// 配送管理
		api.POST("/deliveries", erpHandler.CreateDelivery)
		api.GET("/deliveries", erpHandler.GetAllDeliveries)
		api.GET("/deliveries/search", erpHandler.FilterDeliveries)
		api.PUT("/deliveries/:id", erpHandler.UpdateDelivery)
		api.DELETE("/deliveries/:id", erpHandler.DeleteDelivery)

		// 志愿者-项目关联
		api.POST("/volunteer-projects", erpHandler.CreateVolunteerProject)
		api.GET("/volunteer-projects", erpHandler.GetAllVolunteerProjects)
		api.GET("/volunteer-projects/search", erpHandler.FilterVolunteerProjects)
		api.PUT("/volunteer-projects/:id", erpHandler.UpdateVolunteerProject)
		api.DELETE("/volunteer-projects/:id", erpHandler.DeleteVolunteerProject)

		// 员工-项目关联
		api.POST("/employee-projects", erpHandler.CreateEmployeeProject)
		api.GET("/employee-projects", erpHandler.GetAllEmployeeProjects)
		api.GET("/employee-projects/search", erpHandler.FilterEmployeeProjects)
		api.PUT("/employee-projects/:id", erpHandler.UpdateEmployeeProject)
		api.DELETE("/employee-projects/:id", erpHandler.DeleteEmployeeProject)

		// 基金-项目关联
		api.POST("/fund-projects", erpHandler.CreateFundProject)
		api.GET("/fund-projects", erpHandler.GetAllFundProjects)
		api.GET("/fund-projects/search", erpHandler.FilterFundProjects)
		api.PUT("/fund-projects/:id", erpHandler.UpdateFundProject)
		api.DELETE("/fund-projects/:id", erpHandler.DeleteFundProject)

		// 捐赠-库存关联
		api.POST("/donation-inventories", erpHandler.CreateDonationInventory)
		api.GET("/donation-inventories", erpHandler.GetAllDonationInventories)
		api.GET("/donation-inventories/search", erpHandler.FilterDonationInventories)
		api.PUT("/donation-inventories/:id", erpHandler.UpdateDonationInventory)
		api.DELETE("/donation-inventories/:id", erpHandler.DeleteDonationInventory)

		// 交付-库存关联
		api.POST("/delivery-inventories", erpHandler.CreateDeliveryInventory)
		api.GET("/delivery-inventories", erpHandler.GetAllDeliveryInventories)
		api.GET("/delivery-inventories/search", erpHandler.FilterDeliveryInventories)
		api.PUT("/delivery-inventories/:id", erpHandler.UpdateDeliveryInventory)
		api.DELETE("/delivery-inventories/:id", erpHandler.DeleteDeliveryInventory)

		// 日程管理
		api.POST("/schedules", erpHandler.CreateSchedule)
		api.GET("/schedules", erpHandler.GetAllSchedules)
		api.GET("/schedules/search", erpHandler.FilterSchedules)
		api.PUT("/schedules/:id", erpHandler.UpdateSchedule)
		api.DELETE("/schedules/:id", erpHandler.DeleteSchedule)
	}

	// 启动服务器（使用配置中的端口）
	log.Printf("Server starting on http://localhost%s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
