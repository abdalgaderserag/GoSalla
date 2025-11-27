package gosalla

import (
	"fmt"
	"time"
)

// OrdersService handles communication with the order-related endpoints
type OrdersService struct {
	client *Client
}

// Order represents a Salla order
type Order struct {
	ID              int                    `json:"id"`
	ReferenceID     string                 `json:"reference_id"`
	Status          string                 `json:"status"`
	PaymentStatus   string                 `json:"payment_status"`
	Amount          OrderAmount            `json:"amount"`
	Customer        OrderCustomer          `json:"customer"`
	ShippingAddress Address                `json:"shipping_address,omitempty"`
	BillingAddress  Address                `json:"billing_address,omitempty"`
	Items           []OrderItem            `json:"items"`
	Payment         OrderPayment           `json:"payment,omitempty"`
	Shipping        OrderShipping          `json:"shipping,omitempty"`
	Notes           string                 `json:"notes,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// OrderAmount represents order monetary values
type OrderAmount struct {
	Total         float64 `json:"total"`
	Subtotal      float64 `json:"subtotal"`
	Tax           float64 `json:"tax"`
	Shipping      float64 `json:"shipping"`
	Discount      float64 `json:"discount"`
	CurrencyCode  string  `json:"currency_code"`
}

// OrderCustomer represents customer information in an order
type OrderCustomer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
}

// Address represents a shipping or billing address
type Address struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address1    string `json:"address_1"`
	Address2    string `json:"address_2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country"`
	Phone       string `json:"phone,omitempty"`
}

// OrderItem represents a product in an order
type OrderItem struct {
	ID         int                    `json:"id"`
	ProductID  int                    `json:"product_id"`
	Name       string                 `json:"name"`
	SKU        string                 `json:"sku,omitempty"`
	Quantity   int                    `json:"quantity"`
	Price      float64                `json:"price"`
	Total      float64                `json:"total"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

// OrderPayment represents payment information
type OrderPayment struct {
	Method      string    `json:"method"`
	Gateway     string    `json:"gateway,omitempty"`
	Transaction string    `json:"transaction,omitempty"`
	PaidAt      time.Time `json:"paid_at,omitempty"`
}

// OrderShipping represents shipping information
type OrderShipping struct {
	Method      string    `json:"method"`
	TrackingNum string    `json:"tracking_number,omitempty"`
	ShippedAt   time.Time `json:"shipped_at,omitempty"`
}

// OrdersListResponse represents the response from listing orders
type OrdersListResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Data       []Order     `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// OrderResponse represents the response for a single order
type OrderResponse struct {
	Success bool  `json:"success"`
	Code    int   `json:"code"`
	Data    Order `json:"data"`
}

// OrderReservation represents an order reservation
type OrderReservation struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderReservationsResponse represents the response from listing order reservations
type OrderReservationsResponse struct {
	Success    bool               `json:"success"`
	Code       int                `json:"code"`
	Data       []OrderReservation `json:"data"`
	Pagination *Pagination        `json:"pagination,omitempty"`
}

// List retrieves all orders with optional pagination
func (s *OrdersService) List(opts *ListOptions) ([]Order, *Pagination, error) {
	path := "/orders"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp OrdersListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}

// Get retrieves an order by ID
func (s *OrdersService) Get(id int) (*Order, error) {
	path := fmt.Sprintf("/orders/%d", id)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp OrderResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// ListReservations retrieves all current order reservations
func (s *OrdersService) ListReservations(opts *ListOptions) ([]OrderReservation, *Pagination, error) {
	path := "/orders/reservations"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp OrderReservationsResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}
