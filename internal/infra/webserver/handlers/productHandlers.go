package handlers

import (
	"GoAPI/internal/dto"
	"GoAPI/internal/entity"
	"GoAPI/internal/infra/database"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	ProductDB database.ProductI
}

func NewProductHandler(db database.ProductI) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	ErrorHand(w, err, http.StatusBadRequest)

	p, err := entity.NewProduct(product.Name, product.Price)
	ErrorHand(w, err, http.StatusBadRequest)

	err = h.ProductDB.Create(p)
	ErrorHand(w, err, http.StatusInternalServerError)

	w.WriteHeader(http.StatusCreated)
}

func ErrorHand(w http.ResponseWriter, err error, status int) {
	if err != nil {
		w.WriteHeader(status)
		return
	}
}
