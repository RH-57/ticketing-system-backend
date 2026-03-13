package structs

type MonthlyTicket struct {
	Month string `json:"month"`
	Total int64  `json:"total"`
}

type DashboardStatsResponse struct {
	TotalTickets    int64           `json:"total_tickets"`
	TicketsThisYear int64           `json:"tickets_this_year"`
	OpenTickets     int64           `json:"open_tickets"`
	ResolvedTickets int64           `json:"resolved_tickets"`
	SLA             float64         `json:"sla"`
	MonthlyTrend    []MonthlyTicket `json:"monthly_trend"`
}
