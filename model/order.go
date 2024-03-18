package model

import (
  "time"
"github.com/google/uuid"
)

type Order struct {
  OrderID unit64 `json:"order_id"`
  CustomerID uuid.UUID `json:"customer_id"`
  LineItems []LineItem `json:"line_items"`
  CreatedAt *time.Time `json:"created_at"`
  ShippedAt *time.Time `json:"shipped_at"`
  CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
LineItem uuid.UUID `json:"line_item"`
Quantity unit `json:"quantity"`
Price uint `json:"price"`
}
