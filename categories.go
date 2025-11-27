package gosalla

import (
	"fmt"
	"time"
)

// CategoriesService handles communication with the category-related endpoints
type CategoriesService struct {
	client *Client
}

// Category represents a Salla product category
type Category struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	ParentID    int                    `json:"parent_id,omitempty"`
	Image       string                 `json:"image,omitempty"`
	Status      string                 `json:"status"`
	SortOrder   int                    `json:"sort_order,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CategoriesListResponse represents the response from listing categories
type CategoriesListResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Data       []Category  `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// CategoryResponse represents the response for a single category
type CategoryResponse struct {
	Success bool     `json:"success"`
	Code    int      `json:"code"`
	Data    Category `json:"data"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	ParentID    int                    `json:"parent_id,omitempty"`
	Image       string                 `json:"image,omitempty"`
	Status      string                 `json:"status,omitempty"`
	SortOrder   int                    `json:"sort_order,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	ParentID    int                    `json:"parent_id,omitempty"`
	Image       string                 `json:"image,omitempty"`
	Status      string                 `json:"status,omitempty"`
	SortOrder   int                    `json:"sort_order,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// List retrieves all categories with optional pagination
func (s *CategoriesService) List(opts *ListOptions) ([]Category, *Pagination, error) {
	path := "/categories"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp CategoriesListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}

// Get retrieves a category by ID
func (s *CategoriesService) Get(id int) (*Category, error) {
	path := fmt.Sprintf("/categories/%d", id)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp CategoryResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Create creates a new category
func (s *CategoriesService) Create(category *CreateCategoryRequest) (*Category, error) {
	path := "/categories"
	
	req, err := s.client.newRequest("POST", path, category)
	if err != nil {
		return nil, err
	}
	
	var resp CategoryResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Update updates an existing category
func (s *CategoriesService) Update(id int, category *UpdateCategoryRequest) (*Category, error) {
	path := fmt.Sprintf("/categories/%d", id)
	
	req, err := s.client.newRequest("PUT", path, category)
	if err != nil {
		return nil, err
	}
	
	var resp CategoryResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Delete deletes a category
func (s *CategoriesService) Delete(id int) error {
	path := fmt.Sprintf("/categories/%d", id)
	
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	
	return s.client.do(req, nil)
}
