package lnurlp

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

type Env struct {
	env string
	url string
	domain string
	name string
	host string
	macaroon string
}

// This is a one field type just to make sure we remember we're using MSats
type InvoiceRequest struct {
    MSats int64 
}

func getEnv() (Env, error) {
	env := Env{
		url:    os.Getenv("VERCEL_URL"),
		env:    os.Getenv("VERCEL_ENV"),
		domain: os.Getenv("ADDRESSLESS_DOMAIN"),
		name:   os.Getenv("ADDRESSLESS_NAME"),
		host:   os.Getenv("LND_HOST"),
		macaroon: os.Getenv("LND_MACAROON"),
	}

	if (env.domain == "" || env.name == "" || env.host == "" || env.macaroon == "") {
		return env, errors.New("Something is configured wrong. Maybe double check your env?")
	};

	return env, nil
}

func makeMetadata(env Env) string {
	var domain string

	if env.domain == "" || env.env != "production" {
		domain = env.url 
	} else {
		domain = env.domain 
	}

	metadata, _ := sjson.Set("[]", "0.0", "text/identifier")
	metadata, _ = sjson.Set(metadata, "0.1", env.name+"@"+domain)

	metadata, _ = sjson.Set(metadata, "1.0", "text/plain")
	metadata, _ = sjson.Set(metadata, "1.1", "Satoshis to "+env.name+"@"+domain)

	return metadata 
}



func getUrl(env Env) string {
	if env.url == "" || env.env == "" {
		return "http://localhost:3000" 
	} else {
		return "https://" + env.url 
	}
}


func makeInvoice(inv InvoiceRequest, env Env) (string, error) {
	// make the lnurlpay description_hash
	description_hash := sha256.Sum256([]byte(makeMetadata(env)))

	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi: inv.MSats,
		// TODO makeinvoice supports other backends than LND we should too
		Backend: makeinvoice.LNDParams{
			Host: env.host,
			Macaroon: env.macaroon,
		},
		DescriptionHash: description_hash[:],
		Label: strconv.FormatInt(time.Now().Unix(), 16),
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got lnurl request\n")

	env, err := getEnv() 

	// fail fast if the env isn't configured correctly
	if err != nil {
		json.NewEncoder(w).Encode(lnurl.ErrorResponse(err.Error()))
		return
	}

	if amount := r.URL.Query().Get("amount"); amount == "" {
		if err != nil {
			json.NewEncoder(w).Encode(lnurl.ErrorResponse(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(lnurl.LNURLPayResponse1{
			LNURLResponse:   lnurl.LNURLResponse{Status: "OK"},
			Callback:        fmt.Sprintf("%s/api/lnurlp", getUrl(env)),
			MinSendable:     1000,
			MaxSendable:     100000000,
			EncodedMetadata: makeMetadata(env),
			CommentAllowed:  0,
			Tag:             "payRequest",
		})
	} else {
		msat, err := strconv.Atoi(amount)
		if err != nil {
			json.NewEncoder(w).Encode(lnurl.ErrorResponse("amount is not integer"))
			return
		}

		bolt11, err := makeInvoice(InvoiceRequest{MSats: int64(msat)}, env)

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