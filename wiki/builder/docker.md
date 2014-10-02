# Docker Info

## Enable connectivity to host with CoreOS

### Important Notes

It is important to note that docker has a default host bridge which is 
`172.17.42.1`. In the below setup of enabling the remote api, I configure
`2375` as the remote tcp port. That port will need to referenced when
communicating to the remote host docker api.

## Customizing docker on CoreOS

These instructions can be found [at this link](https://coreos.com/docs/launching-containers/building/customizing-docker/)

The docker systemd unit can be customized by overriding the unit that ships with the default CoreOS settings. Common use-cases for doing this are covered below.

### Enable the Remote API on a New Socket

Create a file called `/etc/systemd/system/docker-tcp.socket` to make docker available on a TCP socket on port 2375.

```ini
[Unit]
Description=Docker Socket for the API

[Socket]
ListenStream=2375
BindIPv6Only=both
Service=docker.service

[Install]
WantedBy=sockets.target
```

Then enable this new socket:

```sh
systemctl enable docker-tcp.socket
systemctl stop docker
systemctl start docker-tcp.socket
systemctl start docker
```

Test that it's working:

```sh
docker -H tcp://127.0.0.1:2375 ps
```

You can also do:
```sh
curl http://172.17.42.1:2375/containers/json?all=1
```

## Testing connectivity to back to host

The following was the dockerfile used to test connectivity to the coreos
os system.

    #
    # Voom Test Dockerfile
    # How to use: (cd into the dir where the Dockerfile is first)
    #   to build: docker build -t voomtest .
    #   to run:   docker run -it voomtest /bin/bash
    #
    FROM ubuntu
    EXPOSE 80
    RUN apt-get update
    RUN apt-get install -y curl
    cmd curl "http://172.17.42.1:2375/containers/json?all=1"

