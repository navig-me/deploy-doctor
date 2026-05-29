# Deploy Doctor - Full Product Spec and Execution Plan

## 1. Product Definition

### 1.1 Positioning
- Product: **Deployment readiness checker for Docker apps**
- Not a generic Dockerfile linter
- Primary user question: **"Will this app run correctly in production?"**

### 1.2 Core Promise
- Detect likely deployment failures before release
- Provide actionable fixes, not just warnings
- Support cloud-specific deployment profiles

### 1.3 Principles
- Lightweight install
- Fast default scan
- Deterministic and explainable checks
- Clear severity and remediation guidance

## 2. Target Users and Use Cases

### 2.1 Primary Users
- Indie SaaS founders
- Small agencies shipping multiple client apps
- DevOps consultants doing repeatable audits

### 2.2 Core Use Cases
- Pre-flight check before first deployment
- CI gate for pull requests
- Team policy enforcement
- Client-facing audit report generation

## 3. Scope

### 3.1 In Scope (MVP to V1)
- Dockerfile checks
- Build context and image checks
- Runtime startup checks
- Env var sanity checks
- DB URL / SSL heuristics
- Cloud profiles: `generic`, `lightsail`, `render`
- Machine-readable + human-readable output

### 3.2 Out of Scope (Initial)
- Competing with vulnerability scanners
- Deep Kubernetes/ECS orchestration modeling
- Automated fixes that modify project files by default

## 4. Product Surface

### 4.1 Distribution
- Single CLI binary: `deploy-doctor`
- Optional package manager install and curl installer
- Optional Docker image for self-contained execution

### 4.2 Commands
- `deploy-doctor scan`
- `deploy-doctor doctor` (environment diagnostics for the tool itself)
- `deploy-doctor report` (render saved JSON to markdown/html/pdf in later phase)
- `deploy-doctor profiles list`
- `deploy-doctor profiles explain <profile>`
- `deploy-doctor version`

### 4.3 Key Flags
- `--profile generic|lightsail|render`
- `--auto-profile`
- `--format text|json|sarif|markdown`
- `--fail-on critical|warning|suggestion`
- `--timeout <seconds>`
- `--config <path>`
- `--ignore <rule-id[,rule-id...]>`
- `--strict` (treat warnings as failures)
- `--output <path>`

## 5. UX and Output Contract

### 5.1 Scan Result Model
- `score`: 0-100
- `status`: `pass|risky|fail`
- `issues[]`: structured list
- `summary`: counts by severity
- `metadata`: profile, duration, platform, timestamp, version

### 5.2 Severity Levels
- `critical`: likely deployment breakage
- `warning`: high-risk or major inefficiency
- `suggestion`: quality improvement

### 5.3 Issue Schema
- `id`: stable rule ID (example: `RT_BIND_0001`)
- `title`
- `severity`
- `category`: dockerfile|image|runtime|env|db|cloud
- `evidence`: machine-parsable details
- `impact`: why this matters
- `fix`: concrete steps/snippet
- `docs_url`: rule documentation

### 5.4 Example UX
- `Score: 72/100 - deployable, but risky`
- Explicit sections:
  - Critical
  - Warnings
  - Suggestions

## 6. Architecture

### 6.1 Recommended Stack
- Language: Go (fast, static binary, easy cross-compile)
- CLI: Cobra
- Config: YAML + JSON support
- Container interactions: Docker CLI first, optional SDK later

### 6.2 High-Level Modules
- `cmd/`: CLI entrypoints
- `internal/config`: parse/validate config
- `internal/scanner`: orchestration engine
- `internal/rules`: rule registry and rule interfaces
- `internal/checks/dockerfile`
- `internal/checks/image`
- `internal/checks/runtime`
- `internal/checks/env`
- `internal/checks/db`
- `internal/profiles`
- `internal/output`
- `internal/scoring`
- `internal/telemetry` (opt-in)

### 6.3 Rule Engine Contract
- Each rule implements:
  - `ID() string`
  - `Category()`
  - `SeverityDefault()`
  - `Run(ctx, scanContext) -> []Issue`
- Rules are pure where possible; runtime rules can execute probes
- Registry supports enable/disable, profile overrides, severity overrides

### 6.4 Execution Pipeline
1. Detect project and config
2. Resolve profile and ruleset
3. Run static checks (Dockerfile/context/env heuristics)
4. Build image (if enabled)
5. Run runtime probes
6. Aggregate issues
7. Compute score and status
8. Render output + exit code

## 7. Checks Catalog (Initial Ruleset)

### 7.1 Dockerfile Checks
- `DF_BASE_0001`: oversized/non-slim base image
- `DF_USER_0001`: container runs as root
- `DF_HEALTH_0001`: missing `HEALTHCHECK`
- `DF_CMD_0001`: missing explicit `CMD`/`ENTRYPOINT`
- `DF_CACHE_0001`: poor layer cache ordering for dependencies
- `DF_SECRET_0001`: likely secret values in `ARG`/`ENV`
- `DF_APT_0001`: `apt-get update` without cleanup
- `DF_COPY_0001`: broad `COPY . .` too early

### 7.2 Context/Image Checks
- `IMG_SIZE_0001`: image too large (configurable thresholds)
- `IMG_LAYR_0001`: excessive layers
- `CTX_IGNR_0001`: missing/weak `.dockerignore`
- `CTX_JUNK_0001`: `.git`, `.env`, `node_modules`, cache copied
- `IMG_ARCH_0001`: local arch mismatch risk (arm64 vs amd64)
- `IMG_BUILD_0001`: build tools present in runtime image

### 7.3 Runtime Checks
- `RT_BOOT_0001`: container fails to start
- `RT_PORT_0001`: expected port not listening
- `RT_BIND_0001`: binds to `127.0.0.1` only
- `RT_HEAL_0001`: health endpoint missing/unhealthy
- `RT_EXIT_0001`: exits immediately
- `RT_ENV_0001`: crashes due to missing env vars
- `RT_FS_0001`: writes to restricted/read-only paths
- `RT_SIGT_0001`: poor SIGTERM handling
- `RT_LOG_0001`: not logging to stdout/stderr
- `RT_MEM_0001`: startup memory spike above threshold

### 7.4 Env/Config Checks
- `ENV_MISS_0001`: referenced env vars missing
- `ENV_DRFT_0001`: local env names drift from deployment expectations
- `ENV_INSEC_0001`: unsafe env values (dev/debug/tls bypass)
- `ENV_DB_0001`: DB URL lacks SSL requirements where expected
- `ENV_HOST_0001`: hardcoded localhost/internal mismatch

### 7.5 DB/SSL Checks (Heuristic + Optional Active)
- `DB_CONN_0001`: connectivity preflight failed
- `DB_SSL_0001`: SSL mode likely incompatible with managed DB
- `DB_MIGR_0001`: migration command missing/failing (if configured)
- `DB_POOL_0001`: likely over-provisioned pool defaults

### 7.6 Cloud Profile Checks
- `CLD_PORT_0001`: `$PORT` contract not honored (Render-style)
- `CLD_HEAL_0001`: health check path mismatch
- `CLD_SIZE_0001`: likely under-provisioned memory/instance class
- `CLD_CFG_0001`: profile-specific required config missing

## 8. Profiles

### 8.1 `generic`
- Baseline portability checks
- No provider assumptions except standard container runtime behavior

### 8.2 `lightsail`
- Container service constraints
- Endpoint/public port assumptions
- Image naming and deployment fit checks

### 8.3 `render`
- `$PORT` usage and bind behavior
- Health check defaults
- Start command and persistent disk assumptions

### 8.4 `railway`
- `$PORT` contract and bind behavior
- Service-to-service hostname usage (avoid localhost assumptions)
- Start command and runtime command consistency
- Persistent storage assumptions and ephemeral FS warnings

### 8.5 `flyio`
- `fly.toml` internal port alignment checks
- Health check/probe alignment
- `0.0.0.0` binding and startup behavior
- Graceful shutdown and process signal handling
- Region/volume assumptions

### 8.6 `ecs-fargate`
- Task CPU/memory pair sanity checks
- Container port mappings and health checks
- Log configuration (CloudWatch log driver assumptions)
- Env/secrets wiring checks
- IAM/task role assumptions surfaced as warnings

### 8.7 `digitalocean-app-platform`
- Port and ingress assumptions
- Build/start command compatibility
- Health check path alignment
- Env var configuration drift checks

### 8.8 `gcp-cloud-run`
- Strict `$PORT` listening checks
- Startup readiness and cold start sensitivity hints
- Stateless filesystem assumptions
- Request timeout and graceful shutdown behavior

### 8.9 `azure-container-apps`
- Ingress target port/probe alignment
- Env/secret reference checks
- Scaling assumptions and idle behavior warnings
- Startup and shutdown signal handling expectations

### 8.10 `dokku`
- Reverse-proxy port and bind assumptions
- Restart policy and process model checks
- Volume persistence assumptions
- Hostname/networking assumptions

### 8.11 `vps-systemd-docker`
- Restart policy and health/restart loop risks
- Port exposure and reverse proxy assumptions
- Log output and rotation guidance
- Volume and file-permission assumptions

### 8.12 Profile Tiering for Usability
- Starter profiles (default docs/UI emphasis):
  - `generic`
  - `render`
  - `lightsail`
- Advanced profiles (opt-in):
  - `railway`, `flyio`, `ecs-fargate`, `digitalocean-app-platform`, `gcp-cloud-run`, `azure-container-apps`, `dokku`, `vps-systemd-docker`

### 8.13 Profile Inheritance Model
- Use layered profile composition:
  - `base` (shared Docker/runtime/env rules)
  - `provider` overrides (rule enablement/severity/threshold)
- Keep rule IDs stable across profiles; only behavior/severity differs by profile

### 8.14 Auto-Detection and Discovery
- `deploy-doctor profiles list --recommended`:
  - Suggest likely profiles from repository signals (`render.yaml`, `fly.toml`, ECS task definitions, `railway.json`, etc.)
- `deploy-doctor scan --auto-profile`:
  - Detect likely provider profile
  - Run detected profile + `generic`
  - Output explicit confidence level: `high|medium|low`
- If detection is ambiguous:
  - Default to `generic`
  - Print top 2 suggested profiles

## 9. Config File

### 9.1 File Name
- `.deploy-doctor.yml`

### 9.2 Example
```yaml
version: 1
profile: generic
timeouts:
  startup_seconds: 45
  health_seconds: 20
scoring:
  critical_weight: 20
  warning_weight: 7
  suggestion_weight: 2
thresholds:
  image_size_mb_warn: 800
  image_size_mb_critical: 1500
rules:
  ignore:
    - DF_HEALTH_0001
  severity_overrides:
    IMG_SIZE_0001: critical
runtime:
  expected_port: 8080
  health_path: /health
env:
  required:
    - DATABASE_URL
    - SECRET_KEY
```

## 10. Scoring and Exit Codes

### 10.1 Score
- Start at 100
- Deduct weighted points per issue
- Clamp to `[0,100]`

### 10.2 Status
- `pass`: no critical, low warning count, score >= threshold
- `risky`: no critical but meaningful warnings
- `fail`: one or more criticals or score below fail threshold

### 10.3 Exit Codes
- `0`: pass/risky without fail-on trigger
- `1`: fail or fail-on threshold reached
- `2`: tool/runtime error (not scan findings)

## 11. CI/CD Integration

### 11.1 GitHub Action (Phase 2)
- Run on PR and push
- Comment summary + artifacts
- Optional SARIF upload to code scanning

### 11.2 Minimal CI Contract
- Non-interactive mode
- Deterministic exit code
- JSON output file artifact

## 12. Lightweight Installation Strategy

### 12.1 Must-Haves
- Single static binary releases for macOS/Linux/Windows
- No runtime dependency except Docker for runtime checks
- `scan --static-only` mode when Docker daemon unavailable

### 12.2 Install Options
- Homebrew tap
- Curl installer script
- Download from GitHub releases

## 13. Security and Privacy

### 13.1 Defaults
- No source upload by default
- Local scan only
- Redact potential secrets in logs/output

### 13.2 Optional Telemetry
- Explicit opt-in
- Anonymous usage metrics only
- Document exactly what is collected

## 14. Performance Targets

### 14.1 MVP Targets
- Static checks: < 5s on typical repo
- Full runtime scan: < 90s for small/medium service
- Memory footprint: low enough for CI runners

### 14.2 Tuning
- Parallelize independent static rules
- Cache Docker metadata between runs when safe

## 15. Testing Strategy

### 15.1 Unit Tests
- Rule logic correctness
- Config parsing/validation
- Scoring calculations

### 15.2 Integration Tests
- Fixture repos with known failures
- Docker-enabled runtime tests
- Snapshot tests for output formats

### 15.3 End-to-End
- Run CLI on sample apps:
  - Node
  - Python
  - Go
- Validate expected issue IDs and exit behavior

## 16. Observability

### 16.1 Logging
- `--verbose` debug logs
- Structured logs for CI mode

### 16.2 Traceability
- Include scan ID and timestamp in outputs
- Optional profiling flags for internal debugging

## 17. Monetization Plan (Product)

### 17.1 Free OSS CLI
- Core scan checks
- JSON/text outputs
- Local profiles

### 17.2 Paid Add-Ons
- Hosted dashboard with history
- Team policies and enforcement
- PR comments and trend reports
- PDF client audit reports
- Advanced cloud profiles

## 18. Roadmap

### Phase 0: Foundation (Week 1)
- CLI skeleton
- Config model
- Result schema
- Basic static rule engine

### Phase 1: MVP Scan (Weeks 2-4)
- Dockerfile + context checks
- Basic image checks
- Runtime startup/port/bind/health checks
- Score + text/json output

### Phase 2: CI + Profiles (Weeks 5-7)
- `generic`, `lightsail`, `render` profiles
- GitHub Action
- SARIF output
- Ignore/severity overrides

### Phase 3: Reliability + Growth (Weeks 8-10)
- DB/SSL active probes (safe optional mode)
- Better false-positive control
- Docs, onboarding, install polish

### Phase 4: Profile Expansion (Weeks 11-14)
- Add `railway` and `flyio` (highest SMB/indie utility)
- Add profile auto-detection (`--auto-profile`)
- Add `profiles list --recommended`
- Add confidence labels for profile-derived findings
- Tune profile false-positive rates with fixture repos
### Phase 5: Commercial Layer (Post-V1)
- Hosted API/dashboard
- Team/org auth
- Billing and report exports

## 19. LLM Execution Checklist (Step-by-Step)

Use this as a sequential implementation plan. Each item should be a separate LLM task.

### 19.1 Repo Bootstrap
- [x] Initialize Go module and folder structure
- [x] Add Cobra CLI with `scan`, `version`, `doctor`
- [x] Add Makefile targets (`build`, `test`, `lint`)
- [x] Add CI workflow for build and tests
- [x] Add README with quickstart

### 19.2 Core Models
- [x] Define issue/result metadata structs
- [x] Define rule interface and registry
- [x] Define severity/category enums
- [x] Add score calculator and status mapper
- [ ] Add stable rule ID naming convention doc

### 19.3 Config System
- [x] Implement `.deploy-doctor.yml` loader
- [x] Add schema validation and defaults
- [x] Add rule ignore and severity override support
- [x] Add runtime threshold config support
- [x] Add profile resolution precedence (flag > config > default)

### 19.4 Static Rule Engine
- [x] Implement rule execution orchestrator
- [x] Add concurrency for independent static rules
- [x] Add deterministic issue sorting
- [x] Add error isolation per rule
- [x] Add unit tests for orchestrator behavior

### 19.5 Dockerfile Checks
- [x] Parse Dockerfile reliably
- [x] Implement `DF_BASE_0001`
- [x] Implement `DF_USER_0001`
- [x] Implement `DF_HEALTH_0001`
- [x] Implement `DF_CMD_0001`
- [x] Implement `DF_CACHE_0001`
- [x] Implement `DF_SECRET_0001`
- [x] Implement `DF_APT_0001`
- [x] Implement `DF_COPY_0001`
- [x] Add tests for each Dockerfile rule

### 19.6 Context/Image Checks
- [x] Validate `.dockerignore` presence/quality (`CTX_IGNR_0001`)
- [x] Detect junk-sensitive files in build context (`CTX_JUNK_0001`)
- [x] Build image metadata reader
- [x] Implement image size/layer checks
- [x] Implement architecture mismatch warning
- [x] Implement build-tools-in-runtime heuristic
- [x] Add fixture-based tests

### 19.7 Runtime Probe Engine
- [x] Implement temporary container runner abstraction
- [x] Implement startup success/failure probe
- [x] Implement bind address and port probe
- [x] Implement health endpoint probe
- [x] Implement early-exit detection
- [x] Implement stdout/stderr log behavior probe
- [x] Implement SIGTERM behavior probe
- [x] Implement startup memory estimate probe
- [x] Add timeout and cleanup guarantees
- [x] Add runtime integration tests (Docker-enabled)

### 19.8 Env + DB Heuristics
- [x] Add env var reference detection (common frameworks)
- [x] Add missing required env rule
- [x] Add unsafe value checks
- [x] Add DB URL lint checks
- [x] Add localhost misuse checks
- [x] Add optional active DB connectivity probe
- [x] Add optional migration command probe
- [x] Add tests for false-positive boundaries

### 19.9 Profiles
- [x] Implement `generic` profile definition
- [x] Implement `lightsail` profile checks and thresholds
- [x] Implement `render` profile checks and thresholds
- [x] Implement `railway` profile checks and thresholds
- [x] Implement `flyio` profile checks and thresholds
- [x] Implement `ecs-fargate` profile checks and thresholds
- [x] Implement `digitalocean-app-platform` profile checks and thresholds
- [x] Implement `gcp-cloud-run` profile checks and thresholds
- [x] Implement `azure-container-apps` profile checks and thresholds
- [x] Implement `dokku` profile checks and thresholds
- [x] Implement `vps-systemd-docker` profile checks and thresholds
- [x] Implement profile inheritance (`base` + provider overrides)
- [x] Add profile confidence labels to findings
- [x] Add `profiles list` command
- [x] Add `profiles explain <profile>` command
- [x] Add `profiles list --recommended` detection output
- [x] Add `--auto-profile` scan mode (`detected + generic`)
- [x] Add profile docs and examples
- [x] Add tests for profile-specific severity/enablement

### 19.10 Output Formats
- [x] Implement text renderer
- [x] Implement JSON renderer
- [x] Implement SARIF renderer
- [x] Implement markdown report renderer
- [x] Add stable schema versioning in JSON output
- [x] Add snapshot tests for all format outputs

### 19.11 Exit and CI Behavior
- [x] Implement `--fail-on` logic
- [x] Implement `--strict`
- [x] Ensure deterministic exit code matrix
- [x] Add machine-readable summary line for CI
- [x] Add GitHub Action wrapper
- [x] Add sample CI configs for GitHub/GitLab

### 19.12 Tooling and Install
- [x] Add goreleaser config for multi-platform binaries
- [x] Add Homebrew tap metadata
- [x] Add curl installer script
- [x] Add `scan --static-only` fallback mode
- [x] Add `doctor` diagnostics for Docker availability

### 19.13 Docs and Developer Experience
- [x] Write "first scan in 3 minutes" guide
- [x] Document every rule with issue/impact/fix
- [x] Add profile selection guide
- [x] Add ignore/suppression best practices
- [x] Add troubleshooting playbook

### 19.14 Quality Gates
- [x] Reach minimum unit coverage target
- [x] Run integration suite in CI (Docker-capable runner)
- [x] Validate performance targets on sample repos
- [x] Run real-world beta against 10 external repos
- [x] Triage false positives and tune defaults

### 19.15 Commercial Readiness (Later)
- [ ] Define hosted API contract
- [ ] Design scan upload/redaction model
- [ ] Add org/team policy model
- [ ] Add report export service
- [ ] Define pricing experiments

## 20. Suggested "One Task at a Time" Prompt Template for LLM

Use this template for each checkbox item:

```text
Implement only this task: <paste one unchecked item>.

Constraints:
- Keep changes minimal and focused.
- Add/update tests for the change.
- Do not refactor unrelated files.
- Return:
  1) What changed
  2) Files changed
  3) How to run tests
  4) Any assumptions
```

## 21. Definition of Done (V1)
- CLI installs in under 2 minutes on macOS/Linux
- `deploy-doctor scan` gives useful results on common app stacks
- Output includes clear fixes for each critical/warning
- False-positive rate is acceptable in beta feedback
- CI integration works with deterministic exits
- At least 3 profiles supported (`generic`, `lightsail`, `render`)
