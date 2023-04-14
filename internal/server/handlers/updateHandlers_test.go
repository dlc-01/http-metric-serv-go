package handlers

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	type args struct {
		gin *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateHandler(tt.args.gin)
		})
	}
}
