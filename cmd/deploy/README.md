# Deploy skupper via YAML

When skupper is managed via git ops and thus, deployed via yaml files, some manual work is needed:

* apply manifests
* create the default secret by applying token-request secret
* get resulting default secret
(* seal secret)
* add secret to git
(* delte original secret from cluster, if secret was sealed)
(* apply sealed secret to cluster)

If secret is sealed via sealed-secrets there are some more steps:

* apply manifests
* create the default secret by applying token-request secret
* get resulting default secret
* seal secret
* add secret to git
* delte original secret from cluster
* apply sealed secret to cluster

> Some of these steps require pre-conditions.

### Token request

Before applying the token request, wait for the deployments `skupper-router` and 
`skupper-service-controller` to be ready. If they are not, the `skupper-site-controller` 
will throw an error and the secret will not contain any data.
