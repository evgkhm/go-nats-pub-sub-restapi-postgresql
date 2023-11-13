package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	//"github.com/DATA-DOG/go-sqlmock"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"testing"
)

func TestRepo_UserBalanceAccrual(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		ctx     context.Context
		userDTO *user.User
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr bool
	}{
		{
			name: "Ok_1",
			mock: func() {
				mock.ExpectExec("UPDATE user_info SET (.+) WHERE (.+)").
					WithArgs(10.0, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: 10,
				},
			},
		},
		{
			name: "Ok_2",
			mock: func() {
				mock.ExpectExec("UPDATE user_info SET (.+) WHERE (.+)").
					WithArgs(2.5, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      2,
					Balance: 2.5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err = r.UserBalanceAccrual(tt.args.ctx, &user.User{ID: tt.args.userDTO.ID, Balance: tt.args.userDTO.Balance})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_GetBalance(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		ctx     context.Context
		userDTO *user.User
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		want    float32
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10.0)
				mock.ExpectQuery("SELECT (.+) FROM user_info WHERE (.+)").
					WithArgs(1).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: 10,
				},
			},
			want: 10,
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"balance"})

				mock.ExpectQuery("SELECT (.+) FROM user_info WHERE (.+)").
					WithArgs(404).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      404,
					Balance: 10,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetBalance(tt.args.ctx, tt.args.userDTO.ID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		ctx     context.Context
		userDTO *user.User
	}
	type mockBehavior func(args args, id int)
	tests := []struct {
		name    string
		args    args
		mock    mockBehavior
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: 10,
				},
			},
			mock: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO user_info").
					WithArgs(1, 10.0).
					WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name: "Duplicate error",
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: 10,
				},
			},
			mock: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, ErrUserAlreadyExist)
				mock.ExpectQuery("INSERT INTO user_info").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args, tt.want)
			err := r.CreateUser(tt.args.ctx, tt.args.userDTO)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
