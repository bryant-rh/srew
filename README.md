# What is Srew?
Srew is a plugin manager similar to the krew and brew command line tools

Supports the deployment of the server, which is used to manage its own operation and maintenance tools within the enterprise


Srew helps you:

+ discover all plugins,
+ install them on your machine,
+ and keep the installed plugins up-to-date.

Srew works across all major platforms, like macOS, Linux and Windows.

Srew also helps kubectl plugin developers: You can package and distribute your plugins on multiple platforms easily and makes them discoverable through a centralized plugin repository with Srew.


# Installing

+ macOS/Linux: bash/zsh, fish
+ Windows

## macOS/Linux
### Bash or ZSH shells
Make sure that git is installed.

1. Run this command to download and install Srew:

```Bash
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  SREW="srew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/bryant-rh/srew/releases/latest/download/${SREW}.tar.gz" &&
  tar zxvf "${SREW}.tar.gz" &&
  ./"${SREW}" install srew
)

```

2. Add the $HOME/.srew/bin directory to your PATH environment variable. To do this, update your .bashrc or .zshrc file and append the following line:
   
```Bash
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
```
and restart your shell.


3. Run srew to check the installation.


### Fish shell

1. Run this command in your terminal to download and install srew:

```Bash
begin
  set -x; set temp_dir (mktemp -d); cd "$temp_dir" &&
  set OS (uname | tr '[:upper:]' '[:lower:]') &&
  set ARCH (uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/') &&
  set SREW srew-$OS"_"$ARCH &&
  curl -fsSLO "https://github.com/bryant-rh/srew/releases/latest/download/${SREW}.tar.gz" &&
  tar zxvf $SREW.tar.gz &&
  ./$SREW install srew &&
  set -e SREW temp_dir &&
  cd -
end
```

2. Add the $HOME/.srew/bin directory to your PATH environment variable. To do this, update your config.fish file and append the following line:
 
```Bash
set -gx PATH $PATH $HOME/.krew/bin
```
and restart your shell.

3. Run srew to check the installation.


## Windows

1. Download srew.exe from the Releases page to a directory.

2. Launch a command prompt (cmd.exe) with administrator privileges (since the installation requires use of symbolic links) and navigate to that directory.

3. Run the following command to install srew:
```Bash
.\srew install srew
```
4. Add the %USERPROFILE%\.srew\bin directory to your PATH environment variable (how?)

5. Launch a new command-line window.

6. Run srew to check the installation.



