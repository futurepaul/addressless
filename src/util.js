import { bech32 } from "bech32";

// from https://github.com/fiatjaf/filemarket/blob/master/client/helpers.js
export function lnurlencode(url) {
  console.log(url);
  return bech32.encode(
    "lnurl",
    bech32.toWords(url.split("").map((c) => c.charCodeAt(0))),
    1500
  );
}
