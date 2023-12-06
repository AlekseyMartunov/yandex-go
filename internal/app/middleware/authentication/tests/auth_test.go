package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mockAuthentication "github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication/tests/mocks"
)

func TestCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorage := mockAuthentication.NewMockuserStorage(ctrl)
	mockUserStorage.EXPECT().SaveNewUser().Return(1, nil)
	mockUserStorage.EXPECT().SaveNewUser().Return(2, nil)

	tesCase := []struct {
		Name       string
		ExpectedID string
	}{
		{
			Name:       "test1",
			ExpectedID: "1",
		},
		{
			Name:       "test2",
			ExpectedID: "2",
		},
	}

	for _, tc := range tesCase {
		t.Run(tc.Name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id := r.Header.Get("UserID")
				assert.Equal(t, id, tc.ExpectedID)
			})

			tk := authentication.NewTokenController(mockUserStorage)
			testNextHandler := tk.CheckToken(testHandler)
			req := httptest.NewRequest("GET", "http://testing", nil)

			testNextHandler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}

}
