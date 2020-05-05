# Secrets folder

## Gitignored contents:
* privkey.pem
* cert.pem

## Generate an RSA keypair with a 2048 bit private key
`openssl genpkey -algorithm RSA -out privkey.pem -pkeyopt rsa_keygen_bits:2048`

## Extract the public key from an RSA keypair
`openssl rsa -pubout -in privkey.pem -out cert.pem`
