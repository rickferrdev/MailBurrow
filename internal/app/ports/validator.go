package ports

type Validator interface {
	Validate(out any) error
}
