package app

import (
	"context"
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/handler/dashboard"
	"inventory-app/backend/logs"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

)

// App struct
type App struct {
	ctx context.Context
}
func (a *App) GetTopItems() ([]model.ItemWithStock, error) {
	return dashboard.LoadTopItems()
}
// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}
func (a *App) GetTurnoverByWarehouse() ([]model.ItemTurnoverByWarehouse, error) {
	return dashboard.LoadTurnoverByWarehouse()
}
// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logs.Init()
	db.Init()

}

func (a *App) GetDashboard() (*dashboard.DashboardData, error) {
	return dashboard.LoadDashboard()
}

func (a *App) GetStock() ([]model.Stock, error) {
	return repository.GetStocks()
}
