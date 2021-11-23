package model

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

type RefillDto struct {
	CustomerId int         `json:"customer_id" binding:"required,gt=0"`
	Amount     PennyAmount `json:"amount" binding:"required,gt=0"`
}

type PennyAmount struct {
	Int int `json:"amount" binding:"required,gt=0"`
}

func (ci *PennyAmount) UnmarshalJSON(data []byte) error {
	float, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return errors.New("PennyAmount: UnmarshalJSON: " + err.Error())
	}
	float *= 100
	bytes := []byte(strconv.Itoa(int(math.Round(float))))
	err = json.Unmarshal(bytes, &ci.Int)
	if err != nil {
		return errors.New("PennyAmount: UnmarshalJSON: " + err.Error())
	}
	return nil
}

func (ci *PennyAmount) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(float64(ci.Int) / 100)
	return bytes, err
}
