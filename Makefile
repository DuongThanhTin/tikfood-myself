# TikFood — verification entrypoints.
#
# Single source of truth for "how do I check my work before committing / opening a PR".
# The Definition of Done (docs/verification/definition-of-done.md) and the per-domain
# recipes (docs/recipes/*.md) point here instead of listing loose commands.
#
#   make verify         # run every app's gate (api + web + runner)
#   make verify-api     # go vet + go test + go build   (apps/api)
#   make verify-web     # typecheck + test + build       (apps/web)
#   make verify-runner  # tsc --noEmit + test            (apps/ai-code-runner)

# Keep the Go build cache inside the repo's gitignored .cache/ dir (matches package.json).
GOCACHE := $(CURDIR)/.cache/go-build

# Record a PASS into .verify-evidence/<surface> — only reached when the target's tools
# all exit 0. verify-evidence-guard.py reads these before allowing commit / PR.
RECORD := bash $(CURDIR)/.claude/hooks/record-verify.sh

.PHONY: verify verify-api verify-web verify-runner verify-claude verify-ai-os verify-hooks

## verify: run every app's verification gate
verify: verify-api verify-web verify-runner verify-claude verify-ai-os verify-hooks

## verify-api: go vet + go test + go build (apps/api)
verify-api:
	cd apps/api && GOCACHE=$(GOCACHE) go vet ./...
	cd apps/api && GOCACHE=$(GOCACHE) go test ./...
	cd apps/api && GOCACHE=$(GOCACHE) go build ./...
	@$(RECORD) api

## verify-web: static check + unit tests + production build (apps/web)
## Note: apps/web has no ESLint config (Next 16 dropped the built-in `next lint`),
## so `tsc --noEmit` (the `typecheck` script) is the static-analysis / lint gate here.
verify-web:
	npm --workspace apps/web run typecheck
	npm --workspace apps/web run test
	npm --workspace apps/web run build
	@$(RECORD) web

## verify-runner: tsc --noEmit + unit tests (apps/ai-code-runner)
verify-runner:
	npm --workspace apps/ai-code-runner run typecheck
	npm --workspace apps/ai-code-runner run test
	@$(RECORD) runner

## verify-claude: validate .claude/ artifacts (frontmatter + link resolution)
verify-claude:
	bash .claude/check-artifacts.sh

## verify-ai-os: AI-OS semantic integrity (skill->recipe anchors + recipe->code drift)
verify-ai-os:
	bash .claude/check-ai-os.sh

## verify-hooks: governance guard unit tests (blocked_commands + protected_paths enforcement)
verify-hooks:
	bash .claude/hooks/governance-guard.test.sh
	bash .claude/hooks/verify-evidence-guard.test.sh
	@$(RECORD) hooks
