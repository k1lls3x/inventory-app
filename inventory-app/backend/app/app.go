package app

import (
	"context"
	"fmt"
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/dashboard"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	db.Init() // ✅ инициализация подключения к PostgreSQL
}

// Greet test method
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}


// GetDashboard fetches total_stock, item_count, monthly_orders
func (a *App) GetDashboard() (*dashboard.DashboardData, error) {
	return dashboard.LoadDashboard()
}
