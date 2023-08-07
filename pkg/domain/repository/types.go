package repository

// Query represents the query criteria for a database access.
type Query struct {
	// Offset is the number of items to skip.
	Offset int
	// Limit is the maximum number of items to return.
	Limit int
	// Keyword is the keyword to search for.
	Keyword string
}
