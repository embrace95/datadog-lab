GitOps Bible

🧭 What is GitOps?

GitOps is a workflow where Git becomes the single source of truth for infrastructure and application state. You describe the desired state of your systems in Git, and a controller automatically ensures that the actual cluster state matches the Git state.

GitOps ≠ "pushing YAMLs manually". It means Git drives everything, including:

Deployments

Infra changes

Rollbacks

Audit trails

🛠️ Core Tools in a GitOps Stack

Tool

Purpose

GitHub/GitLab

Version control & source of truth

ArgoCD/Flux

GitOps controller to sync Git with K8s

Helm/Kustomize

Template & customize manifests

VS Code

Your local GitOps authoring environment

📁 Repo Layout Example (Your Current Setup)

├── Ai-Ops-playbook.md
├── Ai-LogAnalysis.md
├── gitops-bible.md   <-- This file
└── obs-k8s-demo/
    ├── Dockerfile
    ├── deployment.sh
    ├── deploy_environment.sh
    ├── ingress-reader-role.yaml
    ├── ingress-reader-rolebinding.yaml
    ├── ingress-reader-sa.yaml
    ├── k3dconfig.yaml
    ├── Makefile
    ├── observk8s-deployment.yaml
    ├── observk8s-service.yaml
    ├── README.md
    └── flowchaos.yaml

🔍 GitOps Case Study: obs-k8s-demo

What happened:

You cloned a project and built a custom container image.

You spun up a K3D Kubernetes cluster using k3dconfig.yaml.

You tried to track the folder in Git, but it was misidentified as a submodule.

You fixed it using:

git rm --cached obs-k8s-demo
rm .gitmodules
git add obs-k8s-demo/
git commit -m "Add demo as regular folder"

You verified the contents were tracked (243 files recognized!)

You used the VS Code Source Control tab to visually confirm.

Lesson:

Submodules are rarely what you want unless you're intentionally pulling in an external repo. If you just have a local project folder — treat it as part of your main Git tree.

📦 GitOps Checklist (Minimal)



🔁 Add This Project to ArgoCD Later (Optional)

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: observk8s
  namespace: argocd
spec:
  destination:
    name: ''
    namespace: default
    server: 'https://kubernetes.default.svc'
  project: default
  source:
    path: obs-k8s-demo
    repoURL: 'https://github.com/YOUR_USERNAME/datadog-lab.git'
    targetRevision: HEAD
  syncPolicy:
    automated:
      prune: true
      selfHeal: true

🧠 Bonus Tips

Add .gitignore and .gitattributes to sanitize commits

Commit descriptive messages

Use GitHub Actions or pre-commit hooks to lint YAMLs

Tag major releases

📚 Suggested Reading

GitOps Principles

ArgoCD Docs

Flux Docs

Kustomize
