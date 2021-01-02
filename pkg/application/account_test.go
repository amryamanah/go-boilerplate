package application

import (
	mockstore "github.com/amryamanah/go-boilerplate/internal/store/mock"
	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/test_util"
	"github.com/golang/mock/gomock"
	"testing"
)

func randomAccount() store.Account {
	return store.Account{
		ID: test_util.RandomInt(1, 1000),
		Owner: test_util.RandomOwner(),
		Balance: test_util.RandomMoney(),
		Currency: test_util.RandomCurrency(),
	}
}
func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockstore.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	NewApplication()
}
