package limit

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	type args struct {
		limit    int
		attempts int
		goodTest bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "testOneLimit",
			args: args{
				limit:    1,
				attempts: 1,
				goodTest: true,
			},
		},
		{name: "testMoreLimit",
			args: args{
				limit:    5,
				attempts: 100,
				goodTest: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var failed int
			var wg sync.WaitGroup

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(Limit(tt.args.limit))
			router.GET("/", func(*gin.Context) {
				time.Sleep(5 * time.Microsecond)
			})
			if tt.args.goodTest {
				w := performRequest("GET", "/", router)
				assert.Equal(t, 200, w.Code)
			} else {
				for i := 0; i < tt.args.attempts; i++ {
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

		})
	}
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
