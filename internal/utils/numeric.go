package utils

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

// NumericToString converts pgtype.Numeric to string
func NumericToString(n pgtype.Numeric) string {
	if !n.Valid {
		return "0.00"
	}

	// pgtype.Numeric uses a big.Int and an exponent
	// value = Int * 10^Exp

	f := new(big.Float).SetInt(n.Int)
	if n.Exp < 0 {
		// Negative exponent means division
		divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil))
		f.Quo(f, divisor)
	} else if n.Exp > 0 {
		// Positive exponent means multiplication
		multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil))
		f.Mul(f, multiplier)
	}

	return f.Text('f', 2)
}

// ToNumeric converts float64 to pgtype.Numeric
func ToNumeric(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	err := n.Scan(fmt.Sprintf("%f", f))
	if err != nil {
		n.Valid = false
		return n
	}
	n.Valid = true
	return n
}

// NumericToFloat converts pgtype.Numeric to float64
func NumericToFloat(n pgtype.Numeric) (float64, error) {
	s := NumericToString(n)
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
