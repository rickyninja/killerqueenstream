# nginx with rtmp module

Install based on [peer5 docs](https://docs.peer5.com/guides/setting-up-hls-live-streaming-server-using-nginx/).

## install nginx depends

`sudo apt-get install build-essential libpcre3 libpcre3-dev libssl-dev zlib1g-dev`

## build and install from source

```
sudo apt install git
mkdir git
cd git
git clone https://github.com/arut/nginx-rtmp-module
git clone https://github.com/nginx/nginx.git
cd nginx
auto/configure --with-http_ssl_module --add-module=../nginx-rtmp-module
make -j `nproc`
sudo make install
```

## install rtmp nginx config file

install [nginx.conf](/etc/nginx/nginx.conf) to /usr/local/nginx/conf/nginx.conf

## use systemd to start rtmp service

install [nginx-rtmp.service](/etc/systemd/system/nginx-rtmp.service) file in /etc/systemd/system/

```
sudo systemctl enable nginx-rtmp.service
sudo systemctl start nginx-rtmp.service
```
