# Kubernetes jobs

## sitesync-job

Run skupper-sitesync in cluster.

```BASH
kubectl apply -n $NAMESPACE -f sitesync-job.yaml
```

> Prequal: Deploy service account skupper-sitesync

```YAML
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: skupper-sitesync
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: skupper-sitesync
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: skupper-sitesync
subjects:
- kind: ServiceAccount
  name: skupper-sitesync
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: skupper-sitesync
```
