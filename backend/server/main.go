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

	chartRepo := repo.NewChartRepository(db)

	// 初始化 Services
	// AuthService 依赖多个 Repository (userRepo, employeeRepo, volunteerRepo, donorRepo)
	authService := services.NewAuthService(userRepo, employeeRepo, volunteerRepo, donorRepo)
	chartService := services.NewChartService(chartRepo)
	donService := services.NewDonService(donorRepo, projectRepo, employeeProjectRepo)

	// 其他 Services 现在都依赖各自的 Repository
	userService := services.NewUserService(userRepo)
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
	chartHandler := handlers.NewChartHandler(chartService)
	donHandler := handlers.NewDonHandler(donService)

	erpHandler := handlers.NewERPHandler(
		userService,
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
	// _ = projectRepo
	// _ = donorRepo
	// _ = donationRepo
	// _ = locationRepo
	// _ = fundRepo
	// _ = expenseRepo
	// _ = transactionRepo
	// _ = purchaseRepo
	// _ = payrollRepo
	// _ = inventoryRepo
	// _ = giftTypeRepo
	// _ = giftRepo
	// _ = inventoryTransactionRepo
	// _ = deliveryRepo
	// _ = volunteerProjectRepo
	// _ = employeeProjectRepo
	// _ = fundProjectRepo
	// _ = donationInventoryRepo
	// _ = deliveryInventoryRepo
	// _ = scheduleRepo

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

		public.GET("/sysadmin", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user-management.html", gin.H{
				"title": "System Admin",
			})
		})

		// Donor/Volunteer/Employee portals (simple pages)
		public.GET("/donor", func(c *gin.Context) {
			c.HTML(http.StatusOK, "donor.html", gin.H{"title": "Donor Portal"})
		})

		public.GET("/volunteer", func(c *gin.Context) {
			c.HTML(http.StatusOK, "volunteer.html", gin.H{"title": "Volunteer Portal"})
		})

		public.GET("/employee", func(c *gin.Context) {
			c.HTML(http.StatusOK, "employee-dashboard.html", gin.H{"title": "Employee Portal"})
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
		authenticated.POST("/auth/logout", authHandler.Logout)
	}

	admin_api := r.Group("/api/v1/dbms/users")
	admin_api.Use(middleware.AuthMiddlewareGin())
	admin_api.Use(middleware.AuthVarifyUserType("employee"))
	{
		// user management
		admin_api.POST("/", erpHandler.CreateUser)
		admin_api.GET("/", erpHandler.GetAllUsers)
		admin_api.GET("/search", erpHandler.FilterUsers)
		admin_api.PUT("/:id", erpHandler.UpdateUser)
		admin_api.DELETE("/:id", erpHandler.DeleteUser)
	}

	// Financial Charts API for employee dashboard
	finchart_api := r.Group("/api/v1/fin/charts")
	finchart_api.Use(middleware.AuthMiddlewareGin())
	finchart_api.Use(middleware.AuthVarifyUserType("employee"))
	{
		finchart_api.GET("/line/fund", chartHandler.FundAllocations)
		finchart_api.GET("/pie/fund", chartHandler.FundAllocationsByProject)
		finchart_api.GET("/line/expenses", chartHandler.Expenses)
		finchart_api.GET("/pie/expenses", chartHandler.ExpensesByProject)
		finchart_api.GET("/line/donations", chartHandler.Donations)
		finchart_api.GET("/pie/donations", chartHandler.DonationsByProject)

	}

	//Donation Charts API for donor dashboard
	don_api := r.Group("/api/v1/donor")
	don_api.Use(middleware.AuthMiddlewareGin())
	don_api.Use(middleware.AuthVarifyUserType("donor"))
	{
		don_api.GET("/charts/line/donations", chartHandler.DonorDonations)
		don_api.GET("/charts/pie/donations", chartHandler.DonorDonationsByProject)
		don_api.GET("/projects", donHandler.GetProjectsByDonor)
		don_api.GET("/donations", donHandler.GetDonationDetails)
	}

	// dbms API for employee
	dbms_api := r.Group("/api/v1/dbms")
	dbms_api.Use(middleware.AuthMiddlewareGin())
	dbms_api.Use(middleware.AuthVarifyUserType("employee"))

	{

		// 项目管理
		dbms_api.POST("/projects", erpHandler.CreateProject)
		dbms_api.GET("/projects", erpHandler.GetAllProjects)
		dbms_api.GET("/projects/search", erpHandler.FilterProjects)
		dbms_api.PUT("/projects/:id", erpHandler.UpdateProject)
		dbms_api.DELETE("/projects/:id", erpHandler.DeleteProject)

		// 捐赠者管理
		dbms_api.POST("/donors", erpHandler.CreateDonor)
		dbms_api.GET("/donors", erpHandler.GetAllDonors)
		dbms_api.GET("/donors/search", erpHandler.FilterDonors)
		dbms_api.PUT("/donors/:id", erpHandler.UpdateDonor)
		dbms_api.DELETE("/donors/:id", erpHandler.DeleteDonor)

		// 捐赠管理
		dbms_api.POST("/donations", erpHandler.CreateDonation)
		dbms_api.GET("/donations", erpHandler.GetAllDonations)
		dbms_api.GET("/donations/search", erpHandler.FilterDonations)
		dbms_api.PUT("/donations/:id", erpHandler.UpdateDonation)
		dbms_api.DELETE("/donations/:id", erpHandler.DeleteDonation)

		// 志愿者管理
		dbms_api.POST("/volunteers", erpHandler.CreateVolunteer)
		dbms_api.GET("/volunteers", erpHandler.GetAllVolunteers)
		dbms_api.GET("/volunteers/search", erpHandler.FilterVolunteers)
		dbms_api.PUT("/volunteers/:id", erpHandler.UpdateVolunteer)
		dbms_api.DELETE("/volunteers/:id", erpHandler.DeleteVolunteer)

		// 员工管理
		dbms_api.POST("/employees", erpHandler.CreateEmployee)
		dbms_api.GET("/employees", erpHandler.GetAllEmployees)
		dbms_api.GET("/employees/search", erpHandler.FilterEmployees)
		dbms_api.PUT("/employees/:id", erpHandler.UpdateEmployee)
		dbms_api.DELETE("/employees/:id", erpHandler.DeleteEmployee)

		// 地点管理
		dbms_api.POST("/locations", erpHandler.CreateLocation)
		dbms_api.GET("/locations", erpHandler.GetAllLocations)
		dbms_api.GET("/locations/search", erpHandler.FilterLocations)
		dbms_api.PUT("/locations/:id", erpHandler.UpdateLocation)
		dbms_api.DELETE("/locations/:id", erpHandler.DeleteLocation)

		// 基金管理
		dbms_api.POST("/funds", erpHandler.CreateFund)
		dbms_api.GET("/funds", erpHandler.GetAllFunds)
		dbms_api.GET("/funds/search", erpHandler.FilterFunds)
		dbms_api.PUT("/funds/:id", erpHandler.UpdateFund)
		dbms_api.DELETE("/funds/:id", erpHandler.DeleteFund)

		// 支出管理
		dbms_api.POST("/expenses", erpHandler.CreateExpense)
		dbms_api.GET("/expenses", erpHandler.GetAllExpenses)
		dbms_api.GET("/expenses/search", erpHandler.FilterExpenses)
		dbms_api.PUT("/expenses/:id", erpHandler.UpdateExpense)
		dbms_api.DELETE("/expenses/:id", erpHandler.DeleteExpense)

		// 交易管理
		dbms_api.POST("/transactions", erpHandler.CreateTransaction)
		dbms_api.GET("/transactions", erpHandler.GetAllTransactions)
		dbms_api.GET("/transactions/search", erpHandler.FilterTransactions)
		dbms_api.PUT("/transactions/:id", erpHandler.UpdateTransaction)
		dbms_api.DELETE("/transactions/:id", erpHandler.DeleteTransaction)

		// 采购管理
		dbms_api.POST("/purchases", erpHandler.CreatePurchase)
		dbms_api.GET("/purchases", erpHandler.GetAllPurchases)
		dbms_api.GET("/purchases/search", erpHandler.FilterPurchases)
		dbms_api.PUT("/purchases/:id", erpHandler.UpdatePurchase)
		dbms_api.DELETE("/purchases/:id", erpHandler.DeletePurchase)

		// 薪资管理
		dbms_api.POST("/payrolls", erpHandler.CreatePayroll)
		dbms_api.GET("/payrolls", erpHandler.GetAllPayrolls)
		dbms_api.GET("/payrolls/search", erpHandler.FilterPayrolls)
		dbms_api.PUT("/payrolls/:id", erpHandler.UpdatePayroll)
		dbms_api.DELETE("/payrolls/:id", erpHandler.DeletePayroll)

		// 库存管理
		dbms_api.POST("/inventory", erpHandler.CreateInventory)
		dbms_api.GET("/inventory", erpHandler.GetAllInventories)
		dbms_api.GET("/inventory/search", erpHandler.FilterInventories)
		dbms_api.PUT("/inventory/:id", erpHandler.UpdateInventory)
		dbms_api.DELETE("/inventory/:id", erpHandler.DeleteInventory)

		// 礼品类型管理
		dbms_api.POST("/gift-types", erpHandler.CreateGiftType)
		dbms_api.GET("/gift-types", erpHandler.GetAllGiftTypes)
		dbms_api.GET("/gift-types/search", erpHandler.FilterGiftTypes)
		dbms_api.PUT("/gift-types/:id", erpHandler.UpdateGiftType)
		dbms_api.DELETE("/gift-types/:id", erpHandler.DeleteGiftType)

		// 礼品管理
		dbms_api.POST("/gifts", erpHandler.CreateGift)
		dbms_api.GET("/gifts", erpHandler.GetAllGifts)
		dbms_api.GET("/gifts/search", erpHandler.FilterGifts)
		dbms_api.PUT("/gifts/:id", erpHandler.UpdateGift)
		dbms_api.DELETE("/gifts/:id", erpHandler.DeleteGift)

		// 库存交易管理
		dbms_api.POST("/inventory-transactions", erpHandler.CreateInventoryTransaction)
		dbms_api.GET("/inventory-transactions", erpHandler.GetAllInventoryTransactions)
		dbms_api.GET("/inventory-transactions/search", erpHandler.FilterInventoryTransactions)
		dbms_api.PUT("/inventory-transactions/:id", erpHandler.UpdateInventoryTransaction)
		dbms_api.DELETE("/inventory-transactions/:id", erpHandler.DeleteInventoryTransaction)

		// 配送管理
		dbms_api.POST("/deliveries", erpHandler.CreateDelivery)
		dbms_api.GET("/deliveries", erpHandler.GetAllDeliveries)
		dbms_api.GET("/deliveries/search", erpHandler.FilterDeliveries)
		dbms_api.PUT("/deliveries/:id", erpHandler.UpdateDelivery)
		dbms_api.DELETE("/deliveries/:id", erpHandler.DeleteDelivery)

		// 志愿者-项目关联
		dbms_api.POST("/volunteer-projects", erpHandler.CreateVolunteerProject)
		dbms_api.GET("/volunteer-projects", erpHandler.GetAllVolunteerProjects)
		dbms_api.GET("/volunteer-projects/search", erpHandler.FilterVolunteerProjects)
		dbms_api.PUT("/volunteer-projects/:id", erpHandler.UpdateVolunteerProject)
		dbms_api.DELETE("/volunteer-projects/:id", erpHandler.DeleteVolunteerProject)

		// 员工-项目关联
		dbms_api.POST("/employee-projects", erpHandler.CreateEmployeeProject)
		dbms_api.GET("/employee-projects", erpHandler.GetAllEmployeeProjects)
		dbms_api.GET("/employee-projects/search", erpHandler.FilterEmployeeProjects)
		dbms_api.PUT("/employee-projects/:id", erpHandler.UpdateEmployeeProject)
		dbms_api.DELETE("/employee-projects/:id", erpHandler.DeleteEmployeeProject)

		// 基金-项目关联
		dbms_api.POST("/fund-projects", erpHandler.CreateFundProject)
		dbms_api.GET("/fund-projects", erpHandler.GetAllFundProjects)
		dbms_api.GET("/fund-projects/search", erpHandler.FilterFundProjects)
		dbms_api.PUT("/fund-projects/:id", erpHandler.UpdateFundProject)
		dbms_api.DELETE("/fund-projects/:id", erpHandler.DeleteFundProject)

		// 捐赠-库存关联
		dbms_api.POST("/donation-inventory", erpHandler.CreateDonationInventory)
		dbms_api.GET("/donation-inventory", erpHandler.GetAllDonationInventories)
		dbms_api.GET("/donation-inventory/search", erpHandler.FilterDonationInventories)
		dbms_api.PUT("/donation-inventory/:id", erpHandler.UpdateDonationInventory)
		dbms_api.DELETE("/donation-inventory/:id", erpHandler.DeleteDonationInventory)

		// 交付-库存关联
		dbms_api.POST("/delivery-inventory", erpHandler.CreateDeliveryInventory)
		dbms_api.GET("/delivery-inventory", erpHandler.GetAllDeliveryInventories)
		dbms_api.GET("/delivery-inventory/search", erpHandler.FilterDeliveryInventories)
		dbms_api.PUT("/delivery-inventory/:id", erpHandler.UpdateDeliveryInventory)
		dbms_api.DELETE("/delivery-inventory/:id", erpHandler.DeleteDeliveryInventory)

		// 日程管理
		dbms_api.POST("/schedules", erpHandler.CreateSchedule)
		dbms_api.GET("/schedules", erpHandler.GetAllSchedules)
		dbms_api.GET("/schedules/search", erpHandler.FilterSchedules)
		dbms_api.PUT("/schedules/:id", erpHandler.UpdateSchedule)
		dbms_api.DELETE("/schedules/:id", erpHandler.DeleteSchedule)
	}

	// 启动服务器（使用配置中的端口）
	log.Printf("Server starting on http://localhost%s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
