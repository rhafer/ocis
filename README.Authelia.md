# Embedded Authelia

This branch contains a proof-of-concept for embedding the Authelia
Authentication Server as the OIDC Identity provider into ocis.

The integration is still very rough and only meant as a proof of concept. It
is my no means production ready. For now it depends on a fork of authelia
which moves some bits of authelia `internal` package to a publically
comsumable package.

# Preparations

## Authelia

* Checkout the `ocis-embed` branch of https://github.com/rhafer/authelia
* Prepare the integrated react web portal:
  ```
  cd web
  pnpm install
  pnpm build
  cd ..
  ```
* Move some required files into place:
  ```
  cp -a api internal/server/public_html/
  ```

## OCIS

Warning: The current implementation still contains a hardcoded rsa key pair
and some other hardcoded secrets. It's by no means intended for production

* Checkout the `authelia-experiment` branch of https://github.com/rhafer/ocis.git
* Create a go workspace using the locally checked out and prepared `authelia`
  tree:
  ```
  go work init ../ocis
  go work use <path-to-your-authelia-checkout>
  ```
  We still need this step for the integrated react webapp of authelia to work
* Build ocis
  ```
  rm ocis/bin/ocis
  make -C ocis build
  ```
* Create an empty configuration file for authelia. (Inside the directory from
  where you'll be running ocis. This is needed to make authelia happy, as it
  checks the file for existence.
  ```
  touch configuration.yml
  ```
* Start OCIS
  ```
  ocis/bin/ocis init --insecure true
  OCIS_OIDC_ISSUER=https://localhost.localdomain:9200/authelia \
  OCIS_URL=https://localhost.localdomain:9200 \
  PROXY_OIDC_ACCESS_TOKEN_VERIFY_METHOD=none \
  ocis/bin/ocis server
  ```
* Open the OCIS web site:
  https://localhost.localdomain:9200

### What is working
Logging in via the Web App an using ocis

### What is NOT working
* running ocis behind a reverse proxy (traefik, ...)
  Authelia relies on the correct X-Forwared-* headers to be correctly set,
  when running ocis behind a reverse proxy, the ocis proxy seems to overwrite
  the X-Forwarded-* header and authelias .well-known/openid-configuration
  respose will point wrong URLs. (The exact problem here still needs to be
  investigated)
* Currently we only generate the OIDC client configuration for `web`. So all
  the other clients (mobile, desktop, ...) will currently not work.
* The `idm` server does not support "RootDSE" queries, you'll see some scary
  warning from authelia at startup because of that.
* logout is still borked
* password resets
* using the default `OCIS_URL=https://localhost:9200`. Authelia
  requires the domain name to have at least one `.` in it.
* Probably a lot of other things ...

