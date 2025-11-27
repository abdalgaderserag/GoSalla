package gosalla

import (
	"fmt"
	"time"
)

// CustomersService handles communication with the customer-related endpoints
type CustomersService struct {
	client *Client
}

// Customer represents a Salla customer
type Customer struct {
	ID          int                    `json:"id"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone,omitempty"`
	Gender      string                 `json:"gender,omitempty"`
	DateOfBirth string                 `json:"date_of_birth,omitempty"`
	Status      string                 `json:"status"`
	Avatar      string                 `json:"avatar,omitempty"`
	Addresses   []CustomerAddress      `json:"addresses,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CustomerAddress represents a customer's address
type CustomerAddress struct {
	ID          int    `json:"id"`
	Type        string `json:"type,omitempty"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address1    string `json:"address_1"`
	Address2    string `json:"address_2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country"`
	Phone       string `json:"phone,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty"`
}

// CustomersListResponse represents the response from listing customers
type CustomersListResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Data       []Customer  `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// CustomerResponse represents the response for a single customer
type CustomerResponse struct {
	Success bool     `json:"success"`
	Code    int      `json:"code"`
	Data    Customer `json:"data"`
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone,omitempty"`
	Password    string                 `json:"password,omitempty"`
	Gender      string                 `json:"gender,omitempty"`
	DateOfBirth string                 `json:"date_of_birth,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateCustomerRequest represents the request to update a customer
type UpdateCustomerRequest struct {
	FirstName   string                 `json:"first_name,omitempty"`
	LastName    string                 `json:"last_name,omitempty"`
	Email       string                 `json:"email,omitempty"`
	Phone       string                 `json:"phone,omitempty"`
	Gender      string                 `json:"gender,omitempty"`
	DateOfBirth string                 `json:"date_of_birth,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// List retrieves all customers with optional pagination
func (s *CustomersService) List(opts *ListOptions) ([]Customer, *Pagination, error) {
	path := "/customers"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp CustomersListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}

// Get retrieves a customer by ID
func (s *CustomersService) Get(id int) (*Customer, error) {
	path := fmt.Sprintf("/customers/%d", id)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp CustomerResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Create creates a new customer
func (s *CustomersService) Create(customer *CreateCustomerRequest) (*Customer, error) {
	path := "/customers"
	
	req, err := s.client.newRequest("POST", path, customer)
	if err != nil {
		return nil, err
	}
	
	var resp CustomerResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Update updates an existing customer
func (s *CustomersService) Update(id int, customer *UpdateCustomerRequest) (*Customer, error) {
	path := fmt.Sprintf("/customers/%d", id)
	
	req, err := s.client.newRequest("PUT", path, customer)
	if err != nil {
		return nil, err
	}
	
	var resp CustomerResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}
