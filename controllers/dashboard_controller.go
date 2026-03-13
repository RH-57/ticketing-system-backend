package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var total, thisYear, open, resolved int64
	var sla float64 = 0
	currentYear := time.Now().Year()

	// 1. Hitung Total Semua Tiket
	database.DB.Model(&models.Ticket{}).Count(&total)

	// 2. Hitung Tiket Khusus Tahun Ini
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)
	database.DB.Model(&models.Ticket{}).
		Where("created_at >= ?", startOfYear).
		Count(&thisYear)

	// 3. Hitung Tiket dengan Status OPEN
	database.DB.Model(&models.Ticket{}).
		Where("status = ?", "OPEN").
		Count(&open)

	// 4. Hitung Tiket dengan Status CLOSED (Resolved)
	database.DB.Model(&models.Ticket{}).
		Where("status = ?", "CLOSED").
		Count(&resolved)

	// 5. Hitung SLA
	if total > 0 {
		sla = (float64(resolved) / float64(total)) * 100
	}

	// 6. Query Tren Bulanan (Versi PostgreSQL)
	var monthlyTrend []structs.MonthlyTicket
	database.DB.Model(&models.Ticket{}).
		// TO_CHAR(created_at, 'Mon') menghasilkan 'Jan', 'Feb', dst.
		// Kita juga ambil angka bulannya untuk sorting yang benar
		Select("TO_CHAR(created_at, 'Mon') as month, count(*) as total").
		Where("EXTRACT(YEAR FROM created_at) = ?", currentYear).
		Group("TO_CHAR(created_at, 'Mon'), EXTRACT(MONTH FROM created_at)").
		Order("EXTRACT(MONTH FROM created_at) ASC").
		Scan(&monthlyTrend)

	// Kirim Response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Dashboard statistics fetched successfully",
		Data: structs.DashboardStatsResponse{
			TotalTickets:    total,
			TicketsThisYear: thisYear,
			OpenTickets:     open,
			ResolvedTickets: resolved,
			SLA:             sla,
			MonthlyTrend:    monthlyTrend,
		},
	})
}
