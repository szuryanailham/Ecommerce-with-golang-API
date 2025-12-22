package orders

import (
	"log"
	"net/http"

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

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request){
	var tempOrder createOrderParams
	if err := json.Read(r,&tempOrder); err != nil {
		log.Println(err)
		http.Error(w,err.Error(), http.StatusBadRequest)
		return 
	}
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w,err.Error(), http.StatusInternalServerError)
	}
	json.Write(w, http.StatusCreated, createdOrder)
}