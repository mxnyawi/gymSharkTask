package model

import (
	"reflect"
	"testing"
)

func TestFindPackages(t *testing.T) {
	packageSizes := &Packages{Sizes: []int{250, 500, 1000, 2000, 5000}}
	tests := []struct {
		name  string
		order *Order
		want  []int
	}{
		{
			name:  "Order 251",
			order: &Order{Amount: 251},
			want:  []int{500},
		},
		{
			name:  "Order 501",
			order: &Order{Amount: 501},
			want:  []int{250, 500},
		},
		{
			name:  "Order 1001",
			order: &Order{Amount: 1001},
			want:  []int{250, 1000},
		},
		{
			name:  "Order 12001",
			order: &Order{Amount: 12001},
			want:  []int{250, 2000, 5000, 5000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packageSizes.FindPackages(tt.order, packageSizes)
			if !reflect.DeepEqual(tt.order.Result, tt.want) {
				t.Errorf("FindPackages() = %v, want %v", tt.order.Result, tt.want)
			}
		})
	}
}
