> # shared:makefile-go-v1
>
> My snippets of GNU Make for Go environment.

## Integration

```
.PHONY: pull-makes
pull-makes:
	rm -rf makes
	git clone git@github.com:kamilsk/shared.git makes
	( \
	  cd makes && \
	  git checkout makefile-go-v1 && \
	  git branch -d master && \
	  echo '- ' $$(cat README.md | head -n1 | awk '{print $$3}') 'at revision' $$(git rev-parse HEAD) \
	)
	rm -rf makes/.git
```

## Useful articles

* [Go tooling essentials](https://rakyll.org/go-tool-flags/)
