# Alibaba Cloud Call Center (云联络中心 / CCC) — Research and Implementation Proposal

> **Status:** Draft for review. No code will be written until you sign off on this document.
>
> **Stack proposed:** Go (backend), React (frontend / agent workbench), MySQL (persistence), PJSIP (SIP signaling + media).
>
> **Sources:** Alibaba Cloud official documentation under `help.aliyun.com/zh/ccs/*` — in particular:
> - 坐席 (Agents): https://help.aliyun.com/zh/ccs/agents
> - 技能组 (Skill Groups): https://help.aliyun.com/zh/ccs/skill-groups
> - 坐席工作台 (Workbench): https://help.aliyun.com/zh/ccs/workbench
> - SIP 对接指引 (SIP Connection Guide): https://help.aliyun.com/zh/ccs/sip-connection-guide
> - 号码管理 (Phone Number Management / Caller ID): https://help.aliyun.com/zh/ccs/phone-number-management
> - 创建复杂云联络中心 (Reference architecture): https://help.aliyun.com/zh/ccs/getting-started/create-a-complex-call-center-instance
> - 新客户须知 (Usage notes, SIP/400 number rules): https://help.aliyun.com/zh/ccs/usage-notes

---

## Part 1 — Research: Core Functionalities of Aliyun CCC

Aliyun CCC ("云联络中心", formerly "云呼叫中心") is a SaaS contact-center platform built on Alibaba Cloud. The product surface is organized around two operational flows — **inbound** and **outbound** — and six core building blocks the user asked about. Each one is summarized below with the behavior I will replicate.

### 1.1 Agents (坐席)

**Purpose.** An *agent* is a user account that can take and place calls inside the contact center. Agents are the terminal endpoint of every interaction — IVRs, queues, and skill groups all eventually deliver a media session to an agent (or to a configured external number).

**Key attributes observed in the product:**

| Field | Notes |
| --- | --- |
| `user_name` | 4–32 chars, immutable after creation; lowercase letters, digits, `.`, `_`, `-`. |
| `display_name` (姓名) | Free text. |
| `employee_id` (工号) | Optional, for call-center operations that track employees by badge ID. |
| `email` | Used for initial-password delivery and login link. |
| `phone` / `landline` | Used in **off-site mode** — inbound calls can be bridged to the agent's mobile when they are not connected to the softphone. |
| `role` | One of `agent` (坐席), `skill_group_leader` (技能组组长), `admin` (管理员). Each role sees a different UI surface and has different API permissions. |
| `work_mode` | `on_site` (场内 — softphone in browser) or `off_site` (场外 — bridge to mobile/landline). |
| `skill_group_memberships` | An agent can belong to multiple skill groups; per-group **level 1–10** (1 = highest priority). |
| `dedicated_outbound_numbers` | An agent can be granted one or more dedicated CLIs to use for outbound calls. |

**Lifecycle operations:** create, bulk import from RAM (Alibaba's IAM), file import (CSV template, max 500/batch), bulk delete, export, role change.

**Runtime state machine (driven by the workbench, see §1.3):**
`Offline → Online (Idle) → Ringing → Talking → After-Call-Work (话后处理) → Idle | Break (小休) | Invisible | Offline`.

### 1.2 Skill Groups (技能组)

**Purpose.** A *skill group* is a queue of agents organized by capability (e.g. "Spanish-speaking billing tier-2"). Inbound calls are delivered to a skill group from an IVR's "transfer to human" node, and the skill group's routing policy chooses which agent gets the call.

**Key attributes:**

| Field | Notes |
| --- | --- |
| `skill_group_id` | 4–64 chars, starts with a letter; immutable. |
| `name`, `description` | Free text. |
| `members` | Each member has a per-group `level` 1–10. |
| `bound_numbers` | A skill group can own one or more inbound numbers — so agents in the group can both receive and originate from those numbers. |

**Routing semantics observed:**
- IVR "transfer to human" node references the skill group and picks an agent based on level (with secondary policies such as longest-idle, least-busy, round-robin — Aliyun exposes these in the IVR designer).
- If no agent is available, the call queues; queue music, position announcements, and timeout-to-voicemail are configured on the skill group.
- Skill groups can also be bound to outbound numbers, in which case the dialed CLI is restricted to numbers the group owns.

### 1.3 Workbench (坐席工作台 / 工作台)

**Purpose.** A browser-based softphone + agent desktop. This is the operational UI an agent lives in.

**Functional surface (verbatim from the docs):**

- **Presence:** Online (上线 / 在线), Offline (离线), Break (小休, with system + custom break reasons), Invisible (隐身).
- **Inbound:** Answer, hang-up (configurable hang-up policy per instance), hold (通话保持) with hold music, retrieve (取回通话), mute, DTMF dialpad.
- **Outbound:** Manual dial; selection of CLI either **automatically** (best-match by callee region — same-city > same-province > random fallback) or **manually** from the agent's allowed list. Redial after disconnect. Outbound is blocked when no eligible CLI is available.
- **In-call controls:** Hold, retrieve, mute, transfer (warm + cold) — including transfer to another agent, skill group, IVR node, or external number; consultative / 3-way call; supervisor monitor / whisper / barge-in.
- **After-Call Work (ACW, 话后处理):** Configurable timer; agent can extend ACW or end it early; during ACW the agent is not eligible for new calls.
- **Ring-no-answer auto-break:** If an agent does not pick up within X seconds the workbench auto-flips them to Break and re-routes the call.
- **Post-call satisfaction (CSAT):** Optionally trigger an IVR survey by pressing a button before hang-up.
- **Browser requirement:** Chrome ≥ 58 (i.e. WebRTC-capable). The workbench in Aliyun is a WebRTC client to an RTC gateway, not a raw SIP client; we will follow the same pattern (see §2.2).

### 1.4 SIP Trunks (SIP 中继)

**Purpose.** A *SIP trunk* is the signaling+media pipe between the contact center and the PSTN. Aliyun lets enterprise customers bring their **own carrier-supplied SIP gateway** so that their existing numbers and rate plans can be used inside CCC.

**Requirements documented by Aliyun for the customer's side of the trunk (we will impose the same on our system):**

1. Full SIP support per RFC 3261, including **loose routing** (`;lr`).
2. **Transparent transit** of the real calling and called numbers (no NAT/CLI rewrites that strip identity).
3. A **fixed, routable public IP** on the customer's gateway.
4. A qualified SIP operator on the customer's side for on-call operations.
5. HA on the customer's side (active/standby gateways).
6. Customer-side **OPTIONS keepalive** so that when our SIP node stops responding the gateway removes it from the inbound load-balancer set, preserving call continuity.

Aliyun's own side publishes a **primary + standby** SIP service endpoint on UDP 5060 with an "SP marker" string identifying the customer.

**Commercial notes from the docs:** SIP integration is a paid add-on (¥30,000 / line, one-time at the time the docs were written) and a single SIP "line" (语音网关) can carry an unlimited number of DIDs for one carrier — i.e. the *trunk* is the cost unit, not the *number*.

### 1.5 Caller IDs / Phone Numbers (主叫号码 / 号码管理)

**Purpose.** A *caller ID* (also called "号码" in the console) is a DID owned by the tenant. Numbers can be bound to inbound flows, outbound usage, agents, and skill groups.

**Per-number configuration observed:**

| Field | Notes |
| --- | --- |
| `number` | E.164 or local format (the docs treat 400-numbers, fixed-line, and mobile DIDs uniformly). |
| `usage` | One or more of `inbound`, `outbound`, `both`. **400-numbers are inbound-only** per Aliyun rules. |
| `inbound_ivr_id` | The IVR flow that owns the call from `INVITE` to "transfer to human" / hang-up. Required when usage includes inbound and the receiver is a human agent or digital agent. |
| `digital_employee_id` | Alternative to `inbound_ivr_id` for AI-only flows. |
| `skill_groups[]` | Skill groups that may originate **outbound** calls using this CLI. |
| `dedicated_agents[]` | Agents granted exclusive use of this CLI for outbound. |
| `group_label` | Optional tag for management. |

**Auto-CLI-selection algorithm for outbound (Aliyun docs § "拨打电话")**:

1. Find numbers eligible to the agent (dedicated → skill-group-bound → tenant-default).
2. Prefer a CLI in the **same city** as the dialed number.
3. Fall back to **same province**.
4. Fall back to **random** eligible CLI.

This is intended to maximize answer rate, since callees are more likely to pick up local numbers.

### 1.6 Trunk Routing (中继路由)

**Purpose.** Once an outbound call exits the agent and the CLI has been chosen, the system has to decide **which SIP trunk** carries the call to the PSTN. *Trunk routing* is the dial-plan layer that maps `(CLI, dialed-number, time-of-day, …)` to an ordered list of trunks with failover.

**Selection criteria documented or inferable from Aliyun's API model:**

- **Number prefix / region match** — dialed `+8610xxxxxxxx` (Beijing fixed) vs. `+8613xxxxxxxxx` (mobile) go to different carriers.
- **CLI carrier match** — a China Mobile DID must egress through a CMobile-anchored trunk to preserve CLI presentation (otherwise the carrier may rewrite it).
- **Tenant / cost-class** — premium routes vs. grey routes.
- **Trunk health** — drop trunks whose OPTIONS keepalive is failing.
- **Failover order** — primary → standby → reject.

This is precisely the layer **we** must build in Go because, unlike Aliyun, we are the carrier-interconnect side as well as the application side.

---

## Part 2 — Technical Implementation Proposal

### 2.1 High-Level Architecture

```
                       ┌──────────────────────────────────────────────┐
                       │              React Admin + Workbench          │
                       │  (Tenant admin UI; agent softphone via WebRTC)│
                       └──────────────┬───────────────────────────────┘
                                      │ HTTPS (REST + WebSocket)
                                      │ WebRTC (SRTP+DTLS, ICE)
                                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          Go API + Control Plane                      │
│  ┌──────────────┐ ┌──────────────┐ ┌───────────────┐ ┌────────────┐  │
│  │ Auth / RBAC  │ │ Agent / SkG  │ │ IVR / Routing │ │  Reporting │  │
│  │  service     │ │  service     │ │   service     │ │  /CDR svc  │  │
│  └──────────────┘ └──────────────┘ └───────────────┘ └────────────┘  │
│                            │  gRPC  │                               │
│                            ▼        ▼                               │
│             ┌──────────────────────────────────────┐                │
│             │      Call Control Service (Go)       │                │
│             │  - presence, ACD, routing decisions  │                │
│             │  - drives SIP nodes via gRPC/AMQP    │                │
│             └────────────┬─────────────────────────┘                │
└──────────────────────────┼──────────────────────────────────────────┘
                           │ gRPC / NATS (low-latency events)
              ┌────────────┴────────────┐
              ▼                         ▼
   ┌────────────────────┐    ┌──────────────────────┐
   │   SIP Edge Node    │    │   SIP Edge Node      │   (N>=2, anycast)
   │   (PJSIP in C++,   │    │   ...                │
   │   Go shim over CGo │    │                      │
   │   or gRPC sidecar) │    │                      │
   └────────┬───────────┘    └───────────┬──────────┘
            │ SIP/RTP                    │
            ▼                            ▼
        Carrier SIP Trunks (China Mobile / Unicom / Telecom / SIP DID providers)
                                     │
                                     ▼
                                    PSTN
                                     │
                                     ▼
                                  Callees

   Workbench media path:
   Browser ──WebRTC──▶ Media Gateway (RTPengine / Jitsi / FreeSWITCH/PJSIP RTP) ──RTP──▶ SIP Edge

   Persistence:
   All services ──▶ MySQL (primary, with read replicas) + Redis (presence, ACD queues) + Kafka (CDR, events)
```

#### Component responsibilities

| Component | Language | Role |
| --- | --- | --- |
| **React app** | TS/React | Two SPAs sharing a component library: (a) **Admin console** for tenants and ops (CRUD on agents, skill groups, numbers, trunks, routing rules, IVR designer). (b) **Agent Workbench** with WebRTC softphone, presence, queues, ACW timer. |
| **API Gateway / BFF** | Go (`gin` or `echo`) | TLS termination, OIDC/JWT auth, RBAC, per-tenant rate limiting, exposes REST + WebSocket for the React apps. |
| **Agent / SkillGroup service** | Go | CRUD + presence broker. Holds the authoritative agent state machine in Redis (with MySQL as durable store). |
| **IVR / Routing service** | Go | Executes IVR flows (a small interpreter over a JSON DAG), evaluates skill-group queues, applies the auto-CLI selection algorithm, and applies trunk-routing rules. |
| **Call Control Service** | Go | The "brain". Receives events from SIP edges (`INVITE-IN`, `ANSWER`, `BYE`, `DTMF`), decides what to do, sends commands back (`PLAY`, `BRIDGE`, `TRANSFER`, `HANGUP`). Stateful, leader-elected per call. |
| **SIP Edge Node** | C++ (PJSIP) wrapped by a **Go sidecar** | Talks SIP+RTP to carriers and to the media gateway. Exposes a gRPC API to Call Control. **This is the integration boundary discussed in §2.4.** |
| **Media Gateway** | Reuse existing (RTPengine or FreeSWITCH or a PJSUA2 mixer) | Bridges WebRTC ⇄ SRTP/RTP ⇄ carrier RTP, handles transcoding (Opus ⇄ G.711). |
| **Reporting / CDR service** | Go | Consumes a Kafka stream of call events, materializes CDRs and aggregates into MySQL (read replicas) for dashboards. |

#### Inbound call flow (mirrors Aliyun's "呼入流转过程")

1. PSTN → carrier → carrier SIP trunk → **SIP Edge** (PJSIP) receives `INVITE`.
2. SIP Edge emits `InviteReceived(tenant, did, cli, callee, src_trunk)` over gRPC to **Call Control**.
3. Call Control looks up the DID in **Phone Number Management**, finds the bound IVR.
4. IVR executes; on "transfer to human" node it calls into the **Skill Group ACD** which selects an idle agent based on level + policy.
5. Call Control instructs the **Media Gateway** to bridge the carrier RTP leg to the agent's WebRTC leg; emits `Ringing` to the agent's Workbench over WebSocket.
6. Agent clicks **Answer** → Call Control sends `200 OK` upstream and connects the bridge.
7. On hang-up the SIP Edge fires `BYE`, Call Control finalizes the CDR, agent goes to ACW.

#### Outbound call flow (mirrors Aliyun's "呼出流转过程")

1. Workbench → REST `POST /calls/outbound {to, cli?}`.
2. Call Control runs **CLI selection** (dedicated → skill-group → tenant; same-city → same-prov → random).
3. Call Control runs **Trunk Routing**: ordered list of trunks matching `(cli_carrier, dest_prefix, time, tenant)`.
4. Call Control instructs the chosen **SIP Edge** to `INVITE` the callee. The Workbench is bridged via the media gateway.
5. CDR records both the chosen CLI and the egress trunk for billing reconciliation.

### 2.2 Frontend (React) — Workbench & Admin

**Workbench softphone:** Modeled after Aliyun's workbench states. Will **not** be a raw SIP-over-WebSocket client — instead it is a WebRTC peer to our media gateway, driven by REST/WebSocket commands. This is exactly the pattern Aliyun uses (workbench requires Chrome ≥ 58, i.e. WebRTC) and avoids shipping SIP credentials to the browser.

Key React modules:
- `presence/` — state machine identical to Aliyun's (Online / Idle / Ringing / Talking / ACW / Break / Invisible / Offline). The state lives server-side; the UI is a thin renderer of websocket-pushed snapshots.
- `softphone/` — `RTCPeerConnection` wrapper, audio-device picker, mute, DTMF (`createDTMFSender`), hold (renegotiate to `sendonly` + server-side hold music).
- `controls/` — Answer, Hang-up, Transfer (with skill-group / agent / external-number picker), Consult, Monitor.
- `dialer/` — Outbound dial pad with CLI selector (auto vs. manual), regional matching hint, redial.
- `acw/` — After-call work timer + wrap-up form (disposition codes, notes).
- `admin/` — Agents, Skill Groups, Numbers, Trunks, Trunk-Routing rules, IVR designer (drag-and-drop DAG editor), Reports.

Tech: React 18 + TypeScript, Vite, TanStack Query, Zustand for local UI state, RTK Query for API, WebSocket via `native` (no socket.io needed). Component library: Ant Design or shadcn/ui — Ant Design maps well to Aliyun's admin-console aesthetic if a similar look is desired.

### 2.3 Backend (Go) — Services and Boundaries

- HTTP framework: `gin` or `echo`; gRPC via `google.golang.org/grpc` for internal RPC.
- ORM: `sqlc` (compile-time-checked SQL) is preferred over GORM for a telephony system because every query is hot — but GORM is acceptable for CRUD-heavy services (agents, skill groups). Mixed usage is fine.
- Migration tooling: `golang-migrate`.
- Auth: OIDC (Keycloak or Auth0) → JWT bearer tokens; RBAC matrix mirroring Aliyun's three roles (agent, skill-group leader, admin) plus a `tenant_admin` super-role.
- Real-time: WebSocket for Workbench; internal events on NATS JetStream (low latency, persistent) and Kafka (for CDR fan-out + analytics).
- Presence + ACD queues: Redis (sorted-set per skill group by "last-idle timestamp" for longest-idle policy; per-agent hash for state).
- Observability: OpenTelemetry traces from edge → call control → media; Prometheus metrics; structured logs via `slog`.

### 2.4 PJSIP Integration — Approach and Trade-offs

PJSIP is a C library. Go cannot call it natively, so the integration choice is the single most important architectural decision in this project. Three viable options are evaluated below.

#### Option A — **Sidecar process** (RECOMMENDED)

PJSIP runs as a **standalone C++ binary** (a thin wrapper around `pjsua2`) on each SIP edge node. The Go SIP-edge service runs alongside it on the same host and communicates over **gRPC over a UNIX domain socket** (low-latency, no kernel TCP overhead) plus a shared-memory ring buffer for RTP statistics.

```
┌──────────────────────────────┐         ┌────────────────────────────┐
│   Go SIP-Edge sidecar        │ ◀──────▶│   pjsua2 C++ process       │
│   - gRPC server to Control   │  gRPC   │   - UAC/UAS                │
│   - REST /healthz, metrics   │  UDS    │   - account mgmt           │
│   - reconciles state with DB │         │   - account auth to trunks │
└──────────────────────────────┘         └────────────────────────────┘
                                                     │ SIP/UDP|TCP|TLS, RTP/SRTP
                                                     ▼
                                                  Carriers
```

**Pros:**
- **Crash isolation.** A PJSIP segfault does not kill the Go runtime. The Go side restarts the sidecar and gracefully fails over calls. This is the single biggest argument for the sidecar: PJSIP is mature but C — any memory bug terminates the whole process, and a Go process is the wrong thing to terminate.
- **Threading model is clean.** PJSIP has its own poll thread, its own worker threads, and its own assumptions about thread-local logging. Mixing those with Go's M:N scheduler is hostile. Keeping them in different OS processes ends the conversation.
- **GC pauses don't perturb SIP timers.** PJSIP's retransmit timers (Timer A, B, F, …) are tight (500 ms / 32 s). Go GC pauses are short but unbounded under load; a 200 ms STW on a busy host could cause a retransmit storm if PJSIP were embedded.
- **Independent upgrades.** PJSIP minor version bumps don't require recompiling Go.
- **Cross-platform builds in Go stay clean.** No CGo, no `pkg-config`, no librari­es to ship.

**Cons:**
- Two processes to deploy, monitor, supervise (use `systemd` or a sidecar container with shared PID namespace).
- IPC overhead per call event (~tens of microseconds over UDS — irrelevant relative to SIP RTTs in the tens of milliseconds).
- We must define a stable gRPC contract between Go and the C++ wrapper (specifying call events, commands, media controls).

#### Option B — **CGo bindings**

Link `libpjsua2` directly into the Go binary via CGo.

**Pros:** single binary, lowest IPC latency.

**Cons (these are dealbreakers for production):**
- CGo introduces a thread-switch penalty per call (~1µs) on **every** PJSIP callback. PJSIP fires callbacks from its own threads, so Go must convert them into goroutines via `cgo` — which acquires the runtime lock. Under load this becomes a contention hot-spot.
- PJSIP threads must NOT call into Go directly; the standard pattern is to push events onto a lock-free queue and have a Go goroutine drain it. This works but is fragile.
- A bug in PJSIP (or in our C glue) crashes the whole Go server, taking down all calls including signaling for unrelated tenants. **Operationally unacceptable.**
- Go's signal handlers conflict with PJSIP's (SIGRTMIN, SIGPIPE); needs careful `signal.Notify` setup.
- `pjsua2` is C++ — CGo bridging C++ requires an additional C shim layer.
- Statically linking PJSIP increases the binary by ~10 MB and forces all builds to have the PJSIP toolchain.

We **reject** this option for production. It can be used as a prototyping shortcut if needed.

#### Option C — **Replace PJSIP with a pure-Go SIP stack**

There are pure-Go SIP libraries (`pion/sip`, `emiago/sipgo`, `livekit/sip`). These avoid the language-boundary problem entirely.

**Pros:** single language, single binary, clean GC story.

**Cons:**
- The user explicitly specified PJSIP, so we keep PJSIP.
- Pure-Go SIP stacks are less battle-tested than PJSIP for edge cases (Chinese carrier quirks, B2BUA semantics, full ICE, SRTP-DTLS).

If at a later phase the team wants to migrate, the sidecar architecture makes the swap a single-binary replacement.

#### Recommendation: **Option A (sidecar)**.

#### Additional PJSIP integration challenges to plan for

1. **B2BUA, not proxy.** Aliyun CCC acts as a back-to-back user agent — it terminates the carrier leg and originates a separate leg to the agent. `pjsua2` natively supports two `Call` objects; we'll mirror that, with the Go side coordinating the bridge.
2. **Media handling.** PJSIP supports RTP/SRTP but **not** WebRTC-grade DTLS-SRTP + ICE end-to-end. The right pattern is:
   - PJSIP handles SIP/RTP towards carriers (G.711 a-/µ-law).
   - A dedicated WebRTC↔SIP media bridge (RTPengine, FreeSWITCH `mod_verto`, or a Pion-based custom bridge) handles the browser side.
   - PJSIP and the media bridge connect via plain RTP on a private interface.
3. **NAT traversal.** Carriers expect us at a fixed, public IP — per Aliyun's own docs requirement #3 — so SIP edges live on a public IP without NAT and PJSIP runs with `bound_addr` and `public_addr` configured. Floating IPs + BGP are the cleanest HA story.
4. **Resilience to PJSIP restarts.** Active calls in PJSIP are in-memory. On restart we lose call legs. Mitigations:
   - Use multiple SIP-edge nodes; carriers OPTIONS-probe and remove the failed one (mirrors the same OPTIONS-keepalive requirement Aliyun imposes on its customers).
   - Drain new calls before restart; let in-flight calls finish.
   - For long-running calls (>1h), advertise `Min-SE` headers and accept session refresh so the bridge can move legs.
5. **Capacity per node.** A single `pjsua2` instance on modern hardware handles roughly 500–1000 concurrent calls (signaling) and a lot less if it also relays RTP. We will **offload RTP** to the media gateway so the SIP edge is signaling-only and can scale much higher.
6. **TLS + SIPS.** Carriers in some regions require SIP-TLS + SRTP. PJSIP supports both but the certificates/SNI configuration is fiddly. We'll script certificate provisioning via cert-manager + a private CA.
7. **Logging / debugging.** PJSIP's own log is verbose; we'll set `pj_log_set_level(3)` in production and stream to a structured log file picked up by the Go sidecar for forwarding to Loki/Elasticsearch. Per-call SIP traces (`pjsua_call_dump`) are gated behind an admin "trace this call" toggle.
8. **License.** PJSIP is GPL/commercial dual-licensed. If we ship a closed-source product we'll need a commercial license from Teluu. If we ship the SIP edge as a standalone OSS-licensed component the GPL terms are easier to satisfy.

### 2.5 MySQL Data Schema (Call Routing + Agent Management)

Notes:
- Multi-tenant; every row carries `tenant_id`.
- Soft-delete via `deleted_at TIMESTAMP NULL` on mutable entities.
- `created_at`/`updated_at` everywhere; `BIGINT UNSIGNED` primary keys (Snowflake IDs from the app — friendly for sharding later).
- All textual IDs that Aliyun documents as "human-readable string IDs" (e.g. `skill_group_id`) are also stored as `VARCHAR(64)` `UNIQUE (tenant_id, code)` alongside the numeric PK.
- Engine: InnoDB. Charset `utf8mb4 / utf8mb4_0900_ai_ci`.
- `BIGINT` for E.164 numbers is **wrong** because of leading zeros and prefixes; use `VARCHAR(32)`.

The schema is broken into four logical clusters: **identity**, **telephony assets**, **routing**, **CDR / runtime**.

#### 2.5.1 Identity cluster

```sql
CREATE TABLE tenants (
  id              BIGINT UNSIGNED PRIMARY KEY,
  code            VARCHAR(64)  NOT NULL UNIQUE,
  display_name    VARCHAR(128) NOT NULL,
  status          ENUM('active','suspended','deleted') NOT NULL DEFAULT 'active',
  created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE users (
  id              BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  user_name       VARCHAR(64)  NOT NULL,           -- immutable, 4-32 chars (Aliyun rule)
  display_name    VARCHAR(128) NOT NULL,
  employee_id     VARCHAR(64),
  email           VARCHAR(255) NOT NULL,
  phone           VARCHAR(32),                     -- off-site / mobile-bridge
  landline        VARCHAR(32),
  password_hash   VARBINARY(255),                  -- nullable when SSO-only
  must_reset_pw   BOOLEAN NOT NULL DEFAULT TRUE,
  role            ENUM('agent','skill_group_leader','admin','tenant_admin') NOT NULL DEFAULT 'agent',
  work_mode       ENUM('on_site','off_site') NOT NULL DEFAULT 'on_site',
  status          ENUM('active','disabled','deleted') NOT NULL DEFAULT 'active',
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at      TIMESTAMP NULL,
  UNIQUE KEY uniq_tenant_username (tenant_id, user_name),
  INDEX idx_tenant_email (tenant_id, email),
  CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Agent-specific extension (kept separate so the user table stays generic)
CREATE TABLE agents (
  user_id         BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  extension       VARCHAR(16),                     -- internal dial-extension
  max_concurrent  TINYINT UNSIGNED NOT NULL DEFAULT 1,
  acw_seconds     INT UNSIGNED NOT NULL DEFAULT 30,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_agents_user FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX idx_agent_tenant (tenant_id)
);

-- Live presence is in Redis; this is the durable mirror for reporting only.
CREATE TABLE agent_presence_log (
  id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  user_id         BIGINT UNSIGNED NOT NULL,
  state           ENUM('offline','online','idle','ringing','talking','acw','break','invisible') NOT NULL,
  reason_code     VARCHAR(64),                     -- break reason: lunch, training, etc.
  entered_at      TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  left_at         TIMESTAMP(3) NULL,
  duration_ms     BIGINT GENERATED ALWAYS AS (TIMESTAMPDIFF(MICROSECOND, entered_at, left_at) / 1000) VIRTUAL,
  INDEX idx_user_time (user_id, entered_at),
  INDEX idx_tenant_state_time (tenant_id, state, entered_at)
);
```

#### 2.5.2 Skill-group cluster

```sql
CREATE TABLE skill_groups (
  id              BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  code            VARCHAR(64)  NOT NULL,            -- Aliyun's "技能组 ID", 4-64 chars
  name            VARCHAR(128) NOT NULL,
  description     VARCHAR(512),
  routing_policy  ENUM('longest_idle','least_busy','round_robin','priority') NOT NULL DEFAULT 'longest_idle',
  max_queue_size  INT UNSIGNED NOT NULL DEFAULT 100,
  max_wait_sec    INT UNSIGNED NOT NULL DEFAULT 300,
  overflow_target ENUM('voicemail','transfer','reject') NOT NULL DEFAULT 'voicemail',
  status          ENUM('active','disabled','deleted') NOT NULL DEFAULT 'active',
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_tenant_code (tenant_id, code),
  CONSTRAINT fk_sg_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE skill_group_members (
  skill_group_id  BIGINT UNSIGNED NOT NULL,
  user_id         BIGINT UNSIGNED NOT NULL,
  level           TINYINT UNSIGNED NOT NULL DEFAULT 5,  -- 1 highest, 10 lowest
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (skill_group_id, user_id),
  INDEX idx_user (user_id),
  CONSTRAINT fk_sgm_sg FOREIGN KEY (skill_group_id) REFERENCES skill_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_sgm_user FOREIGN KEY (user_id) REFERENCES users(id),
  CHECK (level BETWEEN 1 AND 10)
);
```

#### 2.5.3 Telephony-assets cluster

```sql
-- A "carrier" is the upstream telco we connect to.
CREATE TABLE carriers (
  id              BIGINT UNSIGNED PRIMARY KEY,
  code            VARCHAR(32) NOT NULL UNIQUE,    -- e.g. 'china_mobile', 'twilio_apac'
  display_name    VARCHAR(128) NOT NULL,
  region          VARCHAR(32),
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- A SIP trunk is a contracted signaling/media pipe to (or from) a carrier or to a customer-owned PBX.
CREATE TABLE sip_trunks (
  id                  BIGINT UNSIGNED PRIMARY KEY,
  tenant_id           BIGINT UNSIGNED NOT NULL,
  carrier_id          BIGINT UNSIGNED NOT NULL,
  code                VARCHAR(64) NOT NULL,
  direction           ENUM('inbound','outbound','both') NOT NULL DEFAULT 'both',

  -- Far end (what we send INVITEs to / receive INVITEs from)
  far_end_host        VARCHAR(255) NOT NULL,        -- public IP or FQDN
  far_end_port        SMALLINT UNSIGNED NOT NULL DEFAULT 5060,
  transport           ENUM('udp','tcp','tls') NOT NULL DEFAULT 'udp',
  sp_marker           VARCHAR(64),                  -- matches Aliyun's "SP 标示" field
  auth_username       VARCHAR(128),
  auth_password_enc   VARBINARY(255),               -- encrypted at rest
  realm               VARCHAR(255),

  -- Local end (what we present)
  local_bind_ip       VARCHAR(64),                  -- pinned for HA / floating IP
  options_keepalive_s INT UNSIGNED NOT NULL DEFAULT 30,
  max_concurrent      INT UNSIGNED NOT NULL DEFAULT 100,
  codecs              JSON NOT NULL,                -- ["PCMU","PCMA","G722","opus"]

  status              ENUM('active','disabled','deleted') NOT NULL DEFAULT 'active',
  health              ENUM('healthy','degraded','down','unknown') NOT NULL DEFAULT 'unknown',
  last_options_ok_at  TIMESTAMP NULL,

  created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uniq_tenant_code (tenant_id, code),
  INDEX idx_carrier (carrier_id),
  CONSTRAINT fk_trunk_carrier FOREIGN KEY (carrier_id) REFERENCES carriers(id),
  CONSTRAINT fk_trunk_tenant  FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Standby/failover legs for a trunk live in the same table referenced by group:
CREATE TABLE sip_trunk_groups (
  id              BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  code            VARCHAR(64) NOT NULL,
  description     VARCHAR(255),
  UNIQUE KEY uniq_tenant_code (tenant_id, code)
);

CREATE TABLE sip_trunk_group_members (
  group_id        BIGINT UNSIGNED NOT NULL,
  trunk_id        BIGINT UNSIGNED NOT NULL,
  priority        SMALLINT NOT NULL DEFAULT 100,    -- lower = earlier in failover order
  weight          SMALLINT NOT NULL DEFAULT 100,    -- for load balancing within same priority
  PRIMARY KEY (group_id, trunk_id),
  CONSTRAINT fk_stgm_group FOREIGN KEY (group_id) REFERENCES sip_trunk_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_stgm_trunk FOREIGN KEY (trunk_id) REFERENCES sip_trunks(id) ON DELETE CASCADE
);

-- Phone numbers / DIDs / Caller IDs
CREATE TABLE phone_numbers (
  id                  BIGINT UNSIGNED PRIMARY KEY,
  tenant_id           BIGINT UNSIGNED NOT NULL,
  number              VARCHAR(32) NOT NULL,         -- E.164 preferred
  display_label       VARCHAR(128),
  region_country      CHAR(2)  NOT NULL DEFAULT 'CN',
  region_province     VARCHAR(32),
  region_city         VARCHAR(64),
  carrier_id          BIGINT UNSIGNED,              -- which carrier owns it (for CLI integrity)
  usage_flags         SET('inbound','outbound') NOT NULL DEFAULT 'inbound,outbound',
  inbound_ivr_id      BIGINT UNSIGNED,              -- nullable; FK to ivr_flows
  digital_employee_id BIGINT UNSIGNED,              -- nullable; alternative to IVR
  group_label         VARCHAR(64),
  status              ENUM('active','disabled','deleted') NOT NULL DEFAULT 'active',
  created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_tenant_number (tenant_id, number),
  INDEX idx_carrier (carrier_id),
  CONSTRAINT fk_pn_tenant   FOREIGN KEY (tenant_id)   REFERENCES tenants(id),
  CONSTRAINT fk_pn_carrier  FOREIGN KEY (carrier_id)  REFERENCES carriers(id)
);

-- Number ↔ skill group (for outbound CLI eligibility)
CREATE TABLE phone_number_skill_groups (
  phone_number_id BIGINT UNSIGNED NOT NULL,
  skill_group_id  BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (phone_number_id, skill_group_id),
  CONSTRAINT fk_pnsg_pn FOREIGN KEY (phone_number_id) REFERENCES phone_numbers(id) ON DELETE CASCADE,
  CONSTRAINT fk_pnsg_sg FOREIGN KEY (skill_group_id)  REFERENCES skill_groups(id)  ON DELETE CASCADE
);

-- Number ↔ dedicated agent (a CLI an agent is allowed/encouraged to use)
CREATE TABLE phone_number_dedicated_agents (
  phone_number_id BIGINT UNSIGNED NOT NULL,
  user_id         BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (phone_number_id, user_id),
  CONSTRAINT fk_pnda_pn   FOREIGN KEY (phone_number_id) REFERENCES phone_numbers(id) ON DELETE CASCADE,
  CONSTRAINT fk_pnda_user FOREIGN KEY (user_id)         REFERENCES users(id)         ON DELETE CASCADE
);
```

#### 2.5.4 Routing cluster

```sql
-- IVR flow definition (DAG stored as JSON; reference-checked by the IVR service).
CREATE TABLE ivr_flows (
  id              BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  code            VARCHAR(64) NOT NULL,
  name            VARCHAR(128) NOT NULL,
  version         INT UNSIGNED NOT NULL DEFAULT 1,
  graph           JSON NOT NULL,                    -- node/edge definition; validated app-side
  status          ENUM('draft','published','archived') NOT NULL DEFAULT 'draft',
  published_at    TIMESTAMP NULL,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_tenant_code_version (tenant_id, code, version)
);

-- Trunk-routing rules: for outbound, choose the trunk group.
-- Rules are evaluated in `priority` order; first match wins.
CREATE TABLE trunk_routing_rules (
  id                 BIGINT UNSIGNED PRIMARY KEY,
  tenant_id          BIGINT UNSIGNED NOT NULL,
  name               VARCHAR(128) NOT NULL,
  priority           INT NOT NULL DEFAULT 100,         -- lower evaluated first

  -- Match clauses (all NULL = wildcard / always-match)
  match_cli_prefix          VARCHAR(32),               -- e.g. '+8613' for CMobile
  match_cli_carrier_id      BIGINT UNSIGNED,
  match_cli_country         CHAR(2),
  match_dest_prefix         VARCHAR(32),               -- dialed number prefix
  match_dest_country        CHAR(2),
  match_dest_region         VARCHAR(64),               -- province or city
  match_time_of_day_start   TIME,
  match_time_of_day_end     TIME,
  match_dow_mask            TINYINT UNSIGNED,          -- bitmask Mon=1..Sun=64
  match_skill_group_id      BIGINT UNSIGNED,
  match_tenant_cost_class   ENUM('premium','standard','grey'),

  -- Action
  target_trunk_group_id     BIGINT UNSIGNED NOT NULL,
  cli_rewrite_rule          VARCHAR(255),              -- optional sed-style: 's/+86/0086/'
  enabled                   BOOLEAN NOT NULL DEFAULT TRUE,

  created_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_tenant_priority (tenant_id, priority),
  INDEX idx_match_cli_prefix (tenant_id, match_cli_prefix),
  INDEX idx_match_dest_prefix (tenant_id, match_dest_prefix),
  CONSTRAINT fk_trr_target FOREIGN KEY (target_trunk_group_id) REFERENCES sip_trunk_groups(id),
  CONSTRAINT fk_trr_cli_carrier FOREIGN KEY (match_cli_carrier_id) REFERENCES carriers(id),
  CONSTRAINT fk_trr_skill_group FOREIGN KEY (match_skill_group_id) REFERENCES skill_groups(id)
);

-- CLI-selection rules: for an outbound, choose the CLI (calling number).
-- Models Aliyun's "auto" picker: dedicated → skill-group → tenant; same-city → same-prov → random.
-- Stored as ordered policy table so admins can tweak.
CREATE TABLE cli_selection_policies (
  id                       BIGINT UNSIGNED PRIMARY KEY,
  tenant_id                BIGINT UNSIGNED NOT NULL,
  name                     VARCHAR(128) NOT NULL,
  -- Ordered list of strategies, each evaluated in turn until a CLI is found:
  -- strategy = JSON like [{"scope":"dedicated"},
  --                       {"scope":"skill_group","match":"same_city"},
  --                       {"scope":"skill_group","match":"same_province"},
  --                       {"scope":"tenant","match":"any"}]
  strategy                 JSON NOT NULL,
  enabled                  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_tenant_name (tenant_id, name)
);
```

#### 2.5.5 Runtime / CDR cluster

```sql
-- One row per call leg. Inbound and outbound are both modeled as "calls" with two legs:
-- leg A (carrier-side) and leg B (agent-side). Joined by call_id.
CREATE TABLE calls (
  id                 CHAR(32)  PRIMARY KEY,         -- UUIDv4 hex, generated by Call Control
  tenant_id          BIGINT UNSIGNED NOT NULL,
  direction          ENUM('inbound','outbound','internal') NOT NULL,
  cli                VARCHAR(32),
  callee             VARCHAR(32),
  did                VARCHAR(32),                   -- our number for inbound
  agent_user_id      BIGINT UNSIGNED,               -- nullable until assigned
  skill_group_id     BIGINT UNSIGNED,
  ingress_trunk_id   BIGINT UNSIGNED,               -- which trunk inbound came in on
  egress_trunk_id    BIGINT UNSIGNED,               -- which trunk outbound went out on
  ivr_flow_id        BIGINT UNSIGNED,
  start_at           TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  answer_at          TIMESTAMP(3) NULL,
  end_at             TIMESTAMP(3) NULL,
  hangup_cause       VARCHAR(64),                   -- carrier-equivalent + 'agent_hangup' / 'caller_hangup' / 'timeout'
  recording_url      VARCHAR(512),
  disposition_code   VARCHAR(64),                   -- agent-supplied wrap-up code
  notes              TEXT,
  created_at         TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_tenant_time (tenant_id, start_at),
  INDEX idx_agent_time  (agent_user_id, start_at),
  INDEX idx_skill_time  (skill_group_id, start_at)
);

CREATE TABLE call_events (
  id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  call_id         CHAR(32) NOT NULL,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  event_type      VARCHAR(48) NOT NULL,             -- 'invite_received','ringing','answer','transfer','hold','bye',...
  payload         JSON,
  occurred_at     TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  INDEX idx_call_time (call_id, occurred_at),
  CONSTRAINT fk_ce_call FOREIGN KEY (call_id) REFERENCES calls(id) ON DELETE CASCADE
);

-- ACD queue snapshots — written periodically for reporting; live data is in Redis.
CREATE TABLE queue_snapshots (
  id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  skill_group_id  BIGINT UNSIGNED NOT NULL,
  waiting_count   INT UNSIGNED NOT NULL,
  longest_wait_s  INT UNSIGNED NOT NULL,
  agents_idle     INT UNSIGNED NOT NULL,
  agents_busy     INT UNSIGNED NOT NULL,
  captured_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_sg_time (skill_group_id, captured_at)
);
```

#### 2.5.6 Indexing & sharding notes

- **Hot tables** are `calls`, `call_events`, `agent_presence_log`. They are append-mostly and grow unbounded; partition by `start_at` / `occurred_at` (monthly RANGE partitions) and archive older partitions to S3/OSS.
- `phone_numbers` carries the highest read frequency on the routing path (`SELECT … WHERE tenant_id=? AND number=?`); covered by the unique key.
- `trunk_routing_rules` is small (hundreds of rows per tenant) — keep entirely in memory at the routing service and reload on changes via a `pubsub` channel.
- Replication: one primary, one synchronous replica for HA, one or more async read replicas for reports.

### 2.6 Deployment Topology

- **3+ SIP edge nodes** in front of carriers, behind anycast public IPs.
- **2+ media gateways** (RTPengine cluster) for WebRTC ↔ RTP bridging.
- **Stateless Go services** behind an internal L7 LB; horizontal autoscaling on CPU + Redis queue depth.
- **MySQL** 1 primary + 1 sync replica + N async replicas.
- **Redis** cluster (3 masters + 3 replicas) for presence and ACD queues.
- **Kafka / NATS** for events and CDR fan-out.
- **CI/CD:** GitHub Actions → container builds → ArgoCD on Kubernetes; PJSIP sidecar image built separately so its release cycle is decoupled.

### 2.7 Phased delivery plan (proposed)

| Phase | Scope | Exit criteria |
| --- | --- | --- |
| 0. Foundation | Repos, CI, infra-as-code, MySQL schema migrations, auth, tenants, users | Devs can sign up, log in, manage agents/skill groups |
| 1. Inbound MVP | One SIP edge, one trunk, basic IVR (play + transfer-to-skill-group), workbench answer/hang-up | A real PSTN call rings an agent and gets answered |
| 2. Outbound MVP | Outbound dial, CLI selection (auto + manual), trunk routing rules | Agent can dial out with correct CLI |
| 3. Workbench v1 | Hold, mute, transfer (cold), ACW, presence states, break reasons | Feature parity with Aliyun's basic workbench |
| 4. Trunks v2 | Multiple trunks, OPTIONS keepalive, failover, health-based eviction | Surviving a planned trunk outage |
| 5. Workbench v2 | Warm transfer, consult, monitor/whisper/barge, CSAT | Supervisor flows complete |
| 6. Reporting + Recording | CDR analytics, call recording, agent KPIs | Dashboards live |
| 7. Scale & Hardening | HA failover drills, load tests (target: 10k concurrent calls), TLS+SRTP on trunks | Pen-test pass; SLO 99.9 % |

### 2.8 Top architectural risks

1. **PJSIP ↔ Go boundary** — addressed by the sidecar pattern (§2.4) but worth a dedicated spike in Phase 1.
2. **WebRTC media bridging** — RTPengine is the well-known answer in OSS; if license/operational concerns require an in-house bridge we will scope that as Phase 4 work, not MVP.
3. **CLI integrity per carrier** — Chinese carriers will rewrite CLIs that don't match the egress trunk's carrier; the trunk-routing rules must enforce CLI/carrier matching or call delivery silently degrades.
4. **Data residency for cross-border deployments** — recording and CDR storage need per-region MySQL clusters; not in MVP scope but the schema is multi-tenant from day one to allow regional sharding.
5. **License of PJSIP** — get explicit clarification from Teluu before shipping; budget for the commercial license if the product is closed-source.

---

## What I'd like your approval on before coding

1. **The component decomposition in §2.1.** In particular, the split between a Go control plane and a C++/PJSIP sidecar with a separate WebRTC media bridge.
2. **The PJSIP integration approach in §2.4** — sidecar over CGo over pure-Go-SIP. If you have a constraint that forces a single-binary deployment we should revisit.
3. **The schema in §2.5.** Specifically the choices to keep:
   - Presence in Redis (with MySQL as a durable log only),
   - `trunk_routing_rules` as a *table* of typed match clauses rather than a free-form expression,
   - `cli_selection_policies` as a JSON-encoded ordered strategy list (so admins can tweak without code changes).
4. **The phased plan in §2.7.** Mainly: are you comfortable with an Inbound-first MVP before outbound, or do you want them in parallel?

I have made **no code changes** and opened **no PR**. Once you approve (or send back changes), I'll start with Phase 0 / Phase 1 scaffolding.
