# Oasis API

## Installing Golang
 1. Download the go.1.11 archive from [here](https://golang.org/dl/) and extract it into /usr/local, creating a Go tree in /usr/local/go. For example: `sudo tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz`
2. Add /usr/local/go/bin to the PATH environment variable. You can do this by adding this line to your /etc/profile (for a system-wide installation) or $HOME/.profile (or .zshrc): `export PATH=$PATH:/usr/local/go/bin` 
3. Create a `go` folder in your $HOME folder with the following command: `cd ~/ && mkdir -p go/bin go/src go/pkg`
4. Add the bin folder to your $HOME/.profile (or .zshrc) file: `export PATH=$PATH:$HOME/go/bin`

## Installing VS Code GO tools
1. `cd $HOME/go`
2. `export GOPATH=$PWD`
3. `code .` (to open code in this folder with the exported variables)
4. Install the vs-code go extension
5. Reload VS Code
6. Press `ctrl+shift+p` and type `GO: Install/Update Tools`
7. Wait for it to finish and you can now open any project

## Installing dep
1. `cd ~/`
2. `mkdir -p go/bin go/src go/pkg` (if not created when installing golang)
3. `cd go`
4. `export GOPATH=$PWD`
5. `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

## Using dep
First you have to go to where our main package is; in our case :
`cd src/com/merkinsio/oasis-api`.
Then you can use any of the options below

### Adding a package
`dep ensure -add github.com/pkg/foo@^1.0.1`

### Installing all packages in `Gopkg.toml`
`dep ensure`

## Installing direnv
1. Download the latest binary build from [here](https://github.com/direnv/direnv/releases) corresponding to your architecture. For Ubuntu 64bits download `direnv.linux-amd64`.
2. Rename it to direnv: `mv direnv.linux-amd64 direnv`
3. Give it execution permission: `chmod +x direnv` 
4. Put it somewhere in your PATH. E.g.: `sudo cp direnv /usr/local/bin`
5. When the following message appears `direnv: error .envrc is blocked. Run `direnv allow` to approve its content.` execute `direnv allow`
6. If using vs-code, install the direnv extension

## Using local mongodb with docker
`docker run --name oasis -p 27017:27017 -d mongo:4.0`

## Working with binary
First you have to go to where our main package is; in our case :
`cd src/com/merkinsio/oasis-api`.
Then you build the project using:
`go build`
For more information about go build click [here](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies) 
For more information about the enviromen variables when building click [here](https://golang.org/cmd/go/#hdr-Environment_variables)
And finally to run the api yo use:
`./oasis-api`