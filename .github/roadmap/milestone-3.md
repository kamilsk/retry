---
code: null
id: MDk6TWlsZXN0b25lMjU5ODM0Ng==
databaseId: 2598346
number: 3
url: https://github.com/kamilsk/retry/milestone/3
title: "retry 3.0"
state: CLOSED
openIssueCount: 0
closedIssueCount: 40
dueOn: 2017-11-30T00:00:00Z
createdAt: 2017-06-22T12:19:20Z
updatedAt: 2018-05-02T18:38:00Z
closedAt: 2017-11-30T18:09:00Z
---

# retry 3.0

Roadmap:
- simplify code and refactor it;
- remove custom `flag` package;
  - ignore argument sequence;
- support go 1.9;
- review Limit strategy;
- sync with stable `egg`;
- move pseudo-generated code to pkg level.

## Issues

- [x] [[issue-57]] — simplify cmd/retry
- [x] [[issue-58]] — add github.com/cenkalti/backoff to research list
- [x] [[issue-59]] — inject go version
- [x] [[issue-60]] — notification after complete
- [x] [[issue-61]] — integration with new services
- [x] [[issue-62]] — improve debugging
- [x] [[issue-63]] — better work with color in tty
- [x] [[issue-64]] — http client example
- [x] [[issue-65]] — pinger example
- [x] [[issue-66]] — add ignore for some linters
- [x] [[issue-67]] — integrate with bash_completion
- [x] [[issue-68]] — add code of conduct
- [x] [[issue-69]] — configure gometalinter
- [x] [[issue-70]] — investigate ldflags
- [x] [[issue-71]] — replace context.Background by context.TODO in examples
- [x] [[issue-72]] — review examples
- [x] [[issue-73]] — add more blackbox tests
- [x] [[issue-74]] — add Go Code Review Comments
- [x] [[issue-75]] — refactor research.yml
- [x] [[issue-77]] — handle os.Interrupt signal
- [x] [[issue-78]] — improve work with stdout
- [x] [[issue-79]] — add example with os.Interrupt
- [x] [[issue-80]] — use Cancellation as a subset of Context
- [x] [[issue-81]] — update travis ci
- [x] [[issue-82]] — remove context package
- [x] [[issue-83]] — sync with semaphore package
- [x] [[issue-85]] — support different separators
- [x] [[issue-87]] — update makes
- [x] [[issue-88]] — move research to separated branch
- [x] [[issue-89]] — add patreon information
- [x] [[issue-90]] — resolve TODO related to ctx.Err
- [x] [[issue-91]] — restore old logic of Limit
- [x] [[issue-92]] — add error checker
- [x] [[issue-93]] — integrate spinner
- [x] [[issue-110]] — add quick start example
- [x] [[issue-111]] — integrate with scrutinizer
- [x] [[issue-112]] — go 1.10
- [x] [[issue-113]] — race in examples
- [x] [[issue-114]] — fix latest build on 1.7.4 instead of 1.10.0
- [x] [[issue-115]] — vendor deps
