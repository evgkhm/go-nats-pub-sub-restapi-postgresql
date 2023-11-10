package usecase

import (
	"context"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres/mocks"
	mocks2 "go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions/mocks"
	"log/slog"
	"testing"
)

func TestUserUseCase_CheckNegativeBalance(t *testing.T) {
	type args struct {
		ctx     context.Context
		userDTO *user.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				ctx: context.Background(),
				userDTO: &user.User{
					ID:      1,
					Balance: -10,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewUserRepository(t)
			mockTx := mocks2.NewTransactor(t)

			u := &UserUseCase{
				userRepo:  mockRepo,
				txService: mockTx,
				logger:    slog.Default(),
			}
			if err := u.CheckNegativeBalance(tt.args.ctx, tt.args.userDTO); (err != nil) != tt.wantErr {
				t.Errorf("CheckNegativeBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
