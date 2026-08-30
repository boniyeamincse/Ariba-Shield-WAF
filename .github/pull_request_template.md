## Description

<!-- Describe the change and link to the related issue/task. -->

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactor / code quality
- [ ] CI / build / tooling
- [ ] Security fix

## Checklist

- [ ] I have read the [contributing guidelines](../docs/architecture/repository-ci-coding-standards.md)
- [ ] `make lint` passes
- [ ] `make test` passes on all affected packages
- [ ] `make gen-check` passes (if schemas changed)
- [ ] `make check-i18n` passes (if messages changed)
- [ ] New routes are reflected in `docs/api/openapi-v0.yaml`
- [ ] Architecture diagrams updated in `docs/architecture/` if trust boundary or data flow changed
- [ ] No secrets, credentials, or `.env` files committed
- [ ] Commit messages follow the project convention

## Related ADRs

<!-- Link to any ADRs affected by this change (e.g. ADR-002). -->