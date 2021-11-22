package internal

type Balance struct {
	Id         int `db:"id"`
	CustomerId int `db:"customer_id" binding:"required"`
	Amount     int `db:"amount" binding:"required"`
}
