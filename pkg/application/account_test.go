package application

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockstore "github.com/amryamanah/go-boilerplate/internal/store/mock"
	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func randomAccount() store.Account {
	return store.Account{
		ID: util.RandomInt(1, 1000),
		Owner: 1,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account store.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount store.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(store.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			accountID: account.ID,
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(store.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			accountID: 0,
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockstore.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			testApp := NewApplication(mockStore)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testApp.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockstore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(mockStore *mockstore.MockStore) {
				arg := store.CreateAccountParams{
					Owner:    account.Owner,
					Currency:  account.Currency,
					Balance: 0,
				}

				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"owner": account.Owner,
				"currency": "invalid",
			},
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidOwner",
			body: gin.H{
				"owner": "",
				"currency": account.Currency,
			},
			buildStubs: func(mockStore *mockstore.MockStore) {
				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockstore.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			app := NewApplication(mockStore)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			app.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


