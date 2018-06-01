<p align="center">
  <img height="300" src="https://sr.ht/2mc0.png">
</p>

<h1 align="center">tor-controller</h1>

[![Build Status](https://img.shields.io/travis-ci/kragniz/tor-controller.svg?style=flat-square)](https://travis-ci.org/kragniz/tor-controller)

Sprinkle some onions on your kubernetes clusters.

Quickstart
----------

Install tor-controller:

    $ kubectl apply -f hack/install.yaml

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

Generate a test private key:

    $ openssl genrsa -out private_key 1024

Put your private key into a secret:

    $ kubectl create secret generic example-onion-key --from-file=private_key

Create an onion service, `onionservice.yaml`:

```yaml
apiVersion: tor.k8s.io/v1alpha1
kind: OnionService
metadata:
  name: example-onion-service
spec:
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
example-onion-service   h7px2yyugjqkztrb.onion
```

tor-controller creates the following resources for each OnionService:

- a service, which is used to send traffic to application pods
- a configmap containing tor configuration pointing at the service
- tor daemon pod, which serves incoming traffic from the tor network

<p align="center">
  <img src="https://sr.ht/6WbX.png">
</p>

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
  ports:
  - publicPort: 80
    targetPort: 80
    name: http
  selector:
    app: nginx-ingress-controller
    name: nginx-ingress-controller
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
