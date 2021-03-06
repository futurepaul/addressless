# Addressless

Put a [Lightning Address](https://lightningaddress.com/) "server" on a serverless host like Vercel. The code is heavily borrowed from [satdress](https://github.com/fiatjaf/satdress), minus the federated part.

The workflow should be: click the "Deploy" button, follow the instructions over on Vercel (it should prompt you to fill in the necessary environment variables with your LND host and macaroon), and (optional) point a domain. You should end up with a splash page, on a url you control, that you can share with people who want to pay you. And in the background Vercel handles the LNURL stuff that Lightning Address needs.

## Deploying to Vercel

Click this deploy button:

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Ffuturepaul%2Faddressless&env=ADDRESSLESS_DOMAIN,ADDRESSLESS_NAME,LND_HOST,LND_MACAROON&envDescription=Env%20vars%20are%20described%20in%20the%20readme&envLink=https%3A%2F%2Fgithub.com%2Ffuturepaul%2Faddressless%2Fblob%2Fmaster%2FREADME.md%23deploying-to-vercel&demo-title=Addressless%20Demo&demo-description=What%20this%20deploy%20looks%20like&demo-url=https%3A%2F%2Faddressless.vercel.app%2F)

Follow the instructions.

When it asks you for your environment variables, here's what it needs:

- `ADDRESSLESS_DOMAIN` = the domain part of your Lightning Address (my_name@**example.com**... if you don't have a domain you may be able to use `your-project-name.vercel.app` if it's not taken)
- `ADDRESSLESS_NAME` = the name part (**my_name**@example.com)
- `LND_HOST` = the url and REST port your LND node can be found at (https://my-node-123.voltageapp.io:8080) RIGHT NOW TOR DOESN'T WORK [WHICH SUCKS SUGGESTIONS WELCOME](https://github.com/futurepaul/addressless/issues/1) 
- `LND_MACAROON` = your node's invoice macaroon (abc123...)

If you need to change any of these environment variables later, go to your project on Vercel -> Settings -> Environment Variables. You can also go there to customize different variables for different environments (production, preview, development)

To use a custom domain name go to Settings -> Domains and add your desired domain.

If you have a non-Vercel website and don't want to replace it I think you should be fine with a couple redirects:
- `your-domain.com/.well-known/lnurlp/<name>` -> `your-vercel-address.com/.well-known/lnurlp/<name>`
- `your-domain.com/api/lnurlp` -> `your-vercel-address.com/api/lnurlp`

## Testing Locally

- clone this repo
- run `npm install`
- install the [Vercel CLI](https://vercel.com/cli) so you can run the serverless go api
- run `vercel dev` and follow its instructions
- to add the necessary environment variables `LND_HOST` and `LND_MACAROON` run `vercel env add` for each and follow the prompt
- `vercel env pull` will put those environment variables in your `.env`
- `source .env`
- `vercel dev` should actually be working now

I'd feel better if you only used a testnet node and def be sure to use your invoice macaroon.

## TODO:

- [x] get this actually working on vercel: https://addressless.vercel.app/
- [x] draw the rest of the owl (implement Lightning Address)
- [x] make a tutorial for how to do this
- [x] default page is a demo of the functionality
- [ ] support tor
- [ ] tutorial for deploying 
- [ ] support multiple names as aliases
