[![Pre-compiled binaries](http://img.shields.io/badge/Precompiled-Download-green.svg)](http://beta.gobuild.io/github.com/tsloughter/hk-slug)

# hk-slug

A plugin to the fast Heroku CLI [`hk`](https://github.com/heroku/hk) for deploying via
the [Heroku Slug API](https://devcenter.heroku.com/articles/platform-api-deploying-slugs).

Based on Bo Jeane's [hk-deploy](https://github.com/bjeanes/hk-deploy) and using Naaman's [slug](https://github.com/naaman/slug).

## Usage

```sh-session
$ hk slug app.tar.gz
Deploying...
Done (v1)
```

## Install

### Pre-compiled Binaries

Pre-compiled binaries are available [here](http://beta.gobuild.io/github.com/tsloughter/hk-slug).

After unarchiving, stick the `slug` binary in `/usr/local/lib/hk/plugin` or your custom `$HKPATH`.

### Source install

Make sure you have Go (only 1.3 has been tested) installed.

```sh-session
$ go get github.com/tsloughter/hk-slug
$ cd $(go env GOPATH)/src/github.com/tsloughter/hk-slug
$ go build
$ mkdir -p /usr/local/lib/hk/plugin # or any custom $HKPATH
$ mv ./hk-slug /usr/local/lib/hk/plugin/slug
$ hk help slug
Usage: hk slug TARBALL

Deploy the specified directory to Heroku
```
