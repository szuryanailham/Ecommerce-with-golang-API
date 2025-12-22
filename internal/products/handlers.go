package products

import (
	"log"
	"net/http"
	"strconv"

	"github.com/szuryanailham/ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter,r *http.Request) {
	 Products, err := h.service.ListProducts(r.Context())
	 if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	 }
	json.Write(w, http.StatusOK, Products)
}

func (h*handler) FindProductByID(w http.ResponseWriter, r *http.Request){
	idStr := r.URL.Query().Get("id")
	if idStr == ""{
		http.Error(w, "missing product id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
	}
	id := int32(idInt)

	product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
	}
	json.Write(w, http.StatusOK, product)
}