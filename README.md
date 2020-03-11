# ownCloud Infinite Scale: DevLDAP

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-devldap/status.svg)](https://cloud.drone.io/owncloud/ocis-devldap)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/afbd4455ca1a4833966fb69edec87cdb)](https://www.codacy.com/manual/owncloud/ocis-devldap?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-devldap&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-devldap?status.svg)](http://godoc.org/github.com/owncloud/ocis-devldap)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-devldap)](http://goreportcard.com/report/github.com/owncloud/ocis-devldap)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-devldap.svg)](http://microbadger.com/images/owncloud/ocis-devldap "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/devldap/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/ocis-devldap/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.12. For the frontend it's also required to have [NodeJS](https://nodejs.org/en/download/package-manager/) and [Yarn](https://yarnpkg.com/lang/en/docs/install/) installed.

```console
git clone https://github.com/owncloud/ocis-devldap.git
cd ocis-devldap

make generate build

./bin/ocis-devldap -h
```

### Updating demo users

The demo users can be found in `assets/data.json`.
After editing the file, please also run `make generate` to embed them into the binary.

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2019 ownCloud GmbH <https://owncloud.com>
```
