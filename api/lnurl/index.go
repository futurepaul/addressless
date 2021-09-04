package lnurl

// I'm doing separate folders for each "package" because vercel gets mad about sibling files

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fiatjaf/go-lnurl"
	"github.com/fiatjaf/makeinvoice"
	"github.com/tidwall/sjson"
)

func makeMetadata() string {
	var domain = os.Getenv("VERCEL_URL")
	
	metadata, _ := sjson.Set("[]", "0.0", "text/identifier")
	metadata, _ = sjson.Set(metadata, "0.1", "paul@"+domain)

	metadata, _ = sjson.Set(metadata, "1.0", "text/plain")
	metadata, _ = sjson.Set(metadata, "1.1", "Satoshis to paul@"+domain)

	return metadata
}

type InvoiceRequest struct {
	Description string
    MSats int64 
}

func getUrl() string {
	// TODO make this less Vercel specific 
	vercel_url := os.Getenv("VERCEL_URL")
	vercel_env := os.Getenv("VERCEL_ENV")
	

	if vercel_url == "" || vercel_env == "" {
		return "http://localhost:3000"
	} else {
		return "https://" + vercel_url
	}

}


func makeInvoice(inv InvoiceRequest) (string, error) {
	host:= os.Getenv("LND_HOST")
	macaroon:= os.Getenv("LND_MACAROON")

	if host == "" || macaroon == "" {
		return "", errors.New("Something is wrong with the credentials.") 
	}

	desc := inv.Description

	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi: inv.MSats,
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
	fmt.Printf("got lnurl request\n")

	if amount := r.URL.Query().Get("amount"); amount == "" {
		// check if the receiver accepts comments
		var commentLength int64 = 0
		// TODO: support webhook comments

		json.NewEncoder(w).Encode(lnurl.LNURLPayResponse1{
			LNURLResponse:   lnurl.LNURLResponse{Status: "OK"},
			Callback:        fmt.Sprintf("%s/.well-known/lnurlp/%s", getUrl(), "paul"),
			MinSendable:     1000,
			MaxSendable:     100000000,
			EncodedMetadata: makeMetadata(),
			CommentAllowed:  commentLength,
			Tag:             "payRequest",
		})

	} else {
		msat, err := strconv.Atoi(amount)
		if err != nil {
			json.NewEncoder(w).Encode(lnurl.ErrorResponse("amount is not integer"))
			return
		}

		bolt11, err := makeInvoice(InvoiceRequest{MSats: int64(msat), Description: "Testing 123"})
		if err != nil {
			json.NewEncoder(w).Encode(
				lnurl.ErrorResponse("failed to create invoice: " + err.Error()))
			return
		}

		json.NewEncoder(w).Encode(lnurl.LNURLPayResponse2{
			LNURLResponse: lnurl.LNURLResponse{Status: "OK"},
			PR:            bolt11,
			Routes:        make([][]lnurl.RouteInfo, 0),
			Disposable:    lnurl.FALSE,
			SuccessAction: lnurl.Action("Payment received!", ""),
		})
	}
}