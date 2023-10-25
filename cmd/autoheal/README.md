# Why?

Sometimes services propagated via skupper are not reachable. Often the issue 
can be solved by deleting and recreating the service.

# Assumptions

Cluster where services are deployed: Cluster A
Cluster where services are propagated via skupper: Cluster B

1. Resources will be automatically recreated on deletion (having ArgoCD, ... managing the cluster.).
2. Prometheus and blackbox exporter are running in cluster B.
3. Prometheus uptime checks are added for skupper services.
4. Prometheus metrics are readable from cluster A.
5. The services which are propagated to different cluster have `-skupper` suffix. Ideally add new services!


# What?

This skript should run as pod in a cluster A.

If a service is found to be disconnected, the service will be deleted.
