# kutti

The kutti CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/kuttiproject/kutti)](https://goreportcard.com/report/github.com/kuttiproject/kutti)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/kuttiproject/kutti?include_prereleases)

This repository contains the CLI for the kutti project. The CLI is inspired by the docker CLI, and allows management of Clusters and Nodes. The physical implementation of underlying networks and hosts is handled via Drivers. Each driver is also responsible for providing a repository to download host templates for supported Kubernetes versions.

The CLI includes a simple SSH client for connecting to nodes it creates.


<img src="https://github.com/kuttiproject/driver-vbox-images/blob/main/attachments/icon/kutta.png?raw=true" width="32" height="32" /> Icon made by [Freepik](https://www.freepik.com) from [Flaticon](http://www.flaticon.com)