# Chirpy



Create a secret for your server and store it in your .env file. This is the secret used to sign and verify JWTs. By keeping it safe, no other servers will be able to create valid JWTs for your server. We will yet again use environment variables. You can generate a nice long random string on the command line like this:

openssl rand -base64 64

Secrets should NOT be stored in Git, just in case anyone malicious gains access to your repository.
