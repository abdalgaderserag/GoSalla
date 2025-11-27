package gosalla

import (
	"fmt"
	"time"
)

// BrandsService handles communication with the brand-related endpoints
type BrandsService struct {
	client *Client
}

// Brand represents a Salla brand
type Brand struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Logo        string                 `json:"logo,omitempty"`
	Website     string                 `json:"website,omitempty"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// BrandsListResponse represents the response from listing brands
type BrandsListResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Data       []Brand     `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// BrandResponse represents the response for a single brand
type BrandResponse struct {
	Success bool  `json:"success"`
	Code    int   `json:"code"`
	Data    Brand `json:"data"`
}

// CreateBrandRequest represents the request to create a brand
type CreateBrandRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Logo        string                 `json:"logo,omitempty"`
	Website     string                 `json:"website,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateBrandRequest represents the request to update a brand
type UpdateBrandRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Logo        string                 `json:"logo,omitempty"`
	Website     string                 `json:"website,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// List retrieves all brands with optional pagination
func (s *BrandsService) List(opts *ListOptions) ([]Brand, *Pagination, error) {
	path := "/brands"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp BrandsListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}

// Get retrieves a brand by ID
func (s *BrandsService) Get(id int) (*Brand, error) {
	path := fmt.Sprintf("/brands/%d", id)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp BrandResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Create creates a new brand
func (s *BrandsService) Create(brand *CreateBrandRequest) (*Brand, error) {
	path := "/brands"
	
	req, err := s.client.newRequest("POST", path, brand)
	if err != nil {
		return nil, err
	}
	
	var resp BrandResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Update updates an existing brand
func (s *BrandsService) Update(id int, brand *UpdateBrandRequest) (*Brand, error) {
	path := fmt.Sprintf("/brands/%d", id)
	
	req, err := s.client.newRequest("PUT", path, brand)
	if err != nil {
		return nil, err
	}
	
	var resp BrandResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Delete deletes a brand
func (s *BrandsService) Delete(id int) error {
	path := fmt.Sprintf("/brands/%d", id)
	
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	
	return s.client.do(req, nil)
}
