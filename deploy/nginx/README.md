# Dedicated-server reverse proxy

Production exposes the CMS application through Kubernetes NodePort `31373`.

## DNS

In one.com, create an `A` record for `cms.urpi.be` pointing to the dedicated
server's public IPv4 address. Add an `AAAA` record only when the server and its
firewall are configured for public IPv6. Wait for:

```bash
dig +short cms.urpi.be
```

to return the dedicated server address before requesting a certificate.

## Nginx and Certbot

The included virtual host assumes Nginx runs on the Kubernetes node and can
reach the NodePort at `127.0.0.1:31373`.

```bash
sudo install -m 0644 deploy/nginx/cms.urpi.be.conf \
  /etc/nginx/sites-available/cms.urpi.be.conf
sudo ln -s /etc/nginx/sites-available/cms.urpi.be.conf \
  /etc/nginx/sites-enabled/cms.urpi.be.conf
sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx -d cms.urpi.be
sudo certbot renew --dry-run
```

Certbot adds the HTTPS listener, certificate paths, and HTTP-to-HTTPS redirect
to this virtual host. Allow inbound TCP ports `80` and `443`. Do not expose
NodePort `31373` to the public internet. If Nginx is not on a Kubernetes node, replace
`127.0.0.1` in `proxy_pass` with a private cluster-node address and allow
`31373` only from the Nginx server.

## Deploy and verify

Deploy the production profile:

```bash
make helm-deploy-production
kubectl -n cms-production get pods
kubectl -n cms-production get service cms
curl --fail http://127.0.0.1:31373/health
curl --fail https://cms.urpi.be/health
```

The `cms` service should show `80:31373/TCP`.
