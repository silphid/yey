# Yey!

Yey is a user-friendly CLI fostering collaboration within and across dev/devops teams by allowing them to easily share and launch various docker container configurations.

# Motivation

Ever since my career converged onto DevOps many years ago, I have always been confronted with the complexity of dealing with multiple environments and Kubernetes clusters, which often required specific versions of tools (ie: a 1.14 Kubernetes cluster cannot be managed by `kubectl` 1.16). I often had to reinstall a different version of `kubectl` depending on the cluster I was connecting to. Also, onboarding new team members required them to follow detailed instructions to properly install and configure their environment, not to mention configuration drifts where some of them ended up with missing or incorrect versions of tools required by our DevOps bash scripts.

The solution seemed obvious: Package up all tools in docker containers and have everyone use them for interacting with clusters and other tasks involving our scripts. While docker containers are mostly leveraged for deploying applications to the cloud, they can also be extremely useful for standardizing and packaging up different sets of tools and environments for development and DevOps daily work.

However, using containers in that way comes with its own set of challenges:

- The `docker` CLI is not particularly friendly for launching multiple containers with sophisticated configs.
- Starting multiple sessions in parallel against same container is particularly teddious, as you have to mind using proper commands depending if you're starting the container for the first time (`docker run ...`), if it's already running and you must just shell into it (`docker exec ...`), etc.
- Changes to container's filesystem are ephemeral, so you must take great care of mounting local files and directories into your containers to have your work persist across sessions. For example, I like to mount my local $HOME dir in a volume, as well as common user files (ie: `~/.ssh`, `~/.gitconfig`...) And I constantly had to `cd` into the equivalent directory in container as I was in on my local machine.
- Passing all proper configurations as arguments to `docker` CLI for different images and use-cases (ie: environments, clusters...) is teddious and error-prone.
- There is no standard way of sharing launch configurations within and across teams. Docker-compose can be of some help here, but it's really more focused on configuring complex sets of interconnected containers.
- Some individuals may also need to override some launch configurations for their own specific needs.

Yey was designed to address all those challenges and more by abstracting all configurations into YAML files and presenting users with simple interactive prompts.

# Installation

## MacOS

```bash
$ brew tap silphid/yey
$ brew install yey
```

## Other platforms

- Download and install latest release for your platform from the GitHub [releases](https://github.com/silphid/yey/releases) page.
- Make sure you place the binary somewhere in your `$PATH`.

# Getting started

## Docker images

Yey comes with its own set of stock docker images - loaded with a bunch of useful DevOps CLI tools - that you can start using right away. However, you are really expected to create your own images for your own specific needs. You can either use yey's stock images as a base for you own images or just as inspiration to create yours from scratch.

The stock images are defined in the https://github.com/silphid/yey-images repo and are structured as follows:

- `yey-base`: base image for all other yey images, including a minimal set of useful tools (ie: `curl`, `git`, `jq`...)
  - `yey-devops`: kubernetes and terraform tools
    - `yey-aws`: AWS tools
    - `yey-gcp`: Google Cloud Platform tools
  - `yey-go`: Go dev environment
  - `yey-node`: Node dev environment
  - `yey-python`: Python dev environment

For more details on those images, see that project's [documentation](https://github.com/silphid/yey-images).

# Contexts

A yey context is a set of configs (image name, env vars, volume bindings, etc) for launching a docker container. You can define multiple contexts.

Here's an example context:

```yaml
image: alpine
mounts:
  ~/: /home
  ~/.ssh: /root/.ssh
env:
  GCP_ZONE: us-east1-a
  TZ: $TZ
entrypoint: zsh
```

## General structure of `.yeyrc.yaml`

## Resolving `.yeyrc.yaml` file(s)

Yey resolves the RC file in following order:

-

# Commands

## Global flags

- -h --help
- -v --verbose
  - Displays detailed output messages
- --dry-run
  - Only displays shell commands that it would execute normally
  - Automatically turns-on verbose mode
