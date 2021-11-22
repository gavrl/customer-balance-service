package service

import (
	"database/sql"
	"errors"
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	servicemock "github.com/gavrl/app/test/gomock/mocks/servicemock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBalanceService(t *testing.T) {
	assertions := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := servicemock.NewMockBalanceRepository(ctrl)
	balanceServiceImpl := NewBalanceService(balanceRepo)

	assertions.Implements((*Balance)(nil), new(BalanceService), "Balance Service Implementation does not honor service definition")
	assertions.NotNil(balanceServiceImpl, "Product Service not initialized")
	assertions.NotNil(balanceServiceImpl.repo, "Product Service dependency not initialized")
}

func TestRefill(t *testing.T) {
	req := require.New(t)

	type deps struct {
		balanceRepo *servicemock.MockBalanceRepository
	}
	type args struct {
		refillDto model.RefillDto
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
				refillDto: model.RefillDto{CustomerId: 1, Amount: model.PennyAmount{Int: 1200}},
			},
			want:         1200,
			isError:      false,
			errorMessage: "",
		},
		{
			name: "found customer",
			prepare: func(f *deps) {
				f.balanceRepo.EXPECT().GetByCustomerId(1).Times(1).Return(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     1200,
				}, nil)
				f.balanceRepo.EXPECT().UpdateAmount(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     2200,
				}).Times(1).Return(nil)
			},
			input: args{
				refillDto: model.RefillDto{CustomerId: 1, Amount: model.PennyAmount{Int: 1000}},
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
				refillDto: model.RefillDto{CustomerId: 1, Amount: model.PennyAmount{Int: 1000}},
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
				refillDto: model.RefillDto{CustomerId: 1, Amount: model.PennyAmount{Int: 1200}},
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
				f.balanceRepo.EXPECT().UpdateAmount(internal.Balance{
					Id:         1,
					CustomerId: 1,
					Amount:     2400,
				}).Times(1).Return(errors.New("unexpected error from database"))
			},
			input: args{
				refillDto: model.RefillDto{CustomerId: 1, Amount: model.PennyAmount{Int: 1200}},
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
				repo: f.balanceRepo,
			}

			res, err := s.Refill(cs.input.refillDto)

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
