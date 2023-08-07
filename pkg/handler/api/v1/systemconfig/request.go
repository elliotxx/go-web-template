package systemconfig

import "github.com/elliotxx/go-web-template/pkg/handler"

// CreateSystemConfigRequest represents the create request structure for
// configuration of a system.
type CreateSystemConfigRequest struct {
	// Tenant or organization that the system belongs to
	Tenant string `json:"tenant" binding:"required"`
	// Environment where the system is deployed (e.g. prod, gray)
	Env string `json:"env" binding:"required"`
	// Type or category of the system (e.g. cache, message queue)
	Type string `json:"type" binding:"required"`
	// Configuration data in JSON or YAML format
	Config string `json:"config" binding:"required"`
	// Description or purpose of the system
	Description string `json:"description"`
	// Username or ID of the user who created the system
	Creator string `json:"creator" binding:"required"`
	// Username or ID of the user who last modified the system
	Modifier string `json:"modifier"`
}

// UpdateSystemConfigRequest represents the update request structure for
// configuration of a system.
type UpdateSystemConfigRequest struct {
	// Unique ID of the system
	ID uint `json:"id" binding:"required"`
	// Tenant or organization that the system belongs to
	Tenant string `json:"tenant"`
	// Environment where the system is deployed (e.g. prod, gray)
	Env string `json:"env"`
	// Type or category of the system (e.g. cache, message queue)
	Type string `json:"type"`
	// Configuration data in JSON or YAML format
	Config string `json:"config"`
	// Description or purpose of the system
	Description string `json:"description"`
	// Username or ID of the user who created the system
	Creator string `json:"creator"`
	// Username or ID of the user who last modified the system
	Modifier string `json:"modifier"`
}

// QuerySystemConfigRequest represents the query request structure for
// configuration of a system.
type QuerySystemConfigRequest struct {
	handler.Pagination
	handler.Search
}
