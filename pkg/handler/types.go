package handler

// Pagination represents the pagination parameters for a request.
type Pagination struct {
	// Page is the page number, starting from 1.
	// Required: true, Minimum value: 1
	Page int `json:"page" binding:"required,gte=1"`
	// PerPage is the number of items per page.
	// Required: true, Minimum value: 1, Maximum value: 300
	PerPage int `json:"perPage" binding:"required,gte=1,lte=300"`
}

// Search represents the search criteria for a request.
type Search struct {
	// Keyword is the keyword to search for.
	// Optional: true
	Keyword string `json:"keyword,omitempty"`
}
