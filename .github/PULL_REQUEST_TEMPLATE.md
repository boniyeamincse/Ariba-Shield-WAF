## Description
Please include a summary of the changes and the related issue. Describe the architectural impacts (if any) and mention which ADR is tied to this PR if you are changing boundaries, stores, or flow.

Fixes # (issue)

## Type of change
Please delete options that are not relevant.
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Checklist:
- [ ] I have read the `CONTRIBUTING.md` guide.
- [ ] I have run `make lint` and it passes cleanly.
- [ ] I have run `make test` and all unit tests pass.
- [ ] If this introduces new API endpoints, I have updated `docs/api/endpoint.md` and `openapi-v0.yaml`.
- [ ] I have verified that there are no secrets, keys, or hardcoded sensitive credentials in my changes.
- [ ] I have verified that this code does NOT block or crash the live traffic path in the Gateway.
- [ ] I have provided a descriptive and imperative commit message.
