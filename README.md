# arabian-nights
A Docker job to unseal Vault instances running on K8s, storing the unseal keys in Kubernetes secrets.

## Background

When starting a Vault cluster with high availability enabled, Vault instances boot _sealed_. If your Vault cluster is running secured in Kubernetes, you would normally port-forward your Vault cluster to your local machine and manually unseal the Vault.

This repository contains a binary for unsealing Vault instances and storing the unseal keys in Kubernetes secrets. The principal idea is that this tool should run from inside of your cluster, and an operator can be trusted to retrieve the unseal keys from the Kubernetes secrets at a later time, cleaning up the secrets once the unseal keys have been downloaded and distributed to the shardbearers.

## When is this useful?

Because the tool can be deployed as a Batch job on Kubernetes, it is useful when your Kubernetes cluster initialization is automated. For example, suppose your cluster is created with an IaC tool like Terraform, Pulumi, or CloudFormation. Because these tools are often running in CI environments, or triggered as part of a GitOps workflow, they are automating the creation of the Vault cluster in an environment that may not have a location to securely save the shards. Using this tool, you can keep the shares in Kubernetes and retrieve them outside of the CI process.

## How it Works

This Go binary requires as input the name of a Kubernetes Service resource running the Vault cluster. The Service must be running in the same namespace as this binary.

The binary will collect the Vault instances that are members of the Vault cluster by inspecting the pods backing the Service.

Then, it will sequentially unseal each of the Vault instances, storing the unseal key ina Kubernetes secret.
