package service

import (
	"database/sql"
	"errors"
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	servicemock "github.com/gavrl/app/test/gomock/mocks/servicemock"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBalanceService(t *testing.T) {
	assertions := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := servicemock.NewMockBalanceRepository(ctrl)
	balanceServiceImpl := NewBalanceService(balanceRepo, validator.New())

	assertions.Implements((*Balance)(nil), new(BalanceService), "Balance Service Implementation does not honor service definition")
	assertions.NotNil(balanceServiceImpl, "Balance Service not initialized")
	assertions.NotNil(balanceServiceImpl.repo, "balance service dependency repository not initialized")
	assertions.NotNil(balanceServiceImpl.validator, "balance service dependency validator not initialized")
}

func TestRefill(t *testing.T) {
	req := require.New(t)

	type deps struct {
		balanceRepo *servicemock.MockBalanceRepository
	}
	type args struct {
		moveMoneyModel model.MoveMoneyModel
	}
	tests := []struct {
		name         string
		prepare      func(f *deps)
		input        args
		want         int
		isError      bool
		errorMessage string
	}{
		{
			name: "not found customer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{}, sql.ErrNoRows)
				f.balanceRepo.EXPECT().Create(1, 1200).Times(1).Return(1200, nil)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1200},
			},
			want:         1200,
			isError:      false,
			errorMessage: "",
		},
		{
			name: "negative amount",
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: -1200},
			},
			want:         0,
			isError:      true,
			errorMessage: "Amount",
		},
		{
			name: "found customer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().UpdateAmount(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     2200,
				}).Times(1).Return(nil)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         2200,
			isError:      false,
			errorMessage: "",
		},
		{
			name: "database error from GetByCustomerId()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(
					internal.Balance{},
					errors.New("unexpected error from database"),
				)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         0,
			isError:      true,
			errorMessage: "error from database",
		},
		{
			name: "database error from create()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{}, sql.ErrNoRows)
				f.balanceRepo.EXPECT().Create(1, 1200).Times(1).
					Return(0, errors.New("unexpected error from database"))
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1200},
			},
			want:         0,
			isError:      true,
			errorMessage: "error from database",
		},
		{
			name: "database error from update()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().UpdateAmount(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     2400,
				}).Times(1).Return(errors.New("unexpected error from database"))
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1200},
			},
			want:         0,
			isError:      true,
			errorMessage: "error from database",
		},
	}

	for _, cs := range tests {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := deps{
				balanceRepo: servicemock.NewMockBalanceRepository(ctrl),
			}
			if cs.prepare != nil {
				cs.prepare(&f)
			}
			s := &BalanceService{
				repo:      f.balanceRepo,
				validator: validator.New(),
			}

			res, err := s.Refill(&cs.input.moveMoneyModel)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errorMessage)
			} else {
				req.NoError(err)
			}

			req.Equal(cs.want, res)
		})
	}

}

func TestWithdraw(t *testing.T) {
	req := require.New(t)

	type deps struct {
		balanceRepo *servicemock.MockBalanceRepository
	}
	type args struct {
		moveMoneyModel model.MoveMoneyModel
	}
	tests := []struct {
		name         string
		prepare      func(f *deps)
		input        args
		want         int
		isError      bool
		errorMessage string
	}{
		{
			name: "not found customer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{}, sql.ErrNoRows)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         0,
			isError:      true,
			errorMessage: "customer with id 1 does not exists",
		},
		{
			name: "negative amount",
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: -1200},
			},
			want:         0,
			isError:      true,
			errorMessage: "Field validation for 'Amount' failed on the 'gt' tag",
		},
		{
			name: "not enough funds to write off",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1000,
				}, nil)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1200},
			},
			want:         0,
			isError:      true,
			errorMessage: "there are not enough funds",
		},
		{
			name: "found customer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().UpdateAmount(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     200,
				}).Times(1).Return(nil)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         200,
			isError:      false,
			errorMessage: "",
		},
		{
			name: "database error from GetByCustomerId()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(
					internal.Balance{},
					errors.New("unexpected error from database"),
				)
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         0,
			isError:      true,
			errorMessage: "error from database",
		},
		{
			name: "database error from update()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().UpdateAmount(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     200,
				}).Times(1).Return(errors.New("unexpected error from database"))
			},
			input: args{
				moveMoneyModel: model.MoveMoneyModel{CustomerId: 1, Amount: 1000},
			},
			want:         0,
			isError:      true,
			errorMessage: "error from database",
		},
	}

	for _, cs := range tests {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := deps{
				balanceRepo: servicemock.NewMockBalanceRepository(ctrl),
			}
			if cs.prepare != nil {
				cs.prepare(&f)
			}
			s := &BalanceService{
				repo:      f.balanceRepo,
				validator: validator.New(),
			}

			res, err := s.Withdraw(&cs.input.moveMoneyModel)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errorMessage)
			} else {
				req.NoError(err)
			}

			req.Equal(cs.want, res)
		})
	}
}

func TestTransfer(t *testing.T) {
	req := require.New(t)

	type deps struct {
		balanceRepo *servicemock.MockBalanceRepository
	}
	type args struct {
		transferMoneyModel model.TransferMoneyModel
	}
	tests := []struct {
		name         string
		prepare      func(f *deps)
		input        args
		isError      bool
		errorMessage string
	}{
		{
			name: "not found customer from",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{}, sql.ErrNoRows)
			},
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: 1000},
			},
			isError:      true,
			errorMessage: "customer with id 1 does not exists",
		},
		{
			name: "not found customer to",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().GetByCustomerId(2).Times(1).Return(internal.Balance{}, sql.ErrNoRows)
			},
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: 1000},
			},
			isError:      true,
			errorMessage: "customer with id 2 does not exists",
		},
		{
			name: "negative amount",
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: -1000},
			},
			isError:      true,
			errorMessage: "Field validation for 'Amount' failed on the 'gt' tag",
		},
		{
			name: "not enough funds to transfer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1000,
				}, nil)
				f.balanceRepo.EXPECT().GetByCustomerId(2).Times(1).Return(internal.Balance{
					Id:         2,
					CustomerId: 2,
					Amount:     1200,
				}, nil)
			},
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: 3000},
			},
			isError:      true,
			errorMessage: "there are not enough funds",
		},
		{
			name: "database error from transfer()",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1000,
				}, nil)
				f.balanceRepo.EXPECT().GetByCustomerId(2).Times(1).Return(internal.Balance{
					Id:         2,
					CustomerId: 2,
					Amount:     1200,
				}, nil)

				f.balanceRepo.EXPECT().TransferMoney(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     500,
				}, &internal.Balance{
					Id:         2,
					CustomerId: 2,
					Amount:     1700,
				}).Times(1).Return(errors.New("unexpected error from database"))
			},
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: 500},
			},
			isError:      true,
			errorMessage: "error from database",
		},
		{
			name: "success transfer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1000,
				}, nil)
				f.balanceRepo.EXPECT().GetByCustomerId(2).Times(1).Return(internal.Balance{
					Id:         2,
					CustomerId: 2,
					Amount:     1200,
				}, nil)

				f.balanceRepo.EXPECT().TransferMoney(&internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     500,
				}, &internal.Balance{
					Id:         2,
					CustomerId: 2,
					Amount:     1700,
				}).Times(1).Return(nil)
			},
			input: args{
				transferMoneyModel: model.TransferMoneyModel{CustomerIdFrom: 1, CustomerIdTo: 2, Amount: 500},
			},
			isError:      false,
			errorMessage: "",
		},
	}

	for _, cs := range tests {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := deps{
				balanceRepo: servicemock.NewMockBalanceRepository(ctrl),
			}
			if cs.prepare != nil {
				cs.prepare(&f)
			}
			s := &BalanceService{
				repo:      f.balanceRepo,
				validator: validator.New(),
			}

			err := s.Transfer(&cs.input.transferMoneyModel)

			if cs.isError {
				req.Error(err)
				req.Contains(err.Error(), cs.errorMessage)
			} else {
				req.NoError(err)
			}
		})
	}
}
