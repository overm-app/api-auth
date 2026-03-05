## Description
<!-- What does this PR do? Reference the issue: "Closes #NNN" -->

Closes #

## Type of change
- [ ] `feat` — new feature
- [ ] `fix` — bug fix
- [ ] `chore` — tooling, dependencies, config
- [ ] `refactor` — code change with no functional difference
- [ ] `docs` — documentation only

---

## Contracts
- [ ] No contract changes required
- [ ] Contract updated in `overm-app/contracts` before this code was written
- [ ] Breaking change — major version bump required
- [ ] Necessary adjustments — contract updated during implementation

**If contract changed, which file?**
<!-- user-auth-api.yaml / user-auth-web.yaml / recipe-catalog.yaml / etc -->

**If necessary adjustment, what changed and why?**
<!-- Be specific — "added updated_at to User schema, DB fills it automatically and was missed during planning" -->

---

## Database
- [ ] No migration required
- [ ] Migration added in `migrations/`
- [ ] Migration is backward compatible
- [ ] Down migration written and tested locally

---

## Checklist
- [ ] Follows clean architecture — no layer imports above itself
- [ ] No raw SQL — squirrel used for all queries
- [ ] All errors use typed codes from `domain/errors/codes.go`
- [ ] No `fmt.Errorf` used where `appErrors` should be used
- [ ] `context.Context` threaded through all repository and usecase calls
- [ ] No env vars read outside `main.go` or config constructors
- [ ] No secrets or `.env` files committed
- [ ] `GET /healthz` returns `{"status":"ok","service":"<name>"}`
- [ ] `GET /metrics` endpoint present (Phase 4 — skip until observability sprint)

---

## Testing
- [ ] Tested locally with `curl` or equivalent
- [ ] No regressions on existing endpoints
- [ ] Unit tests added (if applicable)

---

## Logs
- [ ] All log calls use `sugar.Infow` / `sugar.Errorw` — no `Infof` or `Errorf`
- [ ] No sensitive data logged (passwords, tokens, secrets)
- [ ] `request_id` included in error logs

---

## Notes for reviewer
<!-- Anything worth calling out — tradeoffs made, things deferred, known issues -->