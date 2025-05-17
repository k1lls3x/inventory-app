package app

import (
	"context"
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/handler/dashboard"
	"inventory-app/backend/internal/handler/export"
	"inventory-app/backend/logs"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"
	"log"
	"os"
	"inventory-app/backend/auth"
	"errors"
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
	  // === ЛОГИРОВАНИЕ В ФАЙЛ ===
    logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        // Если не получилось открыть файл, логируем в стандартный вывод
        log.Println("Ошибка открытия файла логов:", err)
    } else {
        log.SetOutput(logFile)
        log.Println("=== Приложение запущено ===")
    }

    a.ctx = ctx

    log.Println("Startup: инициализация логирования")
  	logs.Init()

    log.Println("Startup: инициализация базы данных")
    db.Init()

    log.Println("Startup: завершён")

}

func (a *App) GetDashboard() (*dashboard.DashboardData, error) {
	return dashboard.LoadDashboard()
}

//-----------------------STOCKS---------------------------------\\

func (a *App) FindStockByWarehouse(warehouseID int) ([]model.ItemWithStock, error){
	return repository.FindStockByWarehouse(warehouseID)
}

func (a *App) ChangeStock(itemID, warehouseID, newQuantity int) error{
	return repository.ChangeStock(itemID, warehouseID, newQuantity)
}

func (a *App) RemoveStock(stockID int) error {
	return repository.RemoveStock(stockID)
}

func (a *App) GetStockDetails() ([]model.ItemWithStock, error){
	return repository.GetStockDetails()
}

func (a *App) GetStocks() ([]model.Stock, error) {
	return repository.GetStocks()
}

func (a *App) AddStock(itemID, quantity, warehouseID int) error {
	return repository.AddStock( itemID, quantity, warehouseID)
}

func (a *App)ExportStockToExcel() (string, error)  {
	return Export.ExportStockToExcel()
}

func (a *App) GetWeeklyStockTrend() ([]model.DailyStock, error) {
	return repository.GetWeeklyStockTrend()
}

//-----------------------Items---------------------------------\\

func (a *App) GetItems() ([]model.Item, error) {
	return repository.GetItems()
}

func (a *App) GetAllItems() ([]model.ItemBrief, error) {
	return repository.GetItemBriefList()
}

func (a *App)  UpdateItem(item model.Item) error {
	return repository.UpdateItem(item)
}

func (a *App) RemoveItem(sku string) error {
	return repository.RemoveItem(sku)
}

func (a *App) AddItem(item model.Item) error {
	return repository.AddItem(item)
}

func (a *App) FindItems(filter model.ItemFilter) ([]model.Item, error) {
	return repository.FindItems(filter)
}

//-----------------------Warehouse---------------------------------\\
func (a *App) GetWarehouses() ([]model.Warehouse, error) {
	return repository.GetWarehouses()
}
//-----------------------Indbound---------------------------------\\

func (a *App) GetInboundDetails()([]model.InboundDetails,error){
	return repository.GetInboundDetails()
}

func (a *App) GetInboundDetailsByDate(date string)([]model.InboundDetails,error){
	return repository.GetInboundDetailsByDate(date)
}

func (a *App) AddInbound(inb model.Inbound) error{
	return repository.AddInbound(inb);
}
func (a *App) DeleteInbound(inboundId int) error{
	return repository.DeleteInbound(inboundId)
}

func (a *App)  EditInbound(inb model.Inbound) error{
	return repository.EditInbound(inb)
}
//-----------------------Supplier---------------------------------\\
func (a *App)  GetSuppliers() ([]model.Supplier, error){
	return repository.GetSuppliers()
}
//-----------------------Auth---------------------------------\\

func (a *App) RegisterUser(username, password, fullName, role string) error {
	user := &model.User{
			Username: username,
			Role:     role,
			FullName: fullName,
	}
	return auth.Register(user, password)
}

func (a *App) LoginUser(username, password string) (*model.User, error) {
	user, ok := auth.Authorize(username, password)
	if !ok {
			return nil, errors.New("Неверный логин или пароль")
	}
	return user, nil
}
