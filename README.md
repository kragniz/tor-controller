<p align="center">
  <img height="300" src="https://sr.ht/2mc0.png">
</p>

tor-controller
==============

Sprinkle some onions on your kubernetes clusters.

Use
---

Get your private key into a secret:

    kubectl create secret generic bmy7nlgozpyn26tv --from-file=private_key

Create an onion service:

```yaml
apiVersion: tor.k8s.io/v1alpha1
kind: OnionService
metadata:
  name: example-onion
spec:
  selector:
    app: httpd
  ports:
    - targetPort: 8080
      publicPort: 80
  privateKeySecret:
    name: bmy7nlgozpyn26tv
    key: private_key
```
