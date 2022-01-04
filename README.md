# garlicshare

<p align="center">
  <img src="https://i.imgur.com/7jifuFY.png" width=100 height=100 >
</p>

# Dependencies
garlicshare requires at least Tor >= 0.3.5

https://www.torproject.org/it/download/

# Installation
garlicshare can be installed in different ways:

## **Go Packages**

<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png" width=95 height=35>

 throught [Golang Packages](https://go.dev/) (golang package manager)
 
```bash
go get github.com/R4yGM/garlicshare
```
**this will work for every platform**

## **Docker**

<img src="https://cdn3.iconfinder.com/data/icons/logos-and-brands-adobe/512/97_Docker-512.png" width=35 height=35>

  if you don't have docker installed you can follow their [guide](https://docs.docker.com/engine/install/)
  
 first of all you have to pull the docker image (only **12.72 MB**) from the docker registry, you can see it [here](https://hub.docker.com/r/r4yan/garlicshare), if you don't want to pull the image you can also clone the repository and then build the image from the Dockerfile
 ```bash
docker pull r4yan/garlicshare:latest
  ```
 
  if you don't want to pull the image you can download or copy the dorkscout Dockerfile that can be found [here](https://github.com/R4yGM/garlicshare/blob/main/Dockerfile) and then build the image from the Dockerfile
  
  then if you want to launch the container you have to first create a volume that contains the files you want to share
  
  ```bash
  docker volume create --name garlicshare_files
  ``` 
  then copy the files you want to share in the volume data folder `/var/lib/docker/volumes/garlicshare_files/_data`
  ```bash
  cp file /var/lib/docker/volumes/garlicshare_files/_data
  ``` 
  a
  then you can start sharing by running :
  ```bash
docker run -v garlicshare_files:/garlicshare r4yan/garlicshare <options>
  ```
  replace the `<options>` with the options/arguments you want to give to garlicshare,
  example :
   ```bash
docker run -v garlicshare_files:/garlicshare r4yan/garlicshare upload -p garlicshare 
  ```
  **the path must be the same as the volume binding path! [read more here](https://docs.docker.com/storage/bind-mounts/#choose-the--v-or---mount-flag)**
  
  **this will work for every platform**

# Usage

garlicshare is very simple to use, you can view the program help with the -h or --help option

```bash
Usage:
  GarlicShare upload [flags]

Flags:
  -h, --help          help for upload
  -k, --key string    Password to download the files
  -p, --path string   Path

Global Flags:
      --config string   config file (default is $HOME/.GarlicShare.yaml)
```

<p align="center">
<a href="https://asciinema.org/a/5YUpQhY76MQE6vXDIVNNyK9T7" target="_blank"><img src="https://asciinema.org/a/5YUpQhY76MQE6vXDIVNNyK9T7.svg" /></a>
</p>
