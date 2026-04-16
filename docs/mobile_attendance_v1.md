# Mobile Attendance V1

This document defines the minimum backend readiness for the first mobile app focused on employee attendance.

## Goals

- Support employee login with the existing auth flow.
- Support employee self-service attendance check-in and check-out.
- Support employee access to their own attendance history.
- Support owner/admin visibility for attendance audit on the web.
- Keep the first API surface small and stable.

## Non-Goals For V1

- Offline sync
- Face recognition / selfie verification
- Advanced approval workflow
- Complex shift roster engine
- Push notification orchestration
- Multi-device trust policy

## V1 Domain Set

The minimum recommended domains before the mobile app starts:

1. Auth
2. Me / employee context
3. Attendance logs
4. Attendance summary
5. Attendance policy
6. Shift snapshot for today

## Route Group

Use the existing route style and keep attendance routes under a single group:

- `GET /attendance-logs`
- `GET /attendance-logs/:id`
- `POST /attendance-logs/check-in`
- `POST /attendance-logs/check-out`

Additional mobile-readiness routes recommended for V1:

- `GET /me`
- `GET /me/employee`
- `GET /attendance-summary/today`
- `GET /attendance-policy`
- `GET /me/shift-today`

## Access Rules

### Owner / Admin

- Can view attendance logs across the tenant.
- Can filter by employee and date range.
- Must not use the employee check-in/check-out endpoints to record attendance for another user in V1.

### Employee

- Can only view their own attendance logs.
- Can only create attendance for themselves.
- `employee_id` must be resolved by the server from the logged-in user.

## Existing V1-Ready Endpoints

These are already aligned with the current attendance implementation direction:

### `GET /attendance-logs`

Purpose:
- List attendance logs

Behavior:
- Owner sees tenant-wide logs
- Employee only sees their own logs

Query:
- `page`
- `limit`
- `q`
- `employee_id`
- `type`
- `source`
- `attendance_date`
- `date_from`
- `date_to`

### `GET /attendance-logs/:id`

Purpose:
- Fetch a single attendance log

Behavior:
- Owner can access logs within tenant
- Employee can access only their own log

### `POST /attendance-logs/check-in`

Purpose:
- Record employee mobile check-in

### `POST /attendance-logs/check-out`

Purpose:
- Record employee mobile check-out

## Missing Endpoints Recommended Before Mobile Starts

### `GET /me`

Purpose:
- Return auth identity plus top-level role/tenant context needed for bootstrapping the app

Suggested response:

```json
{
  "data": {
    "user_id": "uuid",
    "user_name": "EMP-ATT-001",
    "tenant_id": "uuid",
    "tenant_name": "PT Beruang",
    "roles": ["employee"]
  },
  "message": "Your request has been successfully processed",
  "success": true
}
```

### `GET /me/employee`

Purpose:
- Return employee profile for the logged-in user

Suggested fields:
- `id`
- `employee_no`
- `full_name`
- `email`
- `org_unit_id`
- `job_position_id`
- `location_id`
- `shift_id`
- `status`

### `GET /attendance-summary/today`

Purpose:
- Drive the main mobile home screen

Suggested response:

```json
{
  "data": {
    "attendance_date": "2026-03-28",
    "today_status": "checked_in",
    "checked_in": true,
    "checked_out": false,
    "check_in_log_id": "uuid",
    "check_out_log_id": null,
    "last_log_type": "check_in",
    "last_logged_at": "2026-03-28T08:05:00+08:00",
    "can_check_in": false,
    "can_check_out": true
  },
  "message": "Your request has been successfully processed",
  "success": true
}
```

### `GET /attendance-policy`

Purpose:
- Expose the current policy the mobile client needs to present and validate against

Suggested fields:
- `timezone`
- `attendance_radius_meters`
- `attendance_check_in_start`
- `attendance_check_in_end`
- `attendance_check_out_start`
- `attendance_check_out_end`

### `GET /me/shift-today`

Purpose:
- Give the employee a simple shift snapshot for the current day

Suggested fields:
- `shift_id`
- `shift_name`
- `start_time`
- `end_time`
- `attendance_date`

## Request / Response Contracts

### Check-In Request

```json
{
  "logged_at": "2026-03-28T08:05:00+08:00",
  "latitude": -5.1477,
  "longitude": 119.4327,
  "address": "Makassar",
  "device_id": "ios-sim-001",
  "device_name": "iPhone Test",
  "notes": "Arrived"
}
```

### Check-Out Request

```json
{
  "logged_at": "2026-03-28T17:15:00+08:00",
  "device_id": "ios-sim-001",
  "device_name": "iPhone Test"
}
```

### Standard Success Response

```json
{
  "data": {
    "id": "uuid"
  },
  "message": "Attendance recorded successfully",
  "success": true
}
```

### Standard Error Cases

- `400` duplicate check-in
- `400` check-out before check-in
- `400` invalid timestamp payload
- `401` missing or invalid token
- `404` employee profile not found

## Business Rules To Freeze Before Mobile Build

These must be explicit before UI implementation:

1. What timezone is authoritative for attendance day calculation?
2. Can employee check-in outside the allowed time window?
3. Can employee check-in outside the allowed radius?
4. Will outside-policy attendance be rejected or stored with a different status?
5. Is duplicate prevention per day or per shift?
6. Does mobile need idempotency protection for retry?

## Recommended Rule Decisions For V1

To keep implementation simple:

1. Use company timezone as the source of truth.
2. Allow only one `check_in` and one `check_out` per employee per attendance date.
3. Reject `check_out` if no `check_in` exists for that date.
4. Do not implement radius/time-window enforcement yet unless business confirms exact behavior.
5. Store device metadata now, even if it is not yet enforced.

## Automated Test Priorities

### Backend

1. Employee can check in successfully.
2. Duplicate check-in is rejected.
3. Check-out without check-in is rejected.
4. Check-out after check-in succeeds.
5. Employee list endpoint only returns own logs.
6. Owner list endpoint returns tenant logs.
7. Employee cannot open another employee's log detail.

### Web

1. Owner attendance page loads and shows logs.
2. Employee attendance page shows only employee logs.
3. Attendance list keeps rendering correctly without hydration mismatch.

## Implementation Backlog

### Phase 1: Ready Before Mobile UI

- [x] Attendance log table and CRUD-style read endpoints
- [x] Employee self check-in/check-out endpoints
- [x] `GET /me`
- [x] `GET /me/employee`
- [x] `GET /attendance-summary/today`
- [x] `GET /attendance-policy`
- [x] `GET /me/shift-today`

### Phase 2: Strongly Recommended After Bootstrapping

- [ ] Idempotency key support for attendance actions
- [ ] Attendance request / correction flow
- [ ] Leave request module
- [ ] Holiday / calendar feed

### Phase 3: Operational Hardening

- [ ] Radius enforcement
- [ ] Time-window enforcement
- [ ] Device registration / session management
- [ ] Evidence upload for attendance correction

## Definition Of Ready For Mobile App Start

The mobile app should not start before:

1. The Phase 1 endpoints exist.
2. Auth flow is stable.
3. Error messages are predictable for expected attendance failures.
4. At least one employee test account exists.
5. Backend test cases for attendance flow exist.
6. Owner web attendance page already works for support/audit visibility.
