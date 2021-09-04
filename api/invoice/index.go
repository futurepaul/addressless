package invoice

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fiatjaf/makeinvoice"
)

type InvoiceRequest struct {
	Description string
    Sats int64 
}

func makeInvoice(inv InvoiceRequest) (string, error) {
	host:= os.Getenv("LND_HOST")
	macaroon:= os.Getenv("LND_MACAROON")

	if host == "" || macaroon == "" {
		return "", errors.New("Something is wrong with the credentials.") 
	}

	msats := inv.Sats * 1000
	desc := inv.Description

	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi: msats,
		Backend: makeinvoice.LNDParams{
			Host: host,
			Macaroon: macaroon,
		},
		Description: desc,
		// TODO: what should this label be?
		Label: strconv.FormatInt(time.Now().Unix(), 16),
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ir InvoiceRequest 
	err := decoder.Decode(&ir)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bolt11, err := makeInvoice(ir);

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(bolt11)
}