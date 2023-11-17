package tests

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication"
	"net/http"
	"net/http/httptest"
	"testing"

	mockAuthentication "github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication/tests/mocks"
	"github.com/golang/mock/gomock"
)

func BenchmarkToken(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockUserStorage := mockAuthentication.NewMockuserStorage(ctrl)
	mockUserStorage.EXPECT().SaveNewUser().Return(1, nil)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	tk := authentication.NewTokenController(mockUserStorage)
	testNextHandler := tk.CheckToken(testHandler)
	req := httptest.NewRequest("GET", "http://testing", nil)

	b.ResetTimer()

	testNextHandler.ServeHTTP(httptest.NewRecorder(), req)

}
