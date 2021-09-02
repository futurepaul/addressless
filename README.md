# Addressless

The goal is to put a [Lightning Address](https://lightningaddress.com/) "server" on a serverless host like Vercel. So far this just creates invoices, using [makeinvoice](https://pkg.go.dev/github.com/fiatjaf/makeinvoice), but the goal is to replicate the functionality of [satdress](https://github.com/fiatjaf/satdress), minus the federated part.

The workflow should be: click the "Deploy" button, follow the instructions over on Vercel (it should prompt you to fill in the necessary environment variables with your LND host and macaroon), and (optional) point a domain. You should end up with a splash page, on a url you control, that you can share with people who want to pay you. And in the background Vercel should be able to handle the LNURL stuff that Lightning Address needs.

## Testing Locally

- clone this repo
- run `npm install`
- install the [Vercel CLI](https://vercel.com/cli) so you can run the serverless go api
- run `vercel dev` and follow its instructions
- to add the necessary environment variables `LND_HOST` and `LND_MACAROON` run `vercel env add` for each and follow the prompt
- `vercel env pull` will put those environment variables in your `.env`
- `vercel dev` should actually be working now

I'd feel better if you only used a testnet node and def be sure to use your invoice macaroon.

## TODO:

- [x] get this actually working on vercel: https://addressless.vercel.app/
- [ ] draw the rest of the owl (implement Lightning Address)
- [ ] make a tutorial for how to do this
- [ ] default page could be the tutorial / demo of the functionality
