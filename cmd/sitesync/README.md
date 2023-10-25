# What is this?

This skript periodically checks the `resourceVersion` of `skupper-site` config map. If a new version 
is deployed, skupper is "restarted" so that resources in the cluster reflect the definitions 
in the config map.

# Deploy skupper via YAML

When skupper is managed via git ops and thus, deployed via yaml files, 
some manual work is needed when the site config is updated.

Changes in the site config might lead to restarting the skupper-router and 
skupper-service-controller, but several changes will not be applied.

Examples are:
- replicas
- resource requests/limits

What is needed to apply the changes:
- delete skupper-router deployment
- delete skupper-service-controller deployment
- restart skupper-site-controller deployment
