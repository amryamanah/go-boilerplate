package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockstore "github.com/amryamanah/go-boilerplate/internal/store/mock"
	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4/zero"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      store.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg store.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockStore *mockstore.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":     user.Email,
				"phone":     user.Phone,
				"password":  password,
				"full_name": user.FullName,
			},
			buildStubs: func(mockStore *mockstore.MockStore) {
				arg := store.CreateUserParams{
					Email:    user.Email,
					Phone:    user.Phone,
					FullName: user.FullName,
				}
				mockStore.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
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

			//Marshall body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			app.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user store.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser store.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.FullName.String, gotUser.FullName.String)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Phone.String, gotUser.Phone.String)
	require.Empty(t, gotUser.HashedPassword)
}

func randomUser(t *testing.T) (user store.User, password string) {
	password = util.RandomString(7)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = store.User{
		FullName:       zero.StringFrom(util.RandomOwner()),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		Phone:          zero.StringFrom("081289361989"),
	}
	return
}
