# Commands

- `yey use [context or "base" or "none"]`
  - validates context
  - saves context in `~/.yey/state.yaml`
  - omitting context prompts for context from list of available contexts + "base" + "none"
- `yey run [context or "base" or "none" or "?"]`
  - reads `~/.yey/state.yaml`
  - omitting context uses default context (if default not set or "?", prompts for context)
  - if no version set yet, does a `yey upgrade`
  - merges env vars and volume mounts from `shared.yaml` and `user.yaml`
  - keep only volume mounts that exist locally
  - if tag is `latest` or `""` does a forced pull
  - starts docker container
- `yey stop [context or "base" or "none" or "all"]`
  - omitting context prompts for context
- `yey get`
  - `contexts`
    - lists available contexts
  - `containers`
    - lists active containers
  - `context`
    - displays current context
- `yey remove/rm [context or "" or "?" or "all"]`
  - omitting context uses default, "?" prompts for context
  - stops running containers
  - deletes yey docker images and containers
- `yey install [ARGS...]`
  - prompts user for all install properties not passed as arguments
    - Do you want to configure a git repo from which to load remote shared configs?
      - Base URL to raw config git repo config files
      - Git token
        - Maybe we can automate that process/flow!
          Something like: https://github.com/chrisdickinson/get-github-token
          There must be a way to pop a browser and have github (or other) prompt
          user to authorize yey and automatically create a token, like
          for google cloud?
  - stores answers in `~/.yey/state.yaml`
  - if remote shared configs configured, copy `user.yaml` to `~/.yey/`
    (only if doesn't already exist)

## Global flags

- -h --help
- -v --verbose
  - Displays detailed output messages
- --dry-run
  - Only displays shell commands that it would execute normally
  - Automatically turns-on verbose mode

# Git repos

- `yey-containers.git`: monorepo with multiple Dockerfiles, all built individually
  and pushed to dockerhub.
  - `/base/`: Docker image with minimum functionality (injection...)
  - `/devops/`: FROM yey-base + kubectl, k9s, helm, helmfile...
  - `/go/`: FROM yey-devops + go sdk
  - `/node/`: FROM yey-devops + node, TS...
  - `/totum/`: FROM yey-devops + go + node...
- `yey-config.git`: repo to be optionally forked and customized by user/company.
  - `/shared.yaml`
  - `/user.yaml`
- `yey-cli.git`: the yey cli go source code, installable via `go install` or `curl ...`.

# Config files structure

- ~/.yey/
  - state.yaml
  - user.yaml: user-specific config overrides
- yey Git Repo Clone Dir
  - config
    - shared.yaml: base configs shared by all users
    - user.yaml: default user config file copied to home during install
- In any parent folder (v2?)
  - .yeyrc

# ~/.yey/state.yaml file format

```yaml
version: 2021.04
remote: https://github.com/asdf/asdf
token: jasd5DasdmYdndIIiD333Mdojasd5DasdmYdndIIiD333Mdodd
context: cluster1
```

# Config files format

```yaml
base:
  image: silphid/yey
  tag: latest
  container: yey # container name for the XXX part: yey-XXX-context-tag

  env:
    CLOUD: aws # supported clouds: aws, gcp
    REGION: us-east-1

  mounts:
    $HOME/.ssh: /root/.ssh
    $HOME/.gitconfig: /root/.gitconfig
    $HOME/.aws: /root/.aws
    $HOME/.config/gh: /root/.config/gh
    $HOME/.cfconfig: /root/.cfconfig

contexts:
  cluster1:
    env:
      KUBE_CONTEXT: cluster1
      # REGION: us-east-2
  cluster2:
    env:
      KUBE_CONTEXT: cluster2
```

# Installation

- User installs yey cli binary, either by:
  - `$ go install github.com/silphid/yey-cli`
  - `$ brew install yey`
  - `$ curl https://github.com/silphid/yey-cli/raw/.../install.sh | bash`
- User then configures yey via:
  - `$ yey install` (see "Commands" section above)
  - (the brew and curl options could maybe run this automatically after installing the binary?)

# Todo

- Remove "none" context
- Add support for multiple images (other than yey)
  - Sub-folders under `~/.yey` and git `/config` that can specify extra
    `user.yaml` and `shared.yaml` files.
  - The sub-folder name can be used as a prefix to the context name (ie:
    `image1/context1`) or we add an extra `yey image {image1}` command to
    select it?
- Detect current directory on local machine, find equivalent directory (if any)
  in container by looking at all bind mounts and use that as initial current directory
  (and try using Docker's option for specifying startup dir instead of old approach)

# Ideas/questions

- yey user support??? can we run our containers as not root via context config or otherwise? SPIKE
- rm --all remove all contexts
- support in context file for mapping ports
- support in context file for passing arbitrary docker flags
