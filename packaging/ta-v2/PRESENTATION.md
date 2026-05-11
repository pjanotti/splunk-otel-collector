# Splunk_TA_otel v2 — Technical Overview

## Agenda

1. [What is a Splunk TA / Modular Input?](#1-what-is-a-splunk-ta--modular-input)
2. [How `splunkd` Starts a Modular Input](#2-how-splunkd-starts-a-modular-input)
3. [The Three Execution Modes](#3-the-three-execution-modes)
4. [Configuration Options for `Splunk_TA_otel`](#4-configuration-options-for-splunk_ta_otel)

---

## 1. What is a Splunk TA / Modular Input?

### Terminology

| Term | Meaning |
| ------ | --------- |
| **Splunk App** | Full UI application inside Splunk |
| **Splunk Add-on (TA)** | Packaged extension — limited UI, adds data inputs or knowledge |
| **Modular Input** | The script/binary inside a TA that collects data |
| **`splunkd`** | The Splunk host process — parent of all modular inputs |

---

### What Does `Splunk_TA_otel` v2 Do?

- Deploys the **Splunk OpenTelemetry Collector** alongside Splunk Universal Forwarder
- Replacement for the older `Splunk_TA_otel` v1 that used scripts (.sh, .cmd, .ps1) to manage a standalone OTel Collector
- Runs as a long-lived child process of `splunkd`
- Same executable — behaves slightly differently when managed by `splunkd` vs run standalone (OTel Collector).

---

### Package Format

A TA is shipped as a `.tgz` file. The root directory name **must** match the modular input name:

```
Splunk_TA_otel/
├── default/
│   ├── app.conf          ← app metadata, version
│   └── inputs.conf       ← default configuration
├── local/
│   └── inputs.conf       ← admin overrides (takes precedence)
├── README/
│   └── inputs.conf.spec  ← parameter schema / documentation
├── configs/
│   └── agent_config.yaml ← OTel collector config
├── static/
│   ├── appIcon.png
│   └── appIcon_2x.png
├── linux_x86_64/bin/
│   └── Splunk_TA_otel        ← Linux binary
└── windows_x86_64/bin/
    └── Splunk_TA_otel.exe    ← Windows binary
```

- The binary lives in **platform/architecture-specific paths**
- Only `x86_64` is supported for Splunk modular inputs
- Install by extracting the `.tgz` into `$SPLUNK_HOME/etc/apps/`

---

### TA-Specific Parameters (`inputs.conf.spec`)

Defined in [assets/README/inputs.conf.spec](assets/README/inputs.conf.spec):

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `splunk_access_token` | **Yes** | _(empty)_ | Access token for Splunk Observability Cloud |
| `splunk_realm` | **Yes** | _(empty)_ | Observability Cloud realm, e.g. `us0`, `eu0` |
| `splunk_config` | Yes | `$SPLUNK_HOME/etc/apps/Splunk_TA_otel/configs/agent_config.yaml` | Path to OTel Collector config YAML |
| `splunk_collector_env_vars` | No | _(empty)_ | Extra env vars: `KEY1=VAL1,KEY2=VAL2` |
| `splunk_collector_cmd_args` | No | _(empty)_ | Extra CLI args passed to the collector |

---

### `splunk_collector_env_vars` Encoding Rules

The comma (`,`) is the delimiter between key=value pairs. If a value contains a literal `,` or `=`, percent-encode them:

| Character | Encoded form |
|-----------|--------------|
| `,` (comma) | `%2C` |
| `=` (equals) | `%3D` |

**Example:**

```ini
splunk_collector_env_vars = MY_URL=https://host%3A8080/path,MY_TAG=a%2Cb
```

Resolves to env vars:
- `MY_URL=https://host:8080/path`
- `MY_TAG=a,b`

---

### `agent_config.yaml` — OTel Collector Configuration

Located at [assets/configs/agent_config.yaml](assets/configs/agent_config.yaml).

The YAML references environment variables set by the TA:

```yaml
extensions:
  health_check:
    endpoint: "${SPLUNK_LISTEN_INTERFACE}:13133"

receivers:
  hostmetrics:
    collection_interval: 10s
    scrapers:
      cpu:
      disk:
      memory:
      network:

exporters:
  sapm:
    access_token: "${SPLUNK_ACCESS_TOKEN}"
    endpoint: "https://ingest.${SPLUNK_REALM}.signalfx.com/v2/trace"
  signalfx:
    access_token: "${SPLUNK_ACCESS_TOKEN}"
    realm: "${SPLUNK_REALM}"
```

- `splunk_access_token` → `SPLUNK_ACCESS_TOKEN`
- `splunk_realm` → `SPLUNK_REALM`
- Additional vars injectable via `splunk_collector_env_vars`

---

## Summary

```
┌──────────────────────────────────────────────────────────────────┐
│                    Configuration Flow                            │
│                                                                  │
│  inputs.conf (local overrides default)                           │
│       │                                                          │
│       ▼                                                          │
│  splunkd reads + merges stanza params                            │
│       │                                                          │
│       ▼                                                          │
│  sends <input> XML via stdin to binary                           │
│       │                                                          │
│       ▼                                                          │
│  binary maps splunk_* params → env vars + cmd args               │
│       │                                                          │
│       ▼                                                          │
│  OTel Collector starts using agent_config.yaml                   │
│  (YAML reads env vars set by binary)                             │
└──────────────────────────────────────────────────────────────────┘
```

---

## References

- [Splunk Modular Inputs Overview](https://dev.splunk.com/enterprise/docs/developapps/manageknowledge/custominputs/modinputsoverview/)
- [Modular Inputs Examples](https://dev.splunk.com/enterprise/docs/developapps/manageknowledge/custominputs/modinputsexamples/)
- [inputs.conf Reference](https://docs.splunk.com/Documentation/Splunk/latest/Admin/Inputsconf)
- [TA inputs.conf.spec](assets/README/inputs.conf.spec)
- [dotnet-instr-deployer-add-on Primer](../dotnet-instr-deployer-add-on/README.md#splunk-modular-input-primer)
