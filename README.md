<p align="center">
  <img height="300" src="https://sr.ht/2mc0.png">
</p>

<h1 align="center">tor-controller</h1>

[![Build Status](https://img.shields.io/travis-ci/kragniz/tor-controller.svg?style=flat-square)](https://travis-ci.org/kragniz/tor-controller)

Sprinkle some onions on your kubernetes clusters.

tor-controller allows you to create `OnionService` resources in kubernetes.
These services are used similarly to standard kubernetes services, but they serve
traffic on the tor network.

tor-controller creates the following resources for each OnionService:

- a service, which is used to send traffic to application pods
- tor pod, which contains a tor daemon to serve incoming traffic from the tor
  network, and a management process that watches the kubernetes API and
  generates tor config, signaling the tor daemon when it changes
- rbac rules

<p align="center">
  <img src="https://sr.ht/6WbX.png">
</p>

Install
-------

Install tor-controller:

    $ kubectl apply -f hack/install.yaml

Quickstart with random address
------------------------------

Create an onion service, `onionservice.yaml`:

```yaml
apiVersion: tor.k8s.io/v1alpha1
kind: OnionService
metadata:
  name: basic-onion-service
spec:
  version: 2
  selector:
    app: example
  ports:
  - publicPort: 80
    targetPort: 80
```

Apply it:

    $ kubectl apply -f onionservice.yaml

View it:

```bash
$ kubectl get onionservices -o=custom-columns=NAME:.metadata.name,HOSTNAME:.status.hostname
NAME                    HOSTNAME
basic-onion-service     h7px2yyugjqkztrb.onion
```

Exposing a deployment with a fixed address
------------------------------------------

Create some deployment to test against, in this example we'll deploy an echoserver. Create `echoserver.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: http-app
  template:
    metadata:
      labels:
        app: http-app
    spec:
      containers:
      - name: http-app
        image: gcr.io/google_containers/echoserver:1.8
        ports:
        - containerPort: 8080
```
Apply it:

    $ kubectl apply -f echoserver.yaml

For a fixed address, we need a private key. This should be kept safe, since
someone can impersonate your onion service if it is leaked.
Generate an RSA private key (only valid for v2 onion services, v3 services use Ed25519 instead):

    $ openssl genrsa -out private_key 1024

Put your private key into a secret:

    $ kubectl create secret generic example-onion-key --from-file=private_key

Create an onion service, `onionservice.yaml`, referencing the private key we just created:

```yaml
apiVersion: tor.k8s.io/v1alpha1
kind: OnionService
metadata:
  name: example-onion-service
spec:
  version: 2
  selector:
    app: http-app
  ports:
    - targetPort: 8080
      publicPort: 80
  privateKeySecret:
    name: example-onion-key
    key: private_key
```

Apply it:

    $ kubectl apply -f onionservice.yaml

List active OnionServices:

```
$ kubectl get onionservices -o=custom-columns=NAME:.metadata.name,HOSTNAME:.status.hostname
NAME                    HOSTNAME
example-onion-service   s2c6qry5bj57vyms.onion
```

This service should now be accessable from any tor client,
for example [Tor Browser](https://www.torproject.org/projects/torbrowser.html.en):

<p align="center">
  <img src="https://sr.ht/FLbP.png">
</p>

Random service names
--------------------

If `spec.privateKeySecret` is not specified, tor-controller will start a service with a random name.
This will remain in use until the tor-daemon pod restarts or is terminated for some other reason.

Onion service versions
----------------------

The `spec.version` field specifies which onion protocol to use.
v2 is the classic and well supported, v3 is the new replacement.

The biggest difference from a user's point of view is the length of addresses. v2
service names are short, like `x3yvl2svtqgzhcyz.onion`. v3 are longer, like
`ljgpby5ba3xi5osslpdvqsumdb4sbclb2amxtm6a3cwnq7w7sj72noid.onion`.

tor-controller defaults to using v3 if `spec.version` is not specified.


Using with nginx-ingress
------------------------

tor-controller on its own simply directs TCP traffic to a backend service.
If you want to serve HTTP stuff, you'll probably want to pair it with
nginx-ingress or some other ingress controller.

To do this, first install nginx-ingress normally. Then point an onion service
at the nginx-ingress-controller, for example:

```yaml
apiVersion: tor.k8s.io/v1alpha1
kind: OnionService
metadata:
  name: nginx-onion-service
spec:
  version: 2
  selector:
    app: nginx-ingress-controller
    name: nginx-ingress-controller
  ports:
  - publicPort: 80
    targetPort: 80
    name: http
  privateKeySecret:
    name: nginx-onion-key
    key: private_key
```

This can then be used in the same way any other ingress is. Here's a full
example, with a default backend and a subdomain:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: http-app
  template:
    metadata:
      labels:
        app: http-app
    spec:
      containers:
      - name: http-app
        image: gcr.io/google_containers/echoserver:1.8
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: http-app
  labels:
    app: http-app
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: http-app
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: http-app
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  backend:
    serviceName: default-http-backend
    servicePort: 80
  rules:
  - host: echoserver.h7px3yyugjqkztrb.onion
    http:
      paths:
      - path: /
        backend:
          serviceName: http-app
          servicePort: 8080
```
