# 🛠️ SRE Survival Playbook

## 📁 Table of Contents

1. [Mental Models](#mental-models)
2. [Kubernetes Core Operations](#kubernetes-core-operations)
3. [Docker & Containerization](#docker--containerization)
4. [Local Cluster Bootstrapping (K3D)](#local-cluster-bootstrapping-k3d)
5. [Helm](#helm)
6. [RBAC & Permissions](#rbac--permissions)
7. [Ingress & Port Forwarding](#ingress--port-forwarding)
8. [Troubleshooting & Gotchas](#troubleshooting--gotchas)
9. [Git Survival Commands](#git-survival-commands)
10. [SRE Wisdom & Mindset](#sre-wisdom--mindset)

---

## Mental Models

* **Your Terminal ≠ The Cluster**

  * The cluster lives in Docker (via `k3d`), not your terminal. Killing the shell doesn't stop Kubernetes.

* **Kubernetes is a Control Plane**

  * You describe desired state via YAML (or Helm), and the controllers try to make it real.

* **K8s is Not Persistent Unless You Make It So**

  * Use volumes for data. Otherwise, pods and configs can die on reboot.

---

## Kubernetes Core Operations

```bash
kubectl get nodes
kubectl get pods -A
kubectl apply -f deployment.yaml
kubectl delete -f deployment.yaml
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl exec -it <pod-name> -- /bin/sh
```

> 🧠 Cheat: Use `k9s` for a fast TUI to see all pods, namespaces, and logs.

---

## Docker & Containerization

* Build:

```bash
docker build -t my-image:0.1 .
```

* Run locally:

```bash
docker run -p 8080:80 my-image:0.1
```

* Inspect image:

```bash
docker inspect my-image:0.1
```

---

## Local Cluster Bootstrapping (K3D)

```bash
k3d cluster create mycluster --agents 2
k3d kubeconfig get mycluster > ~/.kube/config
kubectl cluster-info
```

> 🧠 K3D runs the entire Kubernetes system inside Docker containers.

---

## Helm

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx ingress-nginx/ingress-nginx
```

> 🧠 Helm is like a package manager (like `apt` or `brew`) but for Kubernetes YAML.

---

## RBAC & Permissions

* You need to create `ServiceAccounts`, `Roles`, and `RoleBindings` to grant access.
* Example:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ingress-reader
  namespace: ingress-nginx
```

> 🧠 Namespace must exist first! Create it with:

```bash
kubectl create ns ingress-nginx
```

---

## Ingress & Port Forwarding

* See services:

```bash
kubectl get svc -A
```

* Port forward:

```bash
kubectl port-forward svc/my-service 8080:80
```

* Or expose via `NodePort` and hit via `localhost:<nodePort>`

> 🧠 Ingress controllers require proper RBAC and service exposure

---

## Troubleshooting & Gotchas

* `Connection refused`: kubeconfig not set, wrong port, or cluster down
* `Forbidden`: Missing RoleBindings or wrong ServiceAccount
* `No external IP`: Ingress service not ready or misconfigured
* `ImagePullBackOff`: Image is private or not pushed to a registry accessible by the cluster

---

## Git Survival Commands

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/yourusername/repo.git
git push -u origin main
```

> 🧠 If push fails, check credentials, remote URL, or if `main` exists already.

---

## SRE Wisdom & Mindset

* Don't memorize. **Build muscle memory.** Use your own playbooks.
* Real engineers use **cheat sheets, not just Google**.
* Always question: *"Is this ephemeral or persistent?"*
* If something works, **note it down immediately**.
* If it feels like magic, trace the pieces until you understand them.

> ✍️ Rule of thumb: If you had to debug this same thing again in 3 months, would Future You know what to do?

---

### ✅ To Do

* Add sections on Prometheus/Grafana, Datadog, CI/CD pipeline flows
* Add real-world examples from work (e.g., Datadog APM tweaks, IAM traps, RDS failovers)
* Include YAML templates and Helm charts for common use cases
