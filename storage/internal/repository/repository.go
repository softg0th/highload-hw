package repository

type Repository struct {
	DB *DataBase
}

func NewRepository(db *DataBase) *Repository {
	return &Repository{
		DB: db,
	}
}
