package http

type Todo struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// for request and response Create.
type (
	CreateRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	CreateResponse struct {
		ID uint64 `json:"id"`
	}
)

// for request and response Delete.
type (
	DeleteRequest struct {
		ID uint64 `json:"id"`
	}

	DeleteResponse struct {
		ID uint64 `json:"id"`
	}
)

// for request and response GetByID.
type (
	GetByIDRequest struct {
		ID uint64 `json:"id"`
	}

	GetByIDResponse struct {
		ID          uint64 `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
)

// for request and response GetWithFilter.
type (
	GetWithFilterRequest struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	GetWithFilterResponse struct {
		Todos []Todo `json:"todos"`
	}
)

// for request and response UpdateStatus.
type (
	UpdateStatusRequest struct {
		Status string `json:"status"`
	}

	UpdateStatusResponse struct {
		ID     uint64 `json:"id"`
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
		ID          uint64 `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
)
