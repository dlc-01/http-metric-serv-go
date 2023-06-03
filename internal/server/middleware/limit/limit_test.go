package limit

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Limit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Limit(1))
	router.GET("/", func(*gin.Context) {
		time.Sleep(500 * time.Microsecond)
	})

	w := performRequest("GET", "/", router)
	assert.Equal(t, 200, w.Code)

}
func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

func TestFulled(t *testing.T) {
	const max = 5

	attempts := 1000
	var failed int
	var wg sync.WaitGroup

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Limit(max))
	router.GET("/", func(*gin.Context) {
		time.Sleep(5 * time.Microsecond)
	})

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := performRequest("GET", "/", router)
			if w.Code == 502 {
				failed++
			}
		}()
	}
	wg.Wait()

	assert.True(t, failed > 0)
}
