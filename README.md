# Addressless

The goal is to put a [Lightning Address](https://lightningaddress.com/) "server" on a serverless host like Vercel. So far this just creates invoices, using [makeinvoice](https://pkg.go.dev/github.com/fiatjaf/makeinvoice), but the goal is to replicate the functionality of [satdress](https://github.com/fiatjaf/satdress), minus the federated part.

The workflow should be: click the "Deploy" button, follow the instructions over on Vercel (it should prompt you to fill in the necessary environment variables with your LND host and macaroon), and (optional) point a domain. You should end up with a splash page, on a url you control, that you can share with people who want to pay you. And in the background Vercel should be able to handle the LNURL stuff that Lightning Address needs.
