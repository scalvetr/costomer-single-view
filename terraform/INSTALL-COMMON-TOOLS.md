# MacOS: Install common tools

Install the following tools:
* Homebrew
* jq
* kubectl
* helm
* kubepfm

```shell
# Install xcode
xcode-select --install

# Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

brew install jq@1.6
brew install kubectl@1.22.2
brew install helm@3.7.1
brew install flowerinthenight/tap/kubepfm@1.5.8
#brew install txn2/tap/kubefwd@1.22.0

echo "source <(kubectl completion zsh)" >> ~/.zshrc
```