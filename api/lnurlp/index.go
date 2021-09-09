package lnurl

// I'm doing separate folders for each "package" because vercel gets mad about sibling files

import (
	"crypto/sha256"
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
	var vercel_env = os.Getenv("VERCEL_ENV")
	var vercel_domain = os.Getenv("VERCEL_URL")

	var pretty_domain = os.Getenv("ADDRESSLESS_DOMAIN")
	var name = os.Getenv("ADDRESSLESS_NAME")

	var domain string

	if pretty_domain == "" || vercel_env != "production" {
		domain = vercel_domain 
	} else {
		domain = pretty_domain
	}
	
	metadata, _ := sjson.Set("[]", "0.0", "text/identifier")
	metadata, _ = sjson.Set(metadata, "0.1", name+"@"+domain)

	metadata, _ = sjson.Set(metadata, "1.0", "text/plain")
	metadata, _ = sjson.Set(metadata, "1.1", "Satoshis to "+name+"@"+domain)

	return metadata
}

type InvoiceRequest struct {
    MSats int64 
}

func getUrl() string {
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

	// make the lnurlpay description_hash
	description_hash := sha256.Sum256([]byte(makeMetadata()))



	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi: inv.MSats,
		Backend: makeinvoice.LNDParams{
			Host: host,
			Macaroon: macaroon,
		},
		DescriptionHash: description_hash[:],
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
			Callback:        fmt.Sprintf("%s/api/lnurlp", getUrl()),
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

		bolt11, err := makeInvoice(InvoiceRequest{MSats: int64(msat)})
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