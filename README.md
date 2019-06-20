# hook2xmpp


[![pipeline status](https://dev.sum7.eu/genofire/hook2xmpp/badges/master/pipeline.svg)](https://dev.sum7.eu/genofire/hook2xmpp/pipelines)
[![coverage report](https://dev.sum7.eu/genofire/hook2xmpp/badges/master/coverage.svg)](https://dev.sum7.eu/genofire/hook2xmpp/pipelines)
[![Go Report Card](https://goreportcard.com/badge/dev.sum7.eu/genofire/hook2xmpp)](https://goreportcard.com/report/dev.sum7.eu/genofire/hook2xmpp)
[![GoDoc](https://godoc.org/dev.sum7.eu/genofire/hook2xmpp?status.svg)](https://godoc.org/dev.sum7.eu/genofire/hook2xmpp)


## Get hook2xmpp

#### Download

Latest Build binary from ci here:

[Download All](https://dev.sum7.eu/genofire/hook2xmpp/-/jobs/artifacts/master/download/?job=build-my-project) (with config example)

[Download Binary](https://dev.sum7.eu/genofire/hook2xmpp/-/jobs/artifacts/master/raw/bin/hook2xmpp?inline=false&job=build-my-project)

#### Build

```bash
go get -u dev.sum7.eu/genofire/hook2xmpp
```

## Configure

see `config_example.toml`

## Start / Boot

_/lib/systemd/system/hook2xmpp.service_ :
```
[Unit]
Description=hook2xmpp
After=network.target
# After=ejabberd.service
# After=prosody.service

[Service]
Type=simple
# User=notRoot
ExecStart=/opt/go/bin/hook2xmpp --config /etc/hook2xmpp.conf
Restart=always
RestartSec=5sec

[Install]
WantedBy=multi-user.target
```

Start: `systemctl start hook2xmpp`
Autostart: `systemctl enable hook2xmpp`
