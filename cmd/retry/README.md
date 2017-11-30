> # cmd/retry
>
> `retry` provides functionality to repeat terminal commands.

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#utilities)
[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![Build Status](https://travis-ci.org/kamilsk/retry.svg?branch=master)](https://travis-ci.org/kamilsk/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/retry)](https://goreportcard.com/report/github.com/kamilsk/retry)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/retry/badge.svg)](https://coveralls.io/github/kamilsk/retry)
[![GoDoc](https://godoc.org/github.com/kamilsk/retry?status.svg)](https://godoc.org/github.com/kamilsk/retry)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](../../LICENSE)

## Concept

```bash
$ retry -limit=3 -backoff=lin:500ms -- curl example.com
```

## Documentation

```
Usage: retry [-timeout Timeout] [--debug] [--notify] [strategy flags] -- command

The strategy flags
    --infinite
        Infinite creates a Strategy that will never stop repeating.

    -limit=X
        Limit creates a Strategy that limits the number of attempts that Retry will
        make.

    -delay=Xs
        Delay creates a Strategy that waits the given duration before the first
        attempt is made.

    -wait=Xs,...
        Wait creates a Strategy that waits the given durations for each attempt after
        the first. If the number of attempts is greater than the number of durations
        provided, then the strategy uses the last duration provided.

    -backoff=:algorithm
        Backoff creates a Strategy that waits before each attempt, with a duration as
        defined by the given backoff.Algorithm.

    -tbackoff=":algorithm :transformation"
        BackoffWithJitter creates a Strategy that waits before each attempt, with a
        duration as defined by the given backoff.Algorithm and jitter.Transformation.

:algorithm
    inc:Xs,Ys
        Incremental creates a Algorithm that increments the initial duration
        by the given increment for each attempt.

    lin:Xs
        Linear creates a Algorithm that linearly multiplies the factor
        duration by the attempt number for each attempt.

    exp:Xs,Y
        Exponential creates a Algorithm that multiplies the factor duration by
        an exponentially increasing factor for each attempt, where the factor is
        calculated as the given base raised to the attempt number.

    binexp:Xs
        BinaryExponential creates a Algorithm that multiplies the factor
        duration by an exponentially increasing factor for each attempt, where the
        factor is calculated as "2" raised to the attempt number (2^attempt).

    fib:Xs
        Fibonacci creates a Algorithm that multiplies the factor duration by
        an increasing factor for each attempt, where the factor is the Nth number in
        the Fibonacci sequence.

:transformation
    full
        Full creates a Transformation that transforms a duration into a result
        duration in [0, n) randomly, where n is the given duration.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://www.awsarchitectureblog.com/2015/03/backoff.html

    equal
        Equal creates a Transformation that transforms a duration into a result
        duration in [n/2, n) randomly, where n is the given duration.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://www.awsarchitectureblog.com/2015/03/backoff.html

    dev:X
        Deviation creates a Transformation that transforms a duration into a result
        duration that deviates from the input randomly by a given factor.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://developers.google.com/api-client-library/java/google-http-java-client/backoff

    ndist:X
        NormalDistribution creates a Transformation that transforms a duration into a
        result duration based on a normal distribution of the input and the given
        standard deviation.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

Examples:
    retry -limit=3 -backoff=lin:10ms -- curl http://example.com
    retry -tbackoff="lin:10s full" --debug -- curl https://example.com
    retry -timeout=500ms --notify --infinite -- git pull

Version 3.0.0 (commit: ..., build date: ..., go version: go1.9, compiler: gc, platform: darwin/amd64)
```

### Complex example

```
$ retry -limit=3 -backoff=lin:10ms -- /bin/sh -c 'echo "trying..."; exit 1'
trying...
#2 attempt at 17.636458ms...
trying...
#3 attempt at 48.287964ms...
trying...
error occurred: "exit status 1"
$ retry -timeout=500ms --infinite -- /bin/sh -c 'echo "trying..."; exit 1'
trying...
...
trying...
#N attempt at 499.691521ms...
error occurred: "context deadline exceeded"
```

## Installation

### Brew

```bash
$ brew install kamilsk/tap/retry
```

### Binary

```bash
$ export SEM_V=3.0.0    # all available versions are on https://github.com/kamilsk/retry/releases
$ export REQ_OS=Linux   # macOS and Windows are also available
$ export REQ_ARCH=64bit # 32bit is also available
$ wget -q -O retry.tar.gz \
      https://github.com/kamilsk/retry/releases/download/${SEM_V}/retry_${SEM_V}_${REQ_OS}-${REQ_ARCH}.tar.gz
$ tar xf retry.tar.gz -C "${GOPATH}"/bin/
$ rm retry.tar.gz
```

### From source code

```bash
$ go get -d github.com/kamilsk/retry
$ cd "${GOPATH}"/src/github.com/kamilsk/retry
$ make cmd-deps-local # or cmd-deps, if you don't have the dep binary but have the docker
$ make cmd-install
```

## Command-line completion

### Useful articles

- [Command-line completion | Docker Documentation](https://docs.docker.com/compose/completion/)

### Bash

Make sure [bash completion](https://github.com/scop/bash-completion) is installed.

- On a current Linux (in a non-minimal installation), bash completion should be available.
- On a macOS, install by the command `brew install bash-completion`.

Place the completion script in `/etc/bash_completion.d/` (or `$(brew --prefix)/etc/bash_completion.d/` on a macOS):

```bash
$ sudo curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/retry.sh \
            -o /etc/bash_completion.d/retry
```

On a macOS, add the following to your `~/.bash_profile`:

```bash
if [ -f $(brew --prefix)/etc/bash_completion ]; then
    source $(brew --prefix)/etc/bash_completion
fi
```

If you're using MacPorts instead of brew you'll need to slightly modify your steps to the following:

- Run `sudo port install bash-completion` to install bash completion.
- Add the following lines to `~/.bash_profile`:
```bash
if [ -f /opt/local/etc/profile.d/bash_completion.sh ]; then
    source /opt/local/etc/profile.d/bash_completion.sh
fi
```

You can source your `~/.bash_profile` or launch a new terminal to utilize completion.

### Zsh

Place the completion script in your `/path/to/zsh/completion`, using, e.g., `~/.zsh/completion/`:

```bash
$ mkdir -p ~/.zsh/completion
$ curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/retry.zsh \
       -o ~/.zsh/completion/_retry
```

Include the directory in your `$fpath`, e.g., by adding in `~/.zshrc`:

```bash
fpath=(~/.zsh/completion $fpath)
```

Make sure `compinit` is loaded or do it by adding in `~/.zshrc`:

```bash
autoload -Uz compinit && compinit -i
```

Then reload your shell:

```bash
exec $SHELL -l
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/retry)
[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

## Notes

- made with ❤️ by [OctoLab](https://www.octolab.org/)

[![Analytics](https://ga-beacon.appspot.com/UA-109817251-1/retry/cmd)](https://github.com/igrigorik/ga-beacon)
