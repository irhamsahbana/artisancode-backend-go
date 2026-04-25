# User Invitation Flow

## Goal

Provide a clear invitation flow so the product does not depend on a generic `Users` menu while still keeping a clean access lifecycle for:

- owner
- admin
- employee

This document is an early product and backend planning artifact. Frontend and backend implementation can follow incrementally.

## Product Direction

Main principles:

- `users` is the identity/auth layer, not the main product menu
- `employees` is the domain object for workforce data
- access should be granted through business-actor-specific flows, not through generic user CRUD

UX implications:

- owner enters through owner login/register
- admin is invited or created by owner
- employee is created through the employee module, then login access is enabled only when needed

## Target Flows

### 1. Owner

- The owner registers tenant and company during initial registration.
- After login, the owner can manage admin access.
- The owner is not created through the normal invitation flow unless a future multi-owner requirement appears.

### 2. Admin

- The owner opens the `Access` or `Admin` page.
- The owner sends an invitation to the admin email.
- The admin receives the link, creates a password, and the invitation status becomes `accepted`.

### 3. Employee

- Owner or admin creates the employee record first.
- The employee can remain active as HR data even without a login account.
- When access is needed, owner or admin sends an invitation from the employee detail or list.
- The employee receives the activation link, creates a password, and is connected to a login account.

## Proposed Domain Model

Recommended minimum model:

- `users`
  - login identity
  - email, password hash, verification status, tenant scope
- `employees`
  - employee domain profile
  - optional `user_id` until the employee activates the account
- `user_roles` or equivalent membership
  - role relations such as owner, admin, employee
- `user_invitations`
  - source of truth for invitations

## Proposed `user_invitations` Fields

Recommended initial columns:

- `id`
- `tenant_id`
- `employee_id` nullable
- `email`
- `role_code`
- `status`
- `token_hash`
- `expires_at`
- `accepted_at` nullable
- `revoked_at` nullable
- `invited_by`
- `last_sent_at`
- `created_at`
- `updated_at`

Minimum `status` values:

- `pending`
- `accepted`
- `expired`
- `revoked`

Notes:

- Store the token hash, not the raw token.
- Invitation email should be immutable per record so the audit trail stays clear.
- For employee invitations, `employee_id` binds the invitation to the correct domain record.

## Backend API Plan

The first phase can focus on these APIs:

### Protected APIs

- `POST /user-invitations`
  - create a new invitation for admin or employee
- `GET /user-invitations`
  - list or filter invitations
- `POST /user-invitations/:id/resend`
  - resend a still-valid invitation or regenerate the token
- `POST /user-invitations/:id/revoke`
  - revoke an invitation

### Public APIs

- `GET /user-invitations/accept`
  - validate the token and show invitation summary
- `POST /user-invitations/accept`
  - set password, verify email, bind role, and bind employee when present

## Acceptance Flow

When an invitation is accepted:

1. validate token, status, tenant, and expiry
2. verify the invitation email still matches
3. create a new user if none exists
4. if a user already exists in the same tenant, reuse that user under strict rules
5. assign the invited role
6. link `employees.user_id` for employee invitations
7. set `email_verified_at`
8. mark the invitation as accepted

## Guardrails

Recommended baseline rules:

- There must not be two active invitations for the same email + role + tenant.
- An employee already linked to an active user must not be re-invited without an explicit reset or reinvite flow.
- Acceptance should be idempotent for the same token after success.
- Admin must not invite a role above their own level.
- High-level admin invitation should be owner-only.

## FE Plan

Recommended frontend phases:

### Phase 1

- hide the `Users` menu from sidebar and dashboard
- keep the internal route if the team still uses it
- add invitation entry points in:
  - employee pages
  - roles/access pages for admin management

### Phase 2

- create an `Invitations` page or an invitation panel inside `Roles`
- show statuses:
  - pending
  - accepted
  - expired
  - revoked
- provide actions:
  - invite
  - resend
  - revoke

### Phase 3

- create a public accept page
- add the set-password form
- display tenant identity, role, and employee summary

## Delivery Plan

Safest implementation order:

1. Add the `user_invitations` table plus base repository and core.
2. Add protected create, list, resend, and revoke APIs.
3. Add the invitation email template.
4. Add the public accept API.
5. Add frontend invitation entry points from employee and access management.
6. Add the frontend accept invitation page.
7. Remove business dependence on the `Users` menu.

## Open Decisions

Product questions that still need agreement:

- whether admin can create invitations for other admins
- whether employee invitations automatically enable app access or require extra approval
- whether one email may be linked to more than one employee across different tenants
- whether future additional owners should use the same invitation flow or a dedicated one

## Recommended MVP

Most effective MVP:

- invitation support only for `admin` and `employee`
- create, resend, revoke, and accept flows
- employee record must already exist before inviting an employee
- owner registration continues through the existing flow

With this approach, the `Users` menu can stay hidden without removing the access control the operations team needs.
