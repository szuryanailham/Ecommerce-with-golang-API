package products

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/szuryanailham/ecom/internal/adapters/sqlc"
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

func (h *handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
    var tempProduct createProductRequest
	if err := json.Read(r,&tempProduct); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdProduct, err := h.service.CreateProduct(r.Context(), repo.CreateProductParams{
		Name: tempProduct.Name,
		PriceCents: tempProduct.PriceCents,
		Quantity: tempProduct.Quantity,
	})
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create product", http.StatusInternalServerError)
	}
	json.Write(w, http.StatusCreated, createdProduct)
}

func ( h*handler)UpdateProductName(w http.ResponseWriter, r *http.Request){
	  idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "missing product id", http.StatusBadRequest)
        return
    }
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid product id", http.StatusBadRequest)
        return
    }
	var tempUpdateProduct UpdateProductNameRequest
	if err :=  json.Read(r, &tempUpdateProduct); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
    err = h.service.UpdateProductName(r.Context(), repo.UpdateProductNameParams{
        ID:   id,
        Name: tempUpdateProduct.Name,
    })
	 if err != nil {
		log.Println(err)
    	http.Error(w, "Failed to update product", http.StatusInternalServerError)
    	return
	 }
	json.Write(w, http.StatusAccepted, map[string]string{
    "message": "Product updated successfully",
})

}


