package internal

type Balance struct {
	Id         int    `json:"-" db:"id"`
	CustomerId string `json:"customer_id" binding:"required"`
	Balance    string `json:"balance" binding:"required"`
}
