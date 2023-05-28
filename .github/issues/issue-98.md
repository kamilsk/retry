---
code: null
id: MDU6SXNzdWUyODA3ODIxODQ=
databaseId: 280782184
number: 98
url: https://github.com/kamilsk/retry/issues/98
title: "fix check-code-quality"
issueType: null
labels:
  - "type:bug"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-10T04:45:31Z
updatedAt: 2017-12-19T08:48:03Z
lastEditedAt: null
closedAt: 2017-12-19T08:48:03Z
---

# fix check-code-quality

```
.PHONY: check-code-quality
check-code-quality: ARGS = \
	--exclude=".*_test\.go:.*error return value not checked.*\(errcheck\)$$" \
	--exclude="duplicate of.*_test.go.*\(dupl\)$$" \
	--vendor --deadline=2m ./...
check-code-quality: docker-tool-gometalinter
```

also define correct deadline: `time make check-code-quality`

blocked by:
- https://github.com/kamilsk/shared/issues/301
- https://github.com/kamilsk/shared/issues/298

<!-- 2017-12-10T12:32Z https://github.com/kamilsk/retry/issues/98#issuecomment-350545127
extend git ignore

```
# go test and code quality report
*.out
```

add code quality report

```
.PHONY: code-quality-report
code-quality-report:
	time make check-code-quality | tail +7 | tee report.out | pbcopy
```
-->
