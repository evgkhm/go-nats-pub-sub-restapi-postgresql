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
