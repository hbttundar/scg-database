package contract

type (
	Seeder interface{ Run(db Connection) error }
)
