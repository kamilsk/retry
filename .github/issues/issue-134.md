---
code: null
id: MDU6SXNzdWU0NzcyMDUwNjE=
databaseId: 477205061
number: 134
url: https://github.com/kamilsk/retry/issues/134
title: "problems with codeclimate.com integration"
issueType: null
labels:
  - "type:bug"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-08-06T07:00:46Z
updatedAt: 2026-11-08T14:02:37Z
lastEditedAt: null
closedAt: 2019-08-08T18:21:50Z
---

# problems with codeclimate.com integration

https://travis-ci.org/kamilsk/retry/jobs/567990296

<!-- 2019-08-09T08:07Z https://github.com/kamilsk/retry/issues/134#issuecomment-519824218
related to https://github.com/codeclimate/test-reporter/issues/378
-->

<!-- 2019-08-09T11:06Z https://github.com/kamilsk/retry/issues/134#issuecomment-519879050
```bash
#!/usr/bin/env bash

SECRET=...

docker run \
       --rm -it \
       -e CC_TEST_REPORTER_ID=${SECRET} \
       -e CODECLIMATE=https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 \
       -e GO111MODULE=on \
       -w /pkg/ \
       -v $(pwd):/pkg/ \
       golang:1.12 \
       bash -c 'curl -L $CODECLIMATE > ./cc-test-reporter; \
                chmod +x ./cc-test-reporter; \
                ./cc-test-reporter before-build; \
                make test-with-coverage-profile; \
                cat c.out && ./cc-test-reporter -v; \
                ./cc-test-reporter after-build -d -t gocov -p $(go list -m)'
```
-->
