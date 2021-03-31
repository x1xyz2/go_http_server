package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Fucker",
		Price: 1.23,
		SKU:   "aaa-bbbccc",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
