package model

import "sort"

type Item struct {
	Name string `json:"name"`
}

// Order represents a customer's order
type Order struct {
	Amount int   `json:"amount"`
	Result []int `json:"result"`
}

// Packages represents all available package sizes
type Packages struct {
	Sizes []int `json:"sizes"`
}

// Function to find the least amount of packages
func (p *Packages) FindPackages(order *Order, packageSizes *Packages) {
	sort.Sort(sort.Reverse(sort.IntSlice(packageSizes.Sizes))) // Sort in descending order for easier comparison

	var result []int
	remainingOrder := order.Amount

	for _, size := range packageSizes.Sizes {
		for remainingOrder >= size {
			remainingOrder -= size
			result = append(result, size)
		}
	}

	// sort in ascending order of the package sizes
	sort.Sort(sort.IntSlice(packageSizes.Sizes))

	// If there's still some order left that needs to be fulfilled
	if remainingOrder > 0 {
		for _, size := range packageSizes.Sizes {
			if remainingOrder <= size {
				result = append(result, size)
				break
			}
		}
	}

	// Optimize the package sizes at the end
	order.Result = p.OptimizePackages(result, packageSizes.Sizes)
}

func (p *Packages) OptimizePackages(packages, packageSizes []int) []int {
	totalSize := 0
	for _, size := range packages {
		totalSize += size
	}

	combinations := make([][]int, totalSize+1)
	combinations[0] = []int{}

	for _, size := range packageSizes {
		for i := 0; i <= totalSize-size; i++ {
			if combinations[i] != nil && (combinations[i+size] == nil || len(combinations[i])+1 < len(combinations[i+size])) {
				combination := make([]int, len(combinations[i]), len(combinations[i])+1)
				copy(combination, combinations[i])
				combination = append(combination, size)
				combinations[i+size] = combination
			}
		}
	}

	return combinations[totalSize]
}
