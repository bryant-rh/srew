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

## First Step

deploy srew-server
```Bash
cd deploy/srew-server

#部署mysql
kubectl apply -f mysql5.7-deploy.yaml

# 连接DB, 导入sql
srew-server.sql

#部署cm-server
kubectl apply -f srew-server-deploy.yaml
```



## Second Step
使用
https://github.com/bryant-rh/srew-rot  

创建index 至数据库


## Third Step

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



# Demo

## Search (搜索当前所有的plugin)
```Bash
$ srew search

+---------------+--------------------------------+----------------+-----------+---------+
|     NAME      |          DESCRIPTION           | LATEST VERSION | INSTALLED | UPGRADE |
+---------------+--------------------------------+----------------+-----------+---------+
| resource-view | Display Resource               | v0.1.0         | no        | yes     |
|               | (CPU/Memory/PodCount) Usage    |                |           |         |
|               | and Request ...                |                |           |         |
+---------------+--------------------------------+----------------+-----------+---------+
| tr            | Show a tree of object          | v0.3.0         | yes       | no      |
|               | hierarchies through            |                |           |         |
|               | ownerReferences                |                |           |         |
+---------------+--------------------------------+----------------+-----------+---------+
| tree          | Show a tree                    | v0.5.0         | no        | yes     |
+---------------+--------------------------------+----------------+-----------+---------+

```

```Bash
$ srew search -A

+---------------+--------------------------------+---------+-----------+
|     NAME      |          DESCRIPTION           | VERSION | INSTALLED |
+---------------+--------------------------------+---------+-----------+
| resource-view | Display Resource               | v0.1.0  | no        |
|               | (CPU/Memory/PodCount) Usage    |         |           |
|               | and Request ...                |         |           |
+---------------+--------------------------------+---------+           +
| tr            | Show a tree of object          | v0.2.0  |           |
|               | hierarchies through            |         |           |
|               | ownerReferences                |         |           |
+               +                                +---------+-----------+
|               |                                | v0.3.0  | yes       |
|               |                                |         |           |
|               |                                |         |           |
+---------------+--------------------------------+---------+-----------+
| tree          | Show a tree of                 | v0.4.0  | no        |
+               +--------------------------------+---------+           +
|               | Show a tree                    | v0.5.0  |           |
+---------------+--------------------------------+---------+-----------+
```


## List （展示本地安装的plugin）

```Bash
$ srew list

+--------+---------+
| PLUGIN | VERSION |
+--------+---------+
| tr     | v0.3.0  |
+--------+---------+
```

```Bash
$ srew list --no-format

PLUGIN  VERSION 
tr      v0.3.0 

```

## Info (查看插件详情)

```Bash
$ srew info tr

NAME    : tr
VERSION : v0.3.0
URI     : https://github.com/ahmetb/kubectl-tree/releases/download/v0.4.2/kubectl-tree_v0.4.2_darwin_amd64.tar.gz
SHA256  : 7369dc8d2d473908e15bf94afa64621e5c170a60eaf5ef1c55d99b03e2bf2d34
HOMEPAGE: https://github.com/ahmetb/kubectl-tree
DESCRIPTION: 
This plugin shows 
```

## Install （安装插件）
```Bash
#默认安装最新版本
srew install tr

# 可指定版本进行安装
srew install tr --version v0.2.0 -v 4

```

## Uninstall （卸载插件）

```Bash
srew uninstall tr
```

## Upgrade (升级插件)

```Bash
#不指定插件名称，会依次更新本地安装的所有插件
srew upgrade 

## 指定插件名称进行更新
srew upgrade tr

```