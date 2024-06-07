package model

type Order struct {
	Samples []Sample `json:"samples"`
}

type Sample struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (o *Order) CalculateTotal() float64 {
	total := 0.0
	for _, sample := range o.Samples {
		total += float64(sample.Quantity) * sample.Price
	}
	return total
}
