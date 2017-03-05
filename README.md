# Introduction

This Repository contains one core piece of forjj cli, to handle the
command line interface and help transferring data to FORJJ plugins.

# Building/Testing modules

**Requirements:**
- docker 1.9 or higher
- git 1.8 or higher

Modules are built thanks to docker in golang 1.7.4 and uses
[glide 0.12.1](https://github.com/Masterminds/glide) for project
dependencies.
So, you do not need them on your workstation.
See [glide/Dockerfile for details](glide/Dockerfile)

To build/test them:

Create the GO environment:
```bash
export GOPATH=~/go
mkdir -p $GOPATH/src
cd $GOPATH/src
git clone https://githib.com/forj-oss/forjj_modules
cd forjj_modules
```

Load forjj_modules build environment:
```bash
source build-env.sh
```

Start building/testing it:
```bash
build.sh
```

This last command installs dependencies in vendor sub directory and
start go test on all modules.

# Packages

Currently, there is no packages for those modules.
