package auth

type Expose struct{}

type Dependency struct{}

func New(_ Dependency) (*Expose, error) {
	return &Expose{}, nil
}
