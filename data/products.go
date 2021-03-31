package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "ess123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          3,
		Name:        "Americano",
		Description: "Korean coffee",
		Price:       1.5,
		SKU:         "kco123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}

func (p *Product) Validate() error {
	v := validator.New()
	v.RegisterValidation("sku", ValidateSKU)
	return v.Struct(p)
}

func ValidateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func AddProducts(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	idx, err := searchProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	productList[idx] = p

	return nil
}

func DeleteProduct(id int) error {
	idx, err := searchProduct(id)

	if err != nil {
		return err
	}

	productList = append(productList[:idx], productList[idx+1:]...)
	return nil
}

var errorNoProd = fmt.Errorf("Product Not Found")

func searchProduct(id int) (int, error) {
	for x, prod := range productList {
		if id == prod.ID {
			return x, nil
		}
	}

	return -1, errorNoProd
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func View_ProductList() {
	for x, product := range productList {
		fmt.Println(x, *product)
	}
}

func GetProducts() Products {
	return productList
}
