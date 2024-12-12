package inbound

type Todo struct {
	ID          uint64 `json:"id,string"`
	UserID      uint64 `json:"user_id,string"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor"`
	HashMore   bool   `json:"has_more"`
}

// for request and response Create.
type (
	CreateRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	CreateResponse struct {
		ID uint64 `json:"id,string"`
	}
)

// for request and response Delete.
type (
	DeleteRequest struct {
		ID uint64 `json:"id"`
	}

	DeleteResponse struct {
		ID uint64 `json:"id,string"`
	}
)

// for request and response Find.
type (
	FindRequest struct {
		ID uint64 `json:"id"`
	}

	FindResponse struct {
		ID          uint64 `json:"id,string"`
		UserID      uint64 `json:"user_id,string"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
)

// for request and response Fetch.
type (
	FetchRequest struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	FetchResponse struct {
		Todos      []Todo     `json:"todos"`
		Pagination Pagination `json:"pagination"`
	}
)

// for request and response UpdateStatus.
type (
	UpdateStatusRequest struct {
		Status string `json:"status"`
	}

	UpdateStatusResponse struct {
		ID     uint64 `json:"id,string"`
		Status string `json:"status"`
	}
)

// for request and response Update.
type (
	UpdateRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	UpdateResponse struct {
		ID          uint64 `json:"id,string"`
		UserID      uint64 `json:"user_id,string"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
)
