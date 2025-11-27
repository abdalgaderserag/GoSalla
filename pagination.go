package gosalla

// Pagination represents pagination metadata from API responses
type Pagination struct {
	CurrentPage int `json:"current_page"`
	From        int `json:"from"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	To          int `json:"to"`
	Total       int `json:"total"`
}

// HasNextPage checks if there are more pages available
func (p *Pagination) HasNextPage() bool {
	if p == nil {
		return false
	}
	return p.CurrentPage < p.LastPage
}

// NextPage returns the next page number, or 0 if there are no more pages
func (p *Pagination) NextPage() int {
	if !p.HasNextPage() {
		return 0
	}
	return p.CurrentPage + 1
}

// HasPreviousPage checks if there is a previous page
func (p *Pagination) HasPreviousPage() bool {
	if p == nil {
		return false
	}
	return p.CurrentPage > 1
}

// PreviousPage returns the previous page number, or 0 if on the first page
func (p *Pagination) PreviousPage() int {
	if !p.HasPreviousPage() {
		return 0
	}
	return p.CurrentPage - 1
}

// ListOptions represents common options for list endpoints
type ListOptions struct {
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
}
