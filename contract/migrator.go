package contract

type (
	Migrator interface {
		Up() error
		Down(steps int) error
		Fresh() error
		Close() (sourceErr, dbErr error)
	}
)
