package structs

type TicketListResponse struct {
	ID           uint   `json:"id"`
	TicketNumber string `json:"ticket_number"`
	Title        string `json:"title"`
	Priority     string `json:"priority"`
	EmployeeName string `json:"employee_name"`
	Status       string `json:"status"`
	CreatedBy    string `json:"created_by"`
}

type TicketDetailResponse struct {
	ID           uint   `json:"id"`
	TicketNumber string `json:"ticket_number"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Priority     string `json:"priority"`

	Employee struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"employee"`

	CreatedBy struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"created_by"`

	Branch     string `json:"branch"`
	Division   string `json:"division"`
	Department string `json:"department"`

	CreatedAt string `json:"created_at"`
}

type TicketCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EmployeeID  uint   `json:"employee_id" binding:"required"`

	Priority string `json:"priority" binding:"required,oneof=LOW MEDIUM HIGH"`
}

type TicketUpdateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EmployeeID  uint   `json:"employee_id" binding:"required"`

	Status   string `json:"status" binding:"required,oneof=OPEN PROCESS PENDING CLOSED"`
	Priority string `json:"priority" binding:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
}

type TicketUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=OPEN PROCESS PENDING CLOSED"`
}
