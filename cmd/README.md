> # retrier/cmd
>
> Package cmd provides Command Line Interface to repeat terminal commands.

## retry

```bash
$ retry --infinite
#       -limit=X
#       -delay=Xs
#       -wait=Xs,Ys
#       -backoff=<algorithm>
#       -tbackoff=<algorithm>,<jitter>
# <algorithm>
#       -backoff=inc[Xs,Ys]
#       -backoff=lin[Xs]
#       -backoff=epx[Xs,Y]
#       -backoff=binexp[Xs]
#       -backoff=fib[Xs]
# <jitter>
#       -tbackoff=...,full
#       -tbackoff=...,equal
#       -tbackoff=...,dev[X]
#       -tbackoff=...,ndist[X]
# full example
$ retry -limit=3 -backoff=lin[10ms] -- curl http://unknown.host
curl: (52) Empty reply from server
curl: (52) Empty reply from server
curl: (52) Empty reply from server
$ retry --infinite -- curl http://unknown.host
curl: (52) Empty reply from server
...
curl: (52) Empty reply from server
...
```
