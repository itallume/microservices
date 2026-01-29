package domain

import "fmt"

type Item struct {
	ProductCode string   `json:"product_code"`
	UnitPrice   float32  `json:"unit_price"`
	Quantity	int32    `json:"quantity"`
}

type Shipping struct{
	ID int64 `json:"id"`
	BillId int64 `json:"bill_id"`
	Items []Item `json:"items"`
	CreatedAt int64 `json:"created_at"`
}

func NewShipping(billId int64, items []Item) Shipping{
	return Shipping{
		BillId: billId,
		Items: items,
	}
}

func (s *Shipping) ItemsQuantity()(int32){
	var totalQuantity int32 = 0
	for _, item := range s.Items{
		totalQuantity += item.Quantity
	}
	return totalQuantity
}

func (s *Shipping) CalculateDeliveryTime()(int32, error){
	totalQuantity := s.ItemsQuantity()
	
	if totalQuantity == 0 {
		return 0, fmt.Errorf("no items in shipping")
	}

	totalQuantity = totalQuantity / 5
	if totalQuantity < 1 {
		return 1, nil
	}
	return totalQuantity, nil
}


