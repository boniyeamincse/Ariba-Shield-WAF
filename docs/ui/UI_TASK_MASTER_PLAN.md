# Ariba Shield WAF — UI Task Master Plan

**Generated:** 2026-08-30  
**Source:** Full codebase analysis (console-web frontend + control-api backend)  
**Priority:** P0 = required for usable application, P1 = important production functionality, P2 = useful enhancement, P3 = polish/future

---

## 1. Project Analysis

### Repository Structure

```
shield-waf/
  apps/
    console-web/          # Next.js 15 App Router frontend
    control-api/          # Go 1.26 REST API backend
  services/               # waf-engine, policy-compiler, event-ingestor, notification-service
  docs/                   # API docs, architecture, SRS, UI, operations
```

### Frontend Technology

| Technology | Version | Purpose |
|---|---|---|
| Next.js | 15.5.24 | Framework (App Router) |
| React | 19.2.8 | UI library |
| TypeScript | 5.9.3 | Language |
| next-intl | 3.26.5 | Internationalization (en, bn) |
| Vitest | 3.2.7 | Test runner (no tests exist) |
| ESLint | 9.x | Linting |
| Prettier | 3.x | Code formatting |

### Backend API

- **349 registered routes** in `router.go` (Go 1.22 ServeMux)
- **54 database tables** across 18 migrations
- **Auth:** Session-cookie (`shield_session`), HTTP-only/Secure/SameSite
- **RBAC:** 7 roles (Super Admin, Platform Admin, Security Admin, App Owner, SOC Analyst, Auditor, Read Only)
- **Permissions:** 14 permission strings (app:read, gateway:write, system:admin, etc.)

---

## 2. Frontend Architecture

### Current State

```
apps/console-web/
  src/
    app/
      page.tsx                      # Redirects to /en
      layout.tsx                    # Root layout (imports wrong globals.css!)
      globals.css                   # Full design system (326 lines) — NEVER IMPORTED
      [locale]/
        page.tsx                    # Overview/Dashboard (mock data fallback)
        layout.tsx                  # Locale layout (no sidebar, no nav, no header)
        login/page.tsx              # Login (functional)
        applications/page.tsx       # Applications list (no sidebar, no CRUD)
        gateways/page.tsx           # Gateways list (no sidebar)
        security-events/page.tsx    # Security events list (raw fetch, no sidebar)
        users/page.tsx              # Users list (mock data fallback)
        audit/page.tsx              # Audit log (resource_type bug)
        (no policies, no integrations, no reports, no settings, no traffic)
    components/
      layout/Sidebar.tsx            # Navigation sidebar
      UserProfileWidget.tsx         # User avatar + logout
      CreateUserButton.tsx          # Create user modal + form
    context/
      AuthContext.tsx               # Auth state (React Context)
    lib/
      api.ts                        # API client (10 functions only)
    messages/
      en.json, bn.json              # i18n messages (55 keys each)
  middleware.ts                     # next-intl locale routing
  globals.css                       # Minimal styles (18 lines) — the one actually imported
```

### Critical Bugs

1. **CSS Import Bug:** Root layout imports `@/globals.css` → resolves to `src/globals.css` (18 lines, minimal). The full design system (326 lines with glassmorphism, sidebar, dashboard) is in `src/app/globals.css` which is **never imported**. Result: no styling on fresh build.
2. **API Port Mismatch:** `.env.local` has `NEXT_PUBLIC_API_BASE=http://localhost:8080` but backend listens on `:8443`. API calls fail.
3. **Audit Event `resource_type` Bug:** Backend returns `"resource"` (JSON key), frontend accesses `ev.resource_type` → undefined.
4. **Login Redirect Ignores Locale:** `router.push("/en")` hardcoded — Bangla users get redirected to English.
5. **6 Broken Nav Links:** `/policies`, `/integrations`, `/reports`, `/settings` have no page files. `/traffic`, `/security_events` have nav translations but no sidebar links.
6. **No Loading/Error/Empty States:** No `loading.tsx`, `error.tsx`, `not-found.tsx` files anywhere.
7. **Security Events Uses Raw `fetch()`:** Bypasses the API client (`api.ts`).

---

## 3. API Inventory

### Full Endpoint Summary

| Module | APIs | UI Exists | Handler Type |
|---|---|---|---|
| Health & Metrics | 2 | No | Public |
| Authentication | 9 | Partial (login) | Public + Auth |
| Applications | 10 | Partial (list) | Dedicated |
| Domains | 2 | No | Dedicated |
| Origins | 2 | No | Dedicated |
| Security Policies | 13 | No | Dedicated |
| Policy Versions | 5 | No | Dedicated |
| Policy Approvals | 5 | No | Dedicated |
| Rules | 9 | No | Dedicated |
| Rule Bundles | 6 | No | Dedicated |
| Gateways | 12 | Partial (list) | Dedicated |
| Security Events | 6 | Partial (list) | Dedicated |
| Audit Events | 3 | Partial (list) | Dedicated |
| IP Lists | 8 | No | Dedicated |
| Rate Limits | 4 | No | Dedicated |
| Virtual Servers | 3 | No | Dedicated |
| Backend Pools | 7 | No | Dedicated |
| Backend Nodes | 9 | No | Dedicated |
| Health Monitors | 3 | No | Dedicated |
| Routes | 5 | No | Generic CRUD |
| Config Versions | 2 | No | Dedicated |
| Config Validation | 1 | No | Dedicated |
| Traffic | 1 | No | Dedicated |
| Backups | 3 | No | Dedicated |
| Users | 4 | Partial (list) | Dedicated |
| **Groups** | 5 | No | Dedicated |
| **Roles** | 2 | No | Dedicated |
| **Permissions** | 1 | No | Generic (bug) |
| User-Role Assignment | 2 | No | Dedicated |
| Webhooks | 2 | No | Dedicated |
| Exceptions | 8 | No | Dedicated |
| Managed Rules | 2 | No | Dedicated |
| Custom Rules | 2 | No | Dedicated |
| Deployments | 2 | No | Dedicated |
| Certificates | 7 | No | Dedicated |
| TLS Profiles | 5 | No | Generic CRUD |
| Threat Intelligence | 12 | No | Generic + Dedicated |
| API Security | 4 | No | Generic CRUD |
| Bot Management | 12 | No | Generic + Dedicated |
| DLP | 4 | No | Generic CRUD |
| **Integrations** | 8 | No | Dedicated + Generic |
| IAM / SSO | 4 | No | Generic CRUD |
| Service Accounts | 3 | No | Generic CRUD |
| API Keys | 3 | No | Generic CRUD |
| Secrets | 3 | No | Generic CRUD |
| Gateway Clusters | 9 | No | Generic + Dedicated |
| **Incidents** | 11 | No | Dedicated |
| Automation | 4 | No | Generic CRUD |
| Clusters | 4 | No | Generic CRUD |
| Caching | 4 | No | Generic CRUD |
| Analytics | 2 | No | Dedicated |
| **Dashboard** | 8 | Partial (overview) | Dedicated |
| Organizations | 5 | No | Generic CRUD |
| Tenants | 5 | No | Generic CRUD |
| **Listeners** | 7 | No | Generic + Dedicated |
| **Sites** | 6 | No | Generic + Dedicated |
| GraphQL Security | 4 | No | Generic CRUD |
| Client-Side Protection | 4 | No | Generic CRUD |
| API Quotas | 4 | No | Generic CRUD |
| ML Baselines | 4 | No | Generic CRUD |
| **Notification Channels** | 6 | No | Dedicated |
| **Reports** | 9 | No | Dedicated |
| **Learning** | 9 | No | Dedicated |
| Network Protection | 4 | No | Generic CRUD |
| **System Settings** | 8 | No | Dedicated |
| **License / Entitlements** | 5 | No | Dedicated |

**Bold** = modules with dedicated handlers written in the backend but **zero** frontend implementation.

---

## 4. Existing UI Inventory

| Page | Route | Status | Sidebar | API Client | Loading | Empty | Error | Pagination | Filters | Mock Data |
|---|---|---|---|---|---|---|---|---|---|---|
| Login | /login | ✅ COMPLETE | No | Yes | Yes | N/A | Yes | N/A | N/A | No |
| Overview | / | ⚠️ PARTIAL | Yes | Yes | No | No | No | No | No | Yes |
| Applications | /applications | ⚠️ PARTIAL | No | Yes | No | Yes | No | No | No | No |
| Gateways | /gateways | ⚠️ PARTIAL | No | Yes | No | Yes | No | No | No | No |
| Security Events | /security-events | ⚠️ PARTIAL | No | Raw fetch | No | Yes | No | No | No | No |
| Users | /users | ⚠️ PARTIAL | Yes | Yes | No | Yes | Yes | No | No | Yes |
| Audit Log | /audit | ⚠️ PARTIAL | Yes | Raw fetch | No | Yes | Yes | No | No | No |

### Missing Pages (sidebar links that 404)

| Nav Link | API Module | Priority |
|---|---|---|
| Policies | Security Policies | P0 |
| Integrations | Integrations | P1 |
| Reports | Reports | P1 |
| Settings | System Settings | P1 |
| Traffic | Traffic | P2 |
| Security Events (sidebar) | Security Events | P2 |

---

## 5. UI Development Tasks

### Foundation Layer

---

## UI-001 — Fix CSS Import Bug

**Priority:** P0  
**Module:** Foundation  
**Type:** Bug Fix  
**Purpose:** The root layout imports `@/globals.css` which resolves to `src/globals.css` (18 lines, minimal). The full design system is in `src/app/globals.css` which is never imported. Fix the import so the design system applies.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/src/app/layout.tsx`
- `apps/console-web/src/globals.css`
- `apps/console-web/src/app/globals.css`

**Dependencies:** None

**Implementation Requirements:**
- Change the import in `layout.tsx` from `@/globals.css` to `../globals.css` (or consolidate the two CSS files)
- Verify `.dashboard-container`, `.sidebar`, `.glass-panel`, `.main-content`, `.metrics-grid` classes are applied on all pages

**Acceptance Criteria:**
- [ ] Overview page renders with full design system styling
- [ ] Users page renders with sidebar and glassmorphism
- [ ] Audit page renders with dashboard layout
- [ ] No CSS class mismatches

**Testing:** Visual inspection of every page.

---

## UI-002 — Fix API Port Mismatch

**Priority:** P0  
**Module:** Foundation  
**Type:** Bug Fix  
**Purpose:** `.env.local` has `NEXT_PUBLIC_API_BASE=http://localhost:8080` but the backend listens on `:8443`. Fix the port so frontend can reach the API.

**Related APIs:** All API endpoints

**Existing Files:**
- `apps/console-web/.env.local`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** None

**Implementation Requirements:**
- Update `.env.local` to use port 8443
- Verify the fallback in `api.ts` also uses the correct port

**Acceptance Criteria:**
- [ ] Frontend can successfully call `GET /api/v1/health`
- [ ] Frontend can successfully call `POST /api/v1/auth/login`
- [ ] All API client functions work

**Testing:** Manual API call verification.

---

## UI-003 — Fix Audit Event `resource_type` Field

**Priority:** P0  
**Module:** Foundation  
**Type:** Bug Fix  
**Purpose:** Backend returns audit events with `"resource"` JSON key, but frontend accesses `ev.resource_type` which is undefined.

**Related APIs:** `GET /api/v1/audit-events`

**Existing Files:**
- `apps/console-web/src/app/[locale]/audit/page.tsx`

**Dependencies:** None

**Implementation Requirements:**
- Change the frontend type and template to use `ev.resource` instead of `ev.resource_type`

**Acceptance Criteria:**
- [ ] Audit page shows the resource value correctly
- [ ] No undefined values in the audit table

**Testing:** Visual inspection of audit page.

---

## UI-004 — Fix Login Redirect Locale

**Priority:** P0  
**Module:** Foundation  
**Type:** Bug Fix  
**Purpose:** Login page hardcodes `router.push("/en")` instead of using the current locale. Bangla users get redirected to English.

**Related APIs:** `POST /api/v1/auth/login`

**Existing Files:**
- `apps/console-web/src/app/[locale]/login/page.tsx`

**Dependencies:** None

**Implementation Requirements:**
- Use `useParams()` or `usePathname()` to get the current locale
- Redirect to `/${locale}` instead of `/en`

**Acceptance Criteria:**
- [ ] Login redirects to `/en` when on `/en/login`
- [ ] Login redirects to `/bn` when on `/bn/login`

**Testing:** Manual locale switching.

---

## UI-005 — Add `loading.tsx` and `error.tsx` Pages

**Priority:** P0  
**Module:** Foundation  
**Type:** Page  
**Purpose:** No pages have loading/error states. Users see a blank page until data is fetched. Add loading and error boundary files.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/src/app/[locale]/`

**Dependencies:** None

**Implementation Requirements:**
- Add `loading.tsx` in the `[locale]` directory with a simple loading skeleton
- Add `error.tsx` in the `[locale]` directory with a retry button
- Add `not-found.tsx` in the `[locale]` directory with a 404 message

**Acceptance Criteria:**
- [ ] Loading state shows during page transitions
- [ ] Error state shows with retry option on API failure
- [ ] 404 page shows for non-existent routes

**Testing:** Navigate between pages, trigger API errors.

---

### Application Shell

---

## UI-006 — Create Dashboard Layout Component

**Priority:** P0  
**Module:** Application Shell  
**Type:** Component  
**Purpose:** Currently, each page manually renders the sidebar + header + main-content structure. Some pages (applications, gateways, security-events) lack the sidebar entirely. Create a reusable `DashboardLayout` component.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/src/app/[locale]/layout.tsx`
- `apps/console-web/src/components/layout/Sidebar.tsx`
- `apps/console-web/src/app/[locale]/page.tsx`
- `apps/console-web/src/app/[locale]/users/page.tsx`
- `apps/console-web/src/app/[locale]/audit/page.tsx`

**Dependencies:** UI-001 (CSS fix)

**Implementation Requirements:**
- Create `src/components/layout/DashboardLayout.tsx` with sidebar + header + main-content
- Wrap pages that need the dashboard layout
- Add a proper header component with page title, user profile, and breadcrumbs
- Move the sidebar and header into the layout

**Acceptance Criteria:**
- [ ] All pages consistently have the sidebar
- [ ] All pages consistently have the header
- [ ] Applications, gateways, security-events pages get the sidebar
- [ ] Login page is NOT wrapped in the dashboard layout

**Testing:** Navigate between all pages, verify layout consistency.

---

## UI-007 — Fix Sidebar Navigation Links

**Priority:** P0  
**Module:** Application Shell  
**Type:** Component  
**Purpose:** 6 sidebar navigation links point to non-existent routes. Add the missing routes or hide links until pages exist.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/src/components/layout/Sidebar.tsx`
- `apps/console-web/src/messages/en.json`
- `apps/console-web/src/messages/bn.json`

**Dependencies:** UI-006

**Implementation Requirements:**
- Add nav links for `/traffic` and `/security-events` in the sidebar
- Keep `/policies`, `/integrations`, `/reports`, `/settings` links but route to placeholder pages
- Ensure all nav links have proper i18n keys

**Acceptance Criteria:**
- [ ] Sidebar matches the actual available pages
- [ ] No broken nav links
- [ ] Active state highlights correctly
- [ ] i18n keys exist for all nav items

**Testing:** Click every sidebar link, verify navigation.

---

### Shared Components

---

## UI-008 — Create DataTable Component

**Priority:** P0  
**Module:** Shared Components  
**Type:** Component  
**Purpose:** Every list page (applications, gateways, security-events, users, audit, etc.) needs a table. Create a reusable `DataTable` component with sorting, pagination, and column configuration.

**Related APIs:** All list endpoints

**Existing Files:** None

**Dependencies:** None

**Implementation Requirements:**
- Create `src/components/shared/DataTable.tsx`
- Support: columns, rows, sorting by column, pagination, loading state, empty state, error state
- Support: custom cell renderers, row click handlers, selection
- Support: `onSort`, `onPageChange`, `onPageSizeChange` callbacks
- Use TypeScript generics for type safety

**Acceptance Criteria:**
- [ ] DataTable renders with proper column headers
- [ ] Sorting works (clickable column headers)
- [ ] Pagination controls work (page numbers, prev/next, page size)
- [ ] Loading state shows skeleton rows
- [ ] Empty state shows a message
- [ ] Error state shows a retry button

**Testing:** Use in multiple pages, verify all states.

---

## UI-009 — Create StatusBadge and SeverityBadge Components

**Priority:** P0  
**Module:** Shared Components  
**Type:** Component  
**Purpose:** Many pages display status (active/inactive/blocked) and severity (critical/high/medium/low). Create reusable badge components with consistent styling.

**Related APIs:** All list endpoints

**Existing Files:** None (hardcoded maps in pages)

**Dependencies:** None

**Implementation Requirements:**
- Create `src/components/shared/StatusBadge.tsx` — colored badge for status values
- Create `src/components/shared/SeverityBadge.tsx` — colored badge for severity levels
- Use consistent color schemes (green=active, red=blocked, yellow=warning, etc.)
- Support customizable colors and labels

**Acceptance Criteria:**
- [ ] StatusBadge renders correctly for all status values
- [ ] SeverityBadge renders correctly for all severity levels
- [ ] Badges are accessible (proper contrast, aria-labels)

**Testing:** Use in multiple pages, verify visual consistency.

---

## UI-010 — Create ConfirmDialog Component

**Priority:** P0  
**Module:** Shared Components  
**Type:** Component  
**Purpose:** Delete operations (applications, policies, users, etc.) need a confirmation dialog. Create a reusable `ConfirmDialog` component.

**Related APIs:** All DELETE endpoints

**Existing Files:** None

**Dependencies:** None

**Implementation Requirements:**
- Create `src/components/shared/ConfirmDialog.tsx`
- Support: title, message, confirm button, cancel button, loading state
- Support: customizable button text (e.g., "Delete", "Disable", "Archive")
- Support: danger variant (red confirm button)
- Use modal overlay with backdrop click to close

**Acceptance Criteria:**
- [ ] Dialog opens with proper title and message
- [ ] Confirm button triggers the action
- [ ] Cancel button closes the dialog
- [ ] Loading state shows while action is in progress
- [ ] Keyboard accessible (Escape to close, Enter to confirm)

**Testing:** Use in delete flows, verify all states.

---

## UI-011 — Create PageHeader Component

**Priority:** P0  
**Module:** Shared Components  
**Type:** Component  
**Purpose:** Every page needs a consistent header with title, description, and action buttons. Create a reusable `PageHeader` component.

**Related APIs:** None

**Existing Files:** None (inline headers in each page)

**Dependencies:** None

**Implementation Requirements:**
- Create `src/components/shared/PageHeader.tsx`
- Support: title, description, action buttons (array of {label, onClick, variant})
- Support: breadcrumbs
- Support: i18n for title and description

**Acceptance Criteria:**
- [ ] PageHeader renders with title and description
- [ ] Action buttons render correctly
- [ ] All variants are styled consistently

**Testing:** Use in all pages, verify consistency.

---

## UI-012 — Create FilterBar Component

**Priority:** P1  
**Module:** Shared Components  
**Type:** Component  
**Purpose:** Many list endpoints support filters (severity, status, kind, date range). Create a reusable `FilterBar` component.

**Related APIs:** All list endpoints with query parameters

**Existing Files:** None

**Dependencies:** UI-008 (DataTable)

**Implementation Requirements:**
- Create `src/components/shared/FilterBar.tsx`
- Support: text search, select dropdowns, date range picker, clear all
- Support: configurable filters (array of {key, label, type, options})
- Support: `onFilterChange` callback

**Acceptance Criteria:**
- [ ] FilterBar renders with configured filters
- [ ] Changing a filter updates the data
- [ ] Clear all resets all filters
- [ ] Filters persist in URL query params

**Testing:** Use in security-events, audit, incidents pages.

---

### Authentication Pages

---

## UI-013 — Add MFA Enable/Verify/Disable Pages

**Priority:** P1  
**Module:** Authentication  
**Type:** Page  
**Purpose:** Backend supports MFA (TOTP) enable, verify, and disable. Create UI pages for MFA management.

**Related APIs:**
- `POST /api/v1/auth/mfa/enable`
- `POST /api/v1/auth/mfa/verify`
- `POST /api/v1/auth/mfa/disable`

**Existing Files:**
- `apps/console-web/src/app/[locale]/login/page.tsx`
- `apps/console-web/src/context/AuthContext.tsx`

**Dependencies:** UI-006 (Dashboard Layout)

**Implementation Requirements:**
- Add MFA page accessible from user profile or settings
- Enable: show QR code (TOTP secret), allow user to scan
- Verify: 6-digit code input, verify before enabling
- Disable: confirm dialog, then disable
- Show MFA status in user profile

**Acceptance Criteria:**
- [ ] MFA enable page shows QR code
- [ ] MFA verify accepts 6-digit code
- [ ] MFA disable works with confirmation
- [ ] MFA status is shown in user profile

**Testing:** Enable MFA, verify with authenticator app, disable.

---

### Dashboard Pages

---

## UI-014 — Replace Dashboard Mock Data with Real API Calls

**Priority:** P0  
**Module:** Dashboard  
**Type:** Integration  
**Purpose:** The Overview page falls back to hardcoded mock data when API is unreachable. Remove mock data and use real API calls with proper error states.

**Related APIs:**
- `GET /api/v1/dashboard/overview`
- `GET /api/v1/dashboard/applications`

**Existing Files:**
- `apps/console-web/src/app/[locale]/page.tsx`

**Dependencies:** UI-001 (CSS fix), UI-005 (loading/error states)

**Implementation Requirements:**
- Add `getDashboardOverview()` and `getDashboardApplications()` to the API client
- Replace hardcoded metric labels ("100% Protected", "1 Node Offline", "All Synced") with real API responses
- Add loading state (skeleton cards)
- Add error state with retry option
- Add empty state for when no data exists

**Acceptance Criteria:**
- [ ] Metric cards show real API data
- [ ] Loading state shows skeleton cards
- [ ] Error state shows retry button
- [ ] No hardcoded mock data remains

**Testing:** Verify with live API, test loading/error states.

---

## UI-015 — Create Dashboard Traffic Widget

**Priority:** P1  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/traffic` with request volume, latency, and status distribution. Create a traffic widget.

**Related APIs:**
- `GET /api/v1/dashboard/traffic`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add `getDashboardTraffic()` to API client
- Create a traffic summary card showing total requests, avg latency, p99 latency
- Create a status distribution chart (e.g., bar chart or pie chart)
- Use a simple chart library (e.g., recharts or chart.js) or CSS-only bars

**Acceptance Criteria:**
- [ ] Traffic card shows real data from API
- [ ] Status distribution chart renders correctly
- [ ] Units are clear (requests, ms)
- [ ] Responsive layout

**Testing:** Verify with live API.

---

## UI-016 — Create Dashboard Security Widget

**Priority:** P1  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/security` with event volume, blocked count, unique IPs, and severity distribution. Create a security widget.

**Related APIs:**
- `GET /api/v1/dashboard/security`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add `getDashboardSecurity()` to API client
- Create security summary card (total events, blocked, unique IPs)
- Create severity distribution chart (colored bars)
- Add a link to full security events page

**Acceptance Criteria:**
- [ ] Security card shows real data from API
- [ ] Severity distribution chart renders correctly
- [ ] Link to security events page works

**Testing:** Verify with live API.

---

## UI-017 — Create Dashboard Incidents Widget

**Priority:** P1  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/overview` that includes active incident count. Create an incidents widget.

**Related APIs:**
- `GET /api/v1/dashboard/overview`
- `GET /api/v1/incidents`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add incident count to the dashboard
- Show active incidents as a metric card
- Add a link to full incidents page

**Acceptance Criteria:**
- [ ] Incident count shows real data from API
- [ ] Link to incidents page works

**Testing:** Verify with live API.

---

## UI-018 — Create Dashboard Top IPs Widget

**Priority:** P2  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/top-ips`. Create a widget showing top client IPs by event volume.

**Related APIs:**
- `GET /api/v1/dashboard/top-ips`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add `getDashboardTopIPs()` to API client
- Create a table showing top IPs with hits and blocked counts
- Limit to top 10 (backend returns top 20)

**Acceptance Criteria:**
- [ ] Top IPs table shows real data from API
- [ ] IPs are clickable (link to security events filtered by IP)

**Testing:** Verify with live API.

---

## UI-019 — Create Dashboard Top Rules Widget

**Priority:** P2  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/top-rules`. Create a widget showing top triggered rules.

**Related APIs:**
- `GET /api/v1/dashboard/top-rules`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add `getDashboardTopRules()` to API client
- Create a table showing top rules with hit counts
- Limit to top 10

**Acceptance Criteria:**
- [ ] Top rules table shows real data from API
- [ ] Rules are clickable (link to rules detail)

**Testing:** Verify with live API.

---

## UI-020 — Create Dashboard Gateway Health Widget

**Priority:** P2  
**Module:** Dashboard  
**Type:** Component  
**Purpose:** Backend has `GET /api/v1/dashboard/gateways`. Create a widget showing gateway fleet status.

**Related APIs:**
- `GET /api/v1/dashboard/gateways`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-014

**Implementation Requirements:**
- Add `getDashboardGateways()` to API client
- Create a gateway status card (total/active/offline)
- Replace hardcoded "1 Node Offline" with real data

**Acceptance Criteria:**
- [ ] Gateway status card shows real data from API
- [ ] No hardcoded metrics remain

**Testing:** Verify with live API.

---

### Core Resource Pages

---

## UI-021 — Complete Applications Page

**Priority:** P0  
**Module:** Applications  
**Type:** Page  
**Purpose:** The applications page is a bare list with no sidebar, no create/edit/delete, no detail view, no search, no filters. Add full CRUD UI.

**Related APIs:**
- `GET /api/v1/applications`
- `POST /api/v1/applications`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `DELETE /api/v1/applications/{id}`
- `GET /api/v1/applications/{id}/domains`
- `GET /api/v1/applications/{id}/origins`
- `GET /api/v1/applications/{id}/traffic`
- `GET /api/v1/applications/{id}/events`
- `GET /api/v1/applications/{id}/incidents`
- `GET /api/v1/applications/{id}/policies`
- `GET /api/v1/applications/{id}/health`

**Existing Files:**
- `apps/console-web/src/app/[locale]/applications/page.tsx`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011, UI-012

**Implementation Requirements:**
- **List:** Use DataTable with columns (name, status, domains, origins, health)
- **Create:** Application form (name, description) + POST /applications
- **Detail:** Application detail page with tabs (overview, domains, origins, events, incidents, policies, health)
- **Edit:** Inline edit or edit form + PATCH /applications/{id}
- **Delete:** ConfirmDialog + DELETE /applications/{id}
- Add filters by status, search by name
- Add pagination (backend supports default 50, but frontend can limit)

**Acceptance Criteria:**
- [ ] Application list shows all applications with proper columns
- [ ] Create application form works with validation
- [ ] Application detail page shows all tabs with real data
- [ ] Edit application works
- [ ] Delete application works with confirmation
- [ ] Search and filters work
- [ ] Loading, empty, error states handled

**Testing:** Full CRUD flow, detail page navigation, verify all tabs.

---

## UI-022 — Create Application Form

**Priority:** P0  
**Module:** Applications  
**Type:** Form  
**Purpose:** Create a form for creating and editing applications. Handles POST and PATCH.

**Related APIs:**
- `POST /api/v1/applications`
- `PATCH /api/v1/applications/{id}`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-021

**Implementation Requirements:**
- Fields: name (required, text), description (optional, textarea)
- Validation: name required, max length
- Create mode: sends POST, shows success toast, redirects to list
- Edit mode: pre-fills form, sends PATCH, shows success toast
- Loading state on submit button
- Error state for API errors
- Cancel button returns to list

**Acceptance Criteria:**
- [ ] Create mode: POSTs new application, shows success
- [ ] Edit mode: pre-fills, PATCHes changes, shows success
- [ ] Validation: name required, max length enforced
- [ ] Loading state: button shows spinner
- [ ] Error state: API errors shown

**Testing:** Create and edit applications, verify API calls.

---

## UI-023 — Create Application Detail Page

**Priority:** P1  
**Module:** Applications  
**Type:** Page  
**Purpose:** Create a detail page for an application with tabs showing different data.

**Related APIs:**
- `GET /api/v1/applications/{id}`
- `GET /api/v1/applications/{id}/domains`
- `GET /api/v1/applications/{id}/origins`
- `GET /api/v1/applications/{id}/traffic`
- `GET /api/v1/applications/{id}/events`
- `GET /api/v1/applications/{id}/incidents`
- `GET /api/v1/applications/{id}/policies`
- `GET /api/v1/applications/{id}/health`

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-021

**Implementation Requirements:**
- Page route: `/applications/{id}`
- Tabs: Overview, Domains, Origins, Traffic, Events, Incidents, Policies, Health
- Overview tab: application info, health status, quick stats
- Domains tab: list of domains (create/delete)
- Origins tab: list of origins (create/delete)
- Traffic tab: traffic data from dashboard API
- Events tab: security events for this application
- Incidents tab: incidents for this application
- Policies tab: policies bound to this application

**Acceptance Criteria:**
- [ ] All tabs load real data from API
- [ ] Navigating between tabs works
- [ ] Application not found shows 404
- [ ] Loading states for each tab

**Testing:** Navigate to application detail, click all tabs.

---

## UI-024 — Complete Gateways Page

**Priority:** P0  
**Module:** Gateways  
**Type:** Page  
**Purpose:** The gateways page is a bare list with no sidebar, no detail view, no actions. Add full CRUD UI.

**Related APIs:**
- `GET /api/v1/gateways`
- `GET /api/v1/gateways/{id}`
- `PATCH /api/v1/gateways/{id}`
- `DELETE /api/v1/gateways/{id}`
- `GET /api/v1/gateways/{id}/config`
- `GET /api/v1/gateways/{id}/status`
- `GET /api/v1/gateways/{id}/metrics`
- `POST /api/v1/gateways/{id}/config/apply`
- `POST /api/v1/gateways/{id}/config/rollback`

**Existing Files:**
- `apps/console-web/src/app/[locale]/gateways/page.tsx`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** Use DataTable (hostname, status, version, last_seen)
- **Detail:** Gateway detail page with tabs (overview, config, metrics, heartbeats)
- **Edit:** Edit gateway (hostname, ip)
- **Delete:** ConfirmDialog + DELETE
- **Actions:** Apply config, rollback config
- Add filters by status, search by hostname

**Acceptance Criteria:**
- [ ] Gateway list shows all gateways
- [ ] Detail page shows config, status, metrics
- [ ] Apply/rollback config works
- [ ] Delete works with confirmation
- [ ] Loading, empty, error states

**Testing:** List, detail, apply config, delete.

---

## UI-025 — Complete Security Events Page

**Priority:** P0  
**Module:** Security Events  
**Type:** Page  
**Purpose:** The security events page is a bare list with no sidebar, no pagination, no filters, and uses raw fetch(). Add full UI with pagination, filters, and detail view.

**Related APIs:**
- `GET /api/v1/security-events` (supports `?limit=&offset=&severity=&ip=&application_id=&from=&to=&action=&rule_id=`)
- `GET /api/v1/security-events/{id}`
- `GET /api/v1/security-events/{id}/matches`
- `GET /api/v1/security-events/{id}/timeline`
- `POST /api/v1/security-events/{id}/mask`
- `POST /api/v1/security-events/{id}/export`

**Existing Files:**
- `apps/console-web/src/app/[locale]/security-events/page.tsx`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-008, UI-009, UI-012

**Implementation Requirements:**
- **List:** Use DataTable (severity, time, method, path, reason, client_ip)
- **Pagination:** Use backend `?limit=&offset=`
- **Filters:** Severity (dropdown), IP (text), date range, action (dropdown), rule_id
- **Detail:** Event detail page with timeline, matches, mask action, export
- **Mask:** Button to mask/suppress an event
- **Export:** Button to export event
- Fix raw `fetch()` to use API client

**Acceptance Criteria:**
- [ ] Security events list uses API client, not raw fetch
- [ ] Pagination works (prev/next, page numbers)
- [ ] Filters work (severity, IP, date range, action, rule_id)
- [ ] Detail page shows event info, timeline, matches
- [ ] Mask action works
- [ ] Export action works
- [ ] Loading, empty, error states

**Testing:** List with filters, pagination, detail view, mask, export.

---

## UI-026 — Create Policy Pages

**Priority:** P1  
**Module:** Security Policies  
**Type:** Page  
**Purpose:** No UI exists for security policies. Need list, create, detail, edit, and lifecycle actions (validate, activate, disable, rollback, clone, version history, diff).

**Related APIs:**
- `GET /api/v1/security-policies`
- `POST /api/v1/security-policies`
- `GET /api/v1/security-policies/{id}`
- `PATCH /api/v1/security-policies/{id}`
- `DELETE /api/v1/security-policies/{id}`
- `POST /api/v1/security-policies/{id}/validate`
- `POST /api/v1/security-policies/{id}/activate`
- `POST /api/v1/security-policies/{id}/disable`
- `POST /api/v1/security-policies/{id}/rollback`
- `POST /api/v1/security-policies/{id}/clone`
- `GET /api/v1/security-policies/{id}/versions`
- `GET /api/v1/security-policies/{id}/diff`
- `POST /api/v1/security-policies/bind`
- `GET /api/v1/policy-versions/diff`
- `GET /api/v1/policy-approvals`
- `POST /api/v1/policy-approvals`
- `POST /api/v1/policy-approvals/{id}/approve`
- `POST /api/v1/policy-approvals/{id}/reject`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011, UI-012

**Implementation Requirements:**
- **List:** DataTable (name, status, enforcement_mode, version, application)
- **Create:** Policy form (name, enforcement_mode, optional application_id)
- **Detail:** Policy detail with tabs (overview, versions, diff, approvals, bindings)
- **Edit:** Edit policy
- **Delete:** ConfirmDialog
- **Actions:** Validate, Activate, Disable, Rollback, Clone
- **Version history:** List of versions with diff viewer
- **Approvals:** List, create, approve, reject approval requests
- **Bind:** Bind policy to application

**Acceptance Criteria:**
- [ ] Policy list shows all policies
- [ ] Create/Edit policy works
- [ ] Activate/Disable/Rollback/Clone actions work
- [ ] Version history shows all versions
- [ ] Diff viewer shows changes between versions
- [ ] Approval workflow works (create, approve, reject)
- [ ] Policy bind works

**Testing:** Full policy lifecycle: create → validate → activate → modify → version → rollback → delete.

---

## UI-027 — Create Incidents Page

**Priority:** P1  
**Module:** Incidents  
**Type:** Page  
**Purpose:** Backend has full incident management (11 endpoints). No UI exists.

**Related APIs:**
- `GET /api/v1/incidents`
- `POST /api/v1/incidents`
- `GET /api/v1/incidents/{id}`
- `PATCH /api/v1/incidents/{id}`
- `DELETE /api/v1/incidents/{id}`
- `POST /api/v1/incidents/{id}/assign`
- `POST /api/v1/incidents/{id}/escalate`
- `POST /api/v1/incidents/{id}/close`
- `POST /api/v1/incidents/{id}/reopen`
- `GET /api/v1/incidents/{id}/events`
- `GET /api/v1/incidents/{id}/timeline`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (title, severity, status, owner, created_at)
- **Create:** Incident form (title, severity, notes, related_events)
- **Detail:** Incident detail with timeline, related events, actions
- **Actions:** Assign (select owner), Escalate, Close, Reopen
- **Timeline:** Status history view
- **Events:** Related security events list
- Filters by severity, status, date range

**Acceptance Criteria:**
- [ ] Incident list shows all incidents
- [ ] Create incident works
- [ ] Detail page shows timeline and related events
- [ ] Assign/Escalate/Close/Reopen actions work
- [ ] Loading, empty, error states

**Testing:** Full incident lifecycle: create → assign → escalate → close → reopen.

---

### Operations Pages

---

## UI-028 — Create Audit Log Page with Pagination and Export

**Priority:** P1  
**Module:** Audit Log  
**Type:** Page  
**Purpose:** The audit log page exists but has no pagination, no filters, and has a `resource_type` bug. Upgrade it.

**Related APIs:**
- `GET /api/v1/audit-events`
- `GET /api/v1/audit-events/{id}`
- `GET /api/v1/audit-events/export`

**Existing Files:**
- `apps/console-web/src/app/[locale]/audit/page.tsx`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-008, UI-009, UI-012, UI-003 (resource_type fix)

**Implementation Requirements:**
- Fix `resource_type` → `resource` bug
- Add pagination
- Add filters by action, date range
- Add export button (CSV download)
- Add detail view for single audit event
- Fix raw `fetch()` to use API client

**Acceptance Criteria:**
- [ ] Audit events show with correct resource field
- [ ] Pagination works
- [ ] Filters work (action, date range)
- [ ] Export button downloads CSV
- [ ] Detail view shows single event
- [ ] Loading, empty, error states

**Testing:** List with filters, pagination, export, detail view.

---

## UI-029 — Create Users Page with Edit and Delete

**Priority:** P1  
**Module:** Users  
**Type:** Page  
**Purpose:** The users page exists but has mock data fallback, no edit/delete, no search. Upgrade it.

**Related APIs:**
- `GET /api/v1/users`
- `POST /api/v1/users`
- `PATCH /api/v1/users/{id}`
- `DELETE /api/v1/users/{id}`
- `POST /api/v1/users/{id}/roles`
- `DELETE /api/v1/users/{id}/roles/{roleId}`
- `GET /api/v1/roles`
- `GET /api/v1/groups`
- `POST /api/v1/groups`
- `GET /api/v1/groups/{id}`
- `PATCH /api/v1/groups/{id}`
- `DELETE /api/v1/groups/{id}`

**Existing Files:**
- `apps/console-web/src/app/[locale]/users/page.tsx`
- `apps/console-web/src/components/CreateUserButton.tsx`
- `apps/console-web/src/lib/api.ts`

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- Remove mock data fallback — use real API with error state
- **List:** DataTable (email, role, status, created_at, actions)
- **Create:** User form (email, password, role) — already exists as CreateUserButton
- **Edit:** Edit user (role, status)
- **Delete:** ConfirmDialog + DELETE
- **Role assignment:** POST/DELETE /users/{id}/roles
- **Groups:** List, create, edit, delete groups
- **Roles:** List roles
- **Permissions:** View permissions
- Search by email, filter by role/status

**Acceptance Criteria:**
- [ ] User list shows real data from API
- [ ] Create user works (already exists)
- [ ] Edit user role works
- [ ] Delete user works with confirmation
- [ ] Groups CRUD works
- [ ] Roles list shows all roles
- [ ] Permissions view works
- [ ] No mock data
- [ ] Loading, empty, error states

**Testing:** Full user lifecycle, group management, role/permission viewing.

---

### Missing Pages

---

## UI-030 — Create Integrations Page

**Priority:** P1  
**Module:** Integrations  
**Type:** Page  
**Purpose:** Backend has full integrations CRUD + test/enable/disable. No UI exists.

**Related APIs:**
- `GET /api/v1/integrations`
- `POST /api/v1/integrations`
- `GET /api/v1/integrations/{id}`
- `PATCH /api/v1/integrations/{id}`
- `DELETE /api/v1/integrations/{id}`
- `POST /api/v1/integrations/{id}/test`
- `POST /api/v1/integrations/{id}/enable`
- `POST /api/v1/integrations/{id}/disable`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (name, type, endpoint, enabled, status)
- **Create:** Integration form (type dropdown, name, endpoint, log_types, config)
- **Detail:** Integration detail with config view
- **Edit:** Edit integration
- **Delete:** ConfirmDialog
- **Actions:** Test, Enable, Disable
- Supported types: splunk_hec, wazuh, syslog, webhook, teams, slack

**Acceptance Criteria:**
- [ ] Integration list shows all integrations
- [ ] Create integration works with type-specific fields
- [ ] Test action works
- [ ] Enable/Disable works
- [ ] Loading, empty, error states

**Testing:** Create each integration type, test, enable/disable, delete.

---

## UI-031 — Create Reports Page

**Priority:** P1  
**Module:** Reports  
**Type:** Page  
**Purpose:** Backend has full report generation (security, traffic, incidents, compliance) + download. No UI exists.

**Related APIs:**
- `GET /api/v1/reports`
- `POST /api/v1/reports`
- `GET /api/v1/reports/{id}`
- `DELETE /api/v1/reports/{id}`
- `POST /api/v1/reports/security`
- `POST /api/v1/reports/traffic`
- `POST /api/v1/reports/incidents`
- `POST /api/v1/reports/compliance`
- `GET /api/v1/reports/{id}/download`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (name, kind, status, created_at, actions)
- **Generate:** Buttons for each report type (security, traffic, incidents, compliance)
- **Detail:** Report detail with summary data viewer
- **Download:** Download button for each completed report
- **Delete:** ConfirmDialog
- Show report status (pending, ready, failed)

**Acceptance Criteria:**
- [ ] Report list shows all generated reports
- [ ] Generate buttons work for each type
- [ ] Detail page shows summary data
- [ ] Download works (returns JSON attachment)
- [ ] Delete works with confirmation
- [ ] Loading, empty, error states

**Testing:** Generate each report type, view detail, download, delete.

---

## UI-032 — Create System Settings Page

**Priority:** P1  
**Module:** Settings  
**Type:** Page  
**Purpose:** Backend has full settings API (general, security, localization, retention). No UI exists.

**Related APIs:**
- `GET /api/v1/settings`
- `PATCH /api/v1/settings`
- `GET /api/v1/settings/security`
- `PATCH /api/v1/settings/security`
- `GET /api/v1/settings/localization`
- `PATCH /api/v1/settings/localization`
- `GET /api/v1/settings/retention`
- `PATCH /api/v1/settings/retention`

**Existing Files:** None

**Dependencies:** UI-006, UI-011

**Implementation Requirements:**
- **Security tab:** Session timeout, password policy, lockout settings
- **Localization tab:** Default language, locale settings
- **Retention tab:** Log retention days, event TTL
- **General tab:** System name, contact info
- Each setting is a key-value form with save button
- Use PATCH with upsert semantics

**Acceptance Criteria:**
- [ ] Settings page loads real data from API
- [ ] Each category tab loads correctly
- [ ] Save button updates settings
- [ ] Loading, error states

**Testing:** Update each setting category, verify persistence.

---

## UI-033 — Create Certificate Management Page

**Priority:** P2  
**Module:** Certificates  
**Type:** Page  
**Purpose:** Backend has full certificate management (upload, import, renew, expiry). No UI exists.

**Related APIs:**
- `GET /api/v1/certificates`
- `POST /api/v1/certificates`
- `GET /api/v1/certificates/{id}`
- `DELETE /api/v1/certificates/{id}`
- `POST /api/v1/certificates/import`
- `POST /api/v1/certificates/{id}/renew`
- `GET /api/v1/certificates/{id}/expiry`
- `GET /api/v1/tls-profiles`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (name, domain, expiry, status)
- **Upload:** File upload for certificate files
- **Import:** Import from PEM/chain
- **Detail:** Certificate detail with expiry date
- **Renew:** Renew button
- **Delete:** ConfirmDialog
- **TLS Profiles:** List, create, edit, delete TLS profiles

**Acceptance Criteria:**
- [ ] Certificate list shows all certificates
- [ ] Upload/Import works
- [ ] Renew works
- [ ] Expiry dates shown correctly
- [ ] TLS profiles CRUD works

**Testing:** Upload certificate, view detail, renew, delete.

---

## UI-034 — Create Backups Page

**Priority:** P2  
**Module:** Backups  
**Type:** Page  
**Purpose:** Backend has backup/restore endpoints. No UI exists.

**Related APIs:**
- `GET /api/v1/backups`
- `POST /api/v1/backups`
- `POST /api/v1/backups/{id}/restore`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (backup name, size, created_at, status)
- **Create:** Trigger backup button
- **Restore:** ConfirmDialog + restore action
- Show backup status (completed, failed, in_progress)

**Acceptance Criteria:**
- [ ] Backup list shows all backups
- [ ] Create backup triggers backup
- [ ] Restore works with confirmation dialog
- [ ] Loading, empty, error states

**Testing:** Create backup, restore, verify.

---

## UI-035 — Create IP Lists and Rate Limiting Pages

**Priority:** P2  
**Module:** IP Lists  
**Type:** Page  
**Purpose:** Backend has IP list and rate limit management. No UI exists.

**Related APIs:**
- `GET /api/v1/ip-lists`
- `POST /api/v1/ip-lists`
- `GET /api/v1/ip-lists/{id}`
- `PATCH /api/v1/ip-lists/{id}`
- `DELETE /api/v1/ip-lists/{id}`
- `GET /api/v1/ip-lists/{id}/entries`
- `POST /api/v1/ip-lists/{id}/entries`
- `DELETE /api/v1/ip-lists/{id}/entries/{entryId}`
- `GET /api/v1/rate-limits`
- `POST /api/v1/rate-limits`
- `PATCH /api/v1/rate-limits/{id}`
- `DELETE /api/v1/rate-limits/{id}`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **IP Lists tab:** List, create, edit, delete IP lists
- **IP Entries:** Add/remove IP/CIDR entries per list
- **Rate Limits tab:** List, create, edit, delete rate limit policies
- Support: IP list type (allow/block), rate limit fields (limit, window, action)

**Acceptance Criteria:**
- [ ] IP lists CRUD works
- [ ] IP entries add/remove works
- [ ] Rate limits CRUD works
- [ ] Loading, empty, error states

**Testing:** Create IP list, add IP entries, create rate limit, verify.

---

## UI-036 — Create Threat Intelligence Page

**Priority:** P2  
**Module:** Threat Intelligence  
**Type:** Page  
**Purpose:** Backend has threat feed management. No UI exists.

**Related APIs:**
- `GET /api/v1/threat-feeds`
- `POST /api/v1/threat-feeds`
- `GET /api/v1/threat-feeds/{id}`
- `PATCH /api/v1/threat-feeds/{id}`
- `DELETE /api/v1/threat-feeds/{id}`
- `POST /api/v1/threat-feeds/{id}/sync`
- `GET /api/v1/threat-feeds/{id}/indicators`
- `POST /api/v1/threat-feeds/{id}/test`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (name, source, type, confidence, status)
- **Create:** Feed form (name, source, type, indicators)
- **Detail:** Feed detail with indicators list
- **Actions:** Sync, Test
- **Delete:** ConfirmDialog
- Show indicator count per feed

**Acceptance Criteria:**
- [ ] Threat feeds CRUD works
- [ ] Sync action triggers sync
- [ ] Test action tests feed
- [ ] Indicators list shows correctly
- [ ] Loading, empty, error states

**Testing:** Create feed, add indicators, sync, test, delete.

---

### Library / Integration

---

## UI-037 — Add Missing API Client Functions

**Priority:** P0  
**Module:** API Client  
**Type:** Integration  
**Purpose:** The API client (`api.ts`) only has 10 functions. Many backend endpoints are missing from the client. Add all necessary API client functions.

**Related APIs:** All backend endpoints (see API inventory)

**Existing Files:**
- `apps/console-web/src/lib/api.ts`

**Dependencies:** None

**Implementation Requirements:**
- Add TypeScript types for all API resources
- Add functions for all CRUD modules:
  - Applications (detail, update, delete, domains, origins, traffic, events, incidents, policies, health)
  - Security Policies (all lifecycle actions)
  - Security Events (detail, matches, timeline, mask, export)
  - Incidents (all actions)
  - Gateways (detail, update, delete, config, status, metrics, apply, rollback)
  - Audit Events (detail, export)
  - Users (update, delete, roles)
  - Groups (CRUD)
  - Roles (list, get)
  - Integrations (CRUD, test, enable, disable)
  - IP Lists (CRUD, entries)
  - Rate Limits (CRUD)
  - Certificates (CRUD, import, renew, expiry)
  - Backups (CRUD, restore)
  - Reports (CRUD, generate, download)
  - Notification Channels (CRUD, test)
  - Settings (get, update per category)
  - License (get, activate, deactivate, usage, entitlements)
  - Dashboard (all 8 endpoints)
  - Learning (sessions, suggestions)
  - Bot Management (policies, events, clients)
  - DLP (profiles)
  - API Security (schemas)
  - Threat Intelligence (feeds)
  - Bot Policies
  - Caching, Clusters, Automation
  - Webhooks, Exceptions
  - Rules, Rule Bundles
  - TLS Profiles
  - Backend Pools, Nodes, Health Monitors
  - Virtual Servers, Listeners, Sites
  - GraphQL Security, Client-Side Protection, API Quotas, ML Baselines, Network Protection
  - Organizations, Tenants
  - Policy Approvals, Policy Versions

**Acceptance Criteria:**
- [ ] Every backend endpoint has a corresponding API client function
- [ ] TypeScript types match backend response shapes
- [ ] Functions handle errors consistently
- [ ] Auth headers are sent correctly

**Testing:** Use each function in a page, verify API calls.

---

## UI-038 — Fix Security Events to Use API Client

**Priority:** P0  
**Module:** API Client  
**Type:** Integration  
**Purpose:** The security events page uses raw `fetch()` instead of the API client. Fix it.

**Related APIs:**
- `GET /api/v1/security-events`

**Existing Files:**
- `apps/console-web/src/app/[locale]/security-events/page.tsx`

**Dependencies:** UI-037

**Implementation Requirements:**
- Replace raw `fetch()` with `listSecurityEvents()` from API client
- Ensure all query parameters are passed through

**Acceptance Criteria:**
- [ ] Security events page uses API client
- [ ] No raw `fetch()` calls in the page
- [ ] All existing functionality preserved

**Testing:** Security events page loads correctly.

---

## UI-039 — Fix Audit Page to Use API Client

**Priority:** P0  
**Module:** API Client  
**Type:** Integration  
**Purpose:** The audit page uses raw `fetch()` instead of the API client. Fix it.

**Related APIs:**
- `GET /api/v1/audit-events`

**Existing Files:**
- `apps/console-web/src/app/[locale]/audit/page.tsx`

**Dependencies:** UI-037

**Implementation Requirements:**
- Replace raw `fetch()` with `listAuditEvents()` from API client
- Ensure all query parameters are passed through

**Acceptance Criteria:**
- [ ] Audit page uses API client
- [ ] No raw `fetch()` calls in the page
- [ ] All existing functionality preserved

**Testing:** Audit page loads correctly.

---

### RBAC UI Tasks

---

## UI-040 — Add Permission-Aware Navigation

**Priority:** P1  
**Module:** RBAC  
**Type:** Component  
**Purpose:** The sidebar shows all navigation links regardless of user permissions. Hide links that the user's role cannot access.

**Related APIs:**
- `GET /api/v1/auth/me` (returns user role/permissions)

**Existing Files:**
- `apps/console-web/src/components/layout/Sidebar.tsx`
- `apps/console-web/src/context/AuthContext.tsx`

**Dependencies:** UI-006

**Implementation Requirements:**
- Create a permission check utility based on the user's role
- Conditionally render nav links based on permissions
- Example: "Auditor" role should only see Audit Log
- Example: "App Owner" should only see Applications, Traffic, Gateways
- Map backend roles to frontend nav visibility

**Acceptance Criteria:**
- [ ] Sidebar shows only accessible links for current user role
- [ ] Unauthorized routes redirect to 403
- [ ] Permission check is consistent with backend

**Testing:** Login as different roles, verify sidebar.

---

## UI-041 — Add Permission-Aware Action Buttons

**Priority:** P1  
**Module:** RBAC  
**Type:** Component  
**Purpose:** Hide or disable action buttons (create, edit, delete) based on user permissions.

**Related APIs:**
- `GET /api/v1/auth/me`

**Existing Files:**
- `apps/console-web/src/context/AuthContext.tsx`
- `apps/console-web/src/components/shared/` (after creation)

**Dependencies:** UI-040

**Implementation Requirements:**
- Create a `usePermission()` hook
- Example: `canCreateApp` / `canDeleteUser` / `canActivatePolicy`
- Map backend permissions to frontend actions
- Hide buttons the user cannot use
- Grey out or hide delete buttons without write permission

**Acceptance Criteria:**
- [ ] Create buttons hidden for read-only roles
- [ ] Edit buttons hidden for read-only roles
- [ ] Delete buttons hidden for non-admin roles
- [ ] Permission check is consistent with backend

**Testing:** Login as different roles, verify action buttons.

---

## UI-042 — Create 403 Forbidden Page

**Priority:** P1  
**Module:** RBAC  
**Type:** Page  
**Purpose:** When a user navigates to a page they don't have permission to access, show a proper 403 page.

**Related APIs:** None

**Existing Files:** None

**Dependencies:** UI-040

**Implementation Requirements:**
- Create `src/app/[locale]/forbidden/page.tsx`
- Show a clear message: "You do not have permission to access this page"
- Include a "Go to Dashboard" button
- Use the dashboard layout

**Acceptance Criteria:**
- [ ] 403 page renders with clear message
- [ ] "Go to Dashboard" button works
- [ ] 403 page is styled consistently

**Testing:** Access a restricted route as a read-only user.

---

### i18n Tasks

---

## UI-043 — Translate Hardcoded English Strings

**Priority:** P1  
**Module:** i18n  
**Type:** Integration  
**Purpose:** Many pages have hardcoded English strings that are not translated. Add them to the message catalogs.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/src/messages/en.json`
- `apps/console-web/src/messages/bn.json`
- `apps/console-web/src/app/[locale]/page.tsx`
- `apps/console-web/src/app/[locale]/users/page.tsx`
- `apps/console-web/src/app/[locale]/audit/page.tsx`
- `apps/console-web/src/components/UserProfileWidget.tsx`
- `apps/console-web/src/components/CreateUserButton.tsx`
- `apps/console-web/src/app/[locale]/login/page.tsx`

**Dependencies:** None

**Implementation Requirements:**
- Add translation keys for:
  - "Logout" (UserProfileWidget)
  - "Authenticating..." (login page)
  - "Add User", "Create New User", "Email address", "Password", "Role", "Cancel", "Create User", "Creating..." (CreateUserButton)
  - "No users found", "Showing cached data — API offline", "Active", "Pending", "Member Since" (users page)
  - "API offline" (audit page)
  - "100% Protected", "1 Node Offline", "All Synced", "View All", "Protected Applications", "Deploy Config" (overview page)
  - "system" (audit page fallback)
- Add Bangla translations for all new keys

**Acceptance Criteria:**
- [ ] All hardcoded English strings use translation keys
- [ ] Bangla translations exist for all new keys
- [ ] No missing translation errors

**Testing:** Switch to each language, verify all strings are translated.

---

### UX/Accessibility Tasks

---

## UI-044 — Add Loading States to All Pages

**Priority:** P1  
**Module:** UX  
**Type:** Integration  
**Purpose:** No pages have loading states. Users see a blank page until data is fetched. Add loading skeletons to all pages.

**Related APIs:** All API endpoints

**Existing Files:** All pages

**Dependencies:** UI-005

**Implementation Requirements:**
- Add `loading.tsx` in the `[locale]` directory
- Add loading skeletons to each page (DataTable skeleton, card skeleton)
- For pages that fetch multiple resources, show loading state per section

**Acceptance Criteria:**
- [ ] All pages show a loading state while fetching data
- [ ] Loading state is visual (skeleton, spinner, or shimmer)
- [ ] Loading state disappears after data loads

**Testing:** Navigate between pages, verify loading states.

---

## UI-045 — Add Empty States to All Pages

**Priority:** P1  
**Module:** UX  
**Type:** Integration  
**Purpose:** Some pages have empty states, but they are inconsistent. Add proper empty states to all list pages.

**Related APIs:** All list endpoints

**Existing Files:** All pages

**Dependencies:** UI-008

**Implementation Requirements:**
- For every list page, show an empty state when no data exists
- Empty state should include: icon, message, action button (e.g., "Create your first application")
- Use the DataTable's empty state

**Acceptance Criteria:**
- [ ] All list pages show empty state when no data
- [ ] Empty state includes a call to action
- [ ] Empty state disappears after data is created

**Testing:** View each page with no data, verify empty state.

---

## UI-046 — Add Error States to All Pages

**Priority:** P1  
**Module:** UX  
**Type:** Integration  
**Purpose:** Most pages silently catch API errors and return empty arrays. Add proper error states with retry.

**Related APIs:** All API endpoints

**Existing Files:** All pages

**Dependencies:** UI-005

**Implementation Requirements:**
- For every page, show an error state when API fails
- Error state should include: icon, error message, retry button
- Use the DataTable's error state
- Only `/users` and `/audit` currently show any error — fix all others

**Acceptance Criteria:**
- [ ] All pages show error state when API fails
- [ ] Error state includes a retry button
- [ ] Error state disappears after retry succeeds

**Testing:** Disconnect API, load each page, verify error state.

---

### Testing Tasks

---

## UI-047 — Add Vitest Configuration and First Tests

**Priority:** P2  
**Module:** Testing  
**Type:** Testing  
**Purpose:** No frontend tests exist. Set up Vitest with configuration and add tests for the most critical components.

**Related APIs:** None

**Existing Files:**
- `apps/console-web/package.json`
- `apps/console-web/tsconfig.json`

**Dependencies:** None

**Implementation Requirements:**
- Create `vitest.config.ts` (or `.mts`)
- Add React Testing Library + jsdom
- Add tests for:
  - Sidebar component (renders nav links, highlights active)
  - UserProfileWidget (shows user info, logout works)
  - CreateUserButton (form validation, API call)
  - DataTable component (renders rows, sorting, pagination)
  - StatusBadge, SeverityBadge (renders correct colors)
  - Utility functions
- Ensure `npm test` runs all tests

**Acceptance Criteria:**
- [ ] Vitest configuration is correct
- [ ] Component tests pass
- [ ] `npm test` runs successfully
- [ ] Tests cover critical paths

**Testing:** `npm test` passes.

---

## UI-048 — Add Page-Level Integration Tests

**Priority:** P3  
**Module:** Testing  
**Type:** Testing  
**Purpose:** Add integration tests for critical pages using Mock Service Worker (MSW) to mock API responses.

**Related APIs:** All API endpoints

**Existing Files:** All pages

**Dependencies:** UI-047

**Implementation Requirements:**
- Add MSW (Mock Service Worker) for API mocking
- Add tests for:
  - Login page (successful login, failed login, validation)
  - Overview page (loads data, shows error)
  - Applications page (loads list, empty state, error)
- Test loading, empty, error states

**Acceptance Criteria:**
- [ ] MSW handlers cover all API endpoints used by tested pages
- [ ] Integration tests pass
- [ ] Tests cover loading, empty, error, and success states

**Testing:** `npm test` passes.

---

### API → UI Gaps

---

## UI-049 — Create Learning Sessions Page

**Priority:** P2  
**Module:** Learning  
**Type:** Page  
**Purpose:** Backend has full learning session management. No UI exists.

**Related APIs:**
- `GET /api/v1/learning/sessions`
- `POST /api/v1/learning/sessions`
- `GET /api/v1/learning/sessions/{id}`
- `POST /api/v1/learning/sessions/{id}/start`
- `POST /api/v1/learning/sessions/{id}/stop`
- `GET /api/v1/learning/suggestions`
- `GET /api/v1/learning/suggestions/{id}`
- `POST /api/v1/learning/suggestions/{id}/accept`
- `POST /api/v1/learning/suggestions/{id}/reject`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **Sessions list:** DataTable (name, source, status, confidence_threshold, created_at)
- **Create session:** Form (name, source must be "trusted", description, confidence_threshold)
- **Session detail:** Session info with controls (start/stop)
- **Suggestions list:** DataTable (rule_id, severity, confidence, status)
- **Suggestion detail:** Suggestion info with accept/reject actions
- **Accept suggestion:** ConfirmDialog → POST /accept
- **Reject suggestion:** ConfirmDialog → POST /reject

**Acceptance Criteria:**
- [ ] Learning sessions CRUD works
- [ ] Start/Stop session works
- [ ] Suggestions list shows correctly
- [ ] Accept/Reject actions work
- [ ] Loading, empty, error states

**Testing:** Create session, start/stop, generate suggestions, accept/reject.

---

## UI-050 — Create Notification Channels Page

**Priority:** P2  
**Module:** Notification Channels  
**Type:** Page  
**Purpose:** Backend has full notification channel management. No UI exists.

**Related APIs:**
- `GET /api/v1/notification-channels`
- `POST /api/v1/notification-channels`
- `GET /api/v1/notification-channels/{id}`
- `PATCH /api/v1/notification-channels/{id}`
- `DELETE /api/v1/notification-channels/{id}`
- `POST /api/v1/notification-channels/{id}/test`

**Existing Files:** None

**Dependencies:** UI-006, UI-008, UI-009, UI-010, UI-011

**Implementation Requirements:**
- **List:** DataTable (name, kind, is_default, enabled, created_at)
- **Create:** Channel form (name, kind dropdown, config fields based on kind)
- **Edit:** Edit channel
- **Delete:** ConfirmDialog
- **Test:** Test button per channel
- Supported kinds: wazuh, syslog, cef, leef, webhook, email, teams, slack, soar

**Acceptance Criteria:**
- [ ] Notification channels CRUD works
- [ ] Test action works
- [ ] Kind-specific config fields shown
- [ ] Loading, empty, error states

**Testing:** Create each channel kind, test, edit, delete.

---

## UI-051 — Create License Management Page

**Priority:** P2  
**Module:** License  
**Type:** Page  
**Purpose:** Backend has license management. No UI exists.

**Related APIs:**
- `GET /api/v1/license`
- `POST /api/v1/license/activate`
- `POST /api/v1/license/deactivate`
- `GET /api/v1/license/usage`
- `GET /api/v1/license/entitlements`

**Existing Files:** None

**Dependencies:** UI-006, UI-011

**Implementation Requirements:**
- **License info:** Current license card (edition, status, seats, expiry)
- **Usage:** Usage bars (gateways, applications) with current/limit/percentage
- **Entitlements:** Feature list with checkmarks
- **Activate:** License key input + activate button
- **Deactivate:** ConfirmDialog + deactivate button

**Acceptance Criteria:**
- [ ] License info shows correctly
- [ ] Usage bars show current usage
- [ ] Entitlements show feature flags
- [ ] Activate/Deactivate works
- [ ] Loading, empty, error states

**Testing:** View license, activate, deactivate, verify usage.

---

## 6. UI → API Mismatches

| File | Issue | Fix |
|---|---|---|
| `audit/page.tsx` | Access `ev.resource_type` but backend returns `"resource"` | Change to `ev.resource` |
| `security-events/page.tsx` | Uses raw `fetch()` instead of API client | Replace with `listSecurityEvents()` |
| `audit/page.tsx` | Uses raw `fetch()` instead of API client | Replace with `listAuditEvents()` |
| `login/page.tsx` | `router.push("/en")` ignores locale | Use current locale from params |
| `api.ts` | Port 8080 (`.env.local`) should be 8443 | Fix `.env.local` |
| `Sidebar.tsx` | Uses `t("title")` from `useTranslations("overview")` for Overview link | Use `navT("overview")` for consistency |
| `api.ts` | `X-User-Role` header is sent for mock auth | Remove when real auth is in place |

---

## 7. Summary

### API Statistics

| Metric | Count |
|---|---|
| Total API endpoints | 349 |
| Unique resource modules | 67 |
| Database tables | 54 |

### Frontend Statistics

| Metric | Count |
|---|---|
| Total frontend routes | 7 (plus 1 root redirect) |
| Total existing pages | 6 |
| Total complete pages | 1 (login) |
| Total partial pages | 6 (overview, apps, gateways, events, users, audit) |
| Total missing pages | 6 (policies, integrations, reports, settings, incidents, traffic) |
| Total broken sidebar links | 6 |

### UI Task Summary

| Priority | Count |
|---|---|
| P0 | 20 |
| P1 | 18 |
| P2 | 11 |
| P3 | 2 |
| **Total UI tasks** | **51** |

### Breakdown

| Category | Count |
|---|---|
| Foundation bugs | 5 |
| Shared components | 5 |
| Dashboard widgets | 7 |
| Core resource pages | 7 |
| Operations pages | 5 |
| Missing pages | 5 |
| API client integration | 3 |
| RBAC UI | 3 |
| i18n | 1 |
| UX/Accessibility | 3 |
| Testing | 2 |
| API → UI gaps | 3 |
| **Total** | **51** |

---

## 8. Recommended Execution Order

```
Phase 1 — Foundation Bugs (P0)
  UI-001  Fix CSS Import Bug
  UI-002  Fix API Port Mismatch
  UI-003  Fix Audit Event resource_type Field
  UI-004  Fix Login Redirect Locale
  UI-005  Add loading.tsx and error.tsx Pages

Phase 2 — Shared Components (P0)
  UI-037  Add Missing API Client Functions
  UI-008  Create DataTable Component
  UI-009  Create StatusBadge and SeverityBadge Components
  UI-010  Create ConfirmDialog Component
  UI-011  Create PageHeader Component
  UI-012  Create FilterBar Component

Phase 3 — Application Shell (P0)
  UI-006  Create Dashboard Layout Component
  UI-007  Fix Sidebar Navigation Links

Phase 4 — Core Pages (P0)
  UI-021  Complete Applications Page
  UI-024  Complete Gateways Page
  UI-025  Complete Security Events Page
  UI-038  Fix Security Events to Use API Client
  UI-039  Fix Audit Page to Use API Client
  UI-028  Create Audit Log Page with Pagination and Export
  UI-029  Create Users Page with Edit and Delete

Phase 5 — Dashboard (P1)
  UI-014  Replace Dashboard Mock Data with Real API Calls
  UI-015  Create Dashboard Traffic Widget
  UI-016  Create Dashboard Security Widget
  UI-017  Create Dashboard Incidents Widget
  UI-018  Create Dashboard Top IPs Widget
  UI-019  Create Dashboard Top Rules Widget
  UI-020  Create Dashboard Gateway Health Widget

Phase 6 — Policy & Incident Management (P1)
  UI-026  Create Policy Pages
  UI-027  Create Incidents Page

Phase 7 — Missing Pages (P1)
  UI-030  Create Integrations Page
  UI-031  Create Reports Page
  UI-032  Create System Settings Page
  UI-033  Create Certificate Management Page
  UI-034  Create Backups Page
  UI-035  Create IP Lists and Rate Limiting Page
  UI-036  Create Threat Intelligence Page

Phase 8 — Learning & Notifications (P2)
  UI-049  Create Learning Sessions Page
  UI-050  Create Notification Channels Page
  UI-051  Create License Management Page

Phase 9 — RBAC & i18n (P1)
  UI-040  Add Permission-Aware Navigation
  UI-041  Add Permission-Aware Action Buttons
  UI-042  Create 403 Forbidden Page
  UI-043  Translate Hardcoded English Strings

Phase 10 — UX Polish (P1)
  UI-044  Add Loading States to All Pages
  UI-045  Add Empty States to All Pages
  UI-046  Add Error States to All Pages
  UI-013  Add MFA Enable/Verify/Disable Pages
  UI-022  Create Application Form
  UI-023  Create Application Detail Page

Phase 11 — Testing (P2)
  UI-047  Add Vitest Configuration and First Tests
  UI-048  Add Page-Level Integration Tests
```

Total: **51 tasks** across 11 phases, ordered by dependency and priority.