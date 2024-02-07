package config

type Database struct {
	TokenTable map[string]bool
}

func NewDatabase() *Database {
	return &Database{
		TokenTable: map[string]bool{
			"8a395ccb-7f3e-4a5a-b35c-4fea034d24f2": true,
		},
	}
}

func (d *Database) FindOne(key string) bool {
	if _, ok := d.TokenTable[key]; !ok {
		return false
	}

	return true
}
