kube-onions
===========

Sprinkle some onions on your kubernetes clusters.

Install
-------

```
$ helm install chart/onion-controller/ --name onion-controller --namespace kube-system --wait
```

Use
---

Get your private key into a secret:

    kubectl create secret generic bmy7nlgozpyn26tv --from-file=private_key

Create an onion service:

```yaml
apiVersion: onion.kragniz.eu/v1alpha1
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
