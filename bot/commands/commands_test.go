package commands

import (
	"github.com/stretchr/testify/assert"
	"money-diff/model"
	"testing"
)

func TestCalculateCombinedTwoUsersBToA(t *testing.T) {
	payments := []model.GroupedPayment{
		{"a", 10},
		{"b", 30},
	}
	directPayments := []model.GroupedDirectPayment{
		{model.DirectUserDTO{WhoOwes: "a", Whom: "b"}, 50},
	}
	res, _ := Calculate(payments, directPayments)
	assert.Equal(t, res, map[string]debt{
		"ab": {"a", "b", 60},
	})
}

func TestCalculateCombinedTwoUsersAToB(t *testing.T) {
	payments := []model.GroupedPayment{
		{"a", 10},
		{"b", 5},
	}
	directPayments := []model.GroupedDirectPayment{
		{model.DirectUserDTO{WhoOwes: "b", Whom: "a"}, 30},
	}
	res, _ := Calculate(payments, directPayments)
	assert.Equal(t, res, map[string]debt{
		"ba": {"b", "a", 32.5},
	})
}

func TestCalculateCombinedThreeUsers(t *testing.T) {
	payments := []model.GroupedPayment{
		{"a", 200},
		{"b", 10},
	}
	directPayments := []model.GroupedDirectPayment{
		{model.DirectUserDTO{WhoOwes: "c", Whom: "a"}, 15},
	}
	res, _ := Calculate(payments, directPayments)
	assert.Equal(t, res, map[string]debt{
		"ba": {"b", "a", 95},
		"ca": {"c", "a", 15},
	})
}

func TestCalculateCombinedManyUsers(t *testing.T) {
	payments := []model.GroupedPayment{
		{"a", 200},
		{"b", 10},
		{"c", 30},
	}
	directPayments := []model.GroupedDirectPayment{
		{model.DirectUserDTO{WhoOwes: "c", Whom: "a"}, 100},
		{model.DirectUserDTO{WhoOwes: "a", Whom: "b"}, 50},
	}
	res, _ := Calculate(payments, directPayments)
	assert.Equal(t, res, map[string]debt{
		"ba": {"b", "a", 20},
		"ca": {"c", "a", 150},
	})
}

func TestCalculateCombinedManyUsers2(t *testing.T) {
	payments := []model.GroupedPayment{
		{"a", 200},
		{"b", 10},
		{"c", 30},
	}
	directPayments := []model.GroupedDirectPayment{
		{model.DirectUserDTO{WhoOwes: "c", Whom: "a"}, 100},
		{model.DirectUserDTO{WhoOwes: "a", Whom: "b"}, 50},
		{model.DirectUserDTO{WhoOwes: "b", Whom: "a"}, 90},
		{model.DirectUserDTO{WhoOwes: "c", Whom: "b"}, 14},
	}
	res, _ := Calculate(payments, directPayments)
	assert.Equal(t, res, map[string]debt{
		"ba": {"b", "a", 110},
		"ca": {"c", "a", 150},
		"cb": {"c", "b", 14},
	})
}
