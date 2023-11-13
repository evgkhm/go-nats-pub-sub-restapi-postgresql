# SQL driver mock for Golang (with <a href="https://github.com/jmoiron/sqlx">jmoiron/sqlx</a> support)

## Forked from [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)

### Added functionality
- `Newx() and NewxWithDNS()` which returns `*sqlx.DB` object instead of `*sql.DB`

## Install

    https://github.com/zhashkevych/go-sqlxmock@master

## Usage Example

Repository Implementation:
```go
type UserRepository interface {
	Insert(user domain.User) (int, error)
	GetById(id int) (domain.User, error)
	Get(username, password string) (domain.User, error)
}


type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Insert(user domain.User) (int, error) {
	var id int
	row := r.db.QueryRow("INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.FirstName, user.LastName, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
```

Unit Tests:

```go
import (
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestUserRepository_Insert(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewUserRepository(db)

	tests := []struct {
		name    string
		s       repository.UserRepository
		user    domain.User
		mock    func()
		want    int
		wantErr bool
	}{
		{
			//When everything works as expected
			name: "OK",
			s:    s,
			user: domain.User{
				FirstName: "first_name",
				LastName:  "last_name",
				Username:  "username",
				Password:  "password",
			},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").WithArgs("first_name", "last_name", "username", "password").WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name:  "Empty Fields",
			s:     s,
			user: domain.User{
				FirstName: "",
				LastName:  "",
				Username:  "username",
				Password:  "password",
			},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"id"}) 
				mock.ExpectQuery("INSERT INTO users").WithArgs("first_name", "last_name", "username", "password").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Insert(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
```