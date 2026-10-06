<sup>🌍 **Language:** 🇩🇪 [German →](README.de.md)</sup>

[![Codeberg Release](https://img.shields.io/codeberg/v/release/tiny-frameworks/nexfact-assets?logo=codeberg&logoColor=white)](https://codeberg.org/tiny-frameworks/nexfact-assets)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

| | |
|---|---|
|![GSF-NexFact-Logo](nexfact.png) | ***GSF-NexFact-Assets***<br>A dedicated central test asset repository for the **GSF-NexFact** e-invoicing ecosystem.|

<sup>***GSF*** stands for ***Go Small Frameworks*** — minimalist tools for robust applications.</sup>

---

### Repository

This project is developed and maintained on Codeberg.

**Canonical repository:**

> https://codeberg.org/tiny-frameworks/nexfact-assets

GitHub is an automated mirror of the Codeberg repository:

**GitHub mirror:**

> https://github.com/tiny-frameworks/nexfact-assets

---

## GSF-NexFact-Assets

`GSF-nexfact-assets` is the dedicated central test asset repository for the **GSF-nexfact** e-invoicing ecosystem. It stores shared test environments (`env`), reference resources (schemas, templates, mock profiles), and sample invoice files used across `nexfact`, `nexfact-server`, and related microservices.

By isolating heavy binary test fixtures into this repository, core source repositories (`nexfact`, `nexfact-server`, `nexutils`) remain lightweight and fast to clone.

This project is part of the **tiny-frameworks** initiative.

---

## Directory Structure

```text
nexfact-assets/
├── setup-go             # creates isolated working-directory and tcp-addresses
├── env/                 # Reference nexfact environment templates & schemas
└── testdata/
    └── resources/       # Test input files (PDFs, XMLs, JSON payloads, images)
        ├── test_generate_fx.json
        ├── test_combine_invoice.pdf
        ├── test_combine_invoice.xml
        └── test_attachment.png

```

## Test Execution & Isolation Pattern

To guarantee that tests remain strictly deterministic and side-effect-free, **tests must never write directly to this asset repository**.

Instead, tests follow the **TempDir Isolation Pattern**:

```text
┌───────────────────────────┐
│     nexfact-assets        │ (Read-Only Master Templates & Fixtures)
└─────────────┬─────────────┘
              │ 1. Copy on Test Startup
              ▼
┌───────────────────────────┐
│       t.TempDir()         │ (Isolated Working Directory)
│  ├── env/                 │
│  └── resources/           │
└─────────────┬─────────────┘
              │ 2. Process, Write Output & Assert
              ▼
┌───────────────────────────┐
│   Automatic OS Cleanup    │ (Triggered by Go test framework upon completion)
└───────────────────────────┘

```

---

## Local Development (`go.work`)

This repository is designed to be developed as part of a multi-module workspace alongside related `tiny-frameworks` projects. 

To set up a local development environment, place the repositories side-by-side in a common workspace directory and initialize a `go.work` file:

```bash
mkdir workspace && cd workspace
git clone [https://codeberg.org/tiny-frameworks/nexfact.git](https://codeberg.org/tiny-frameworks/nexfact.git)
git clone [https://codeberg.org/tiny-frameworks/nexfact-server.git](https://codeberg.org/tiny-frameworks/nexfact-server.git)
git clone [https://codeberg.org/tiny-frameworks/nexutils.git](https://codeberg.org/tiny-frameworks/nexutils.git)
git clone [https://codeberg.org/tiny-frameworks/nexfact-assets.git](https://codeberg.org/tiny-frameworks/nexfact-assets.git)

go work init ./nexfact ./nexfact-server ./nexutils

```

This setup ensures that Go tools (and tests) resolve local inter-module dependencies seamlessly without needing local `replace` directives in `go.mod`.

---

## Organizational & Standards

* **Copyright:** © 2026 Georg Hagn.
* **Repository:** `codeberg.org/tiny-frameworks/nexfact-assets`
* **License:** Apache License, Version 2.0.

*GSF-NexFact-Assets is an independent open-source project and is not affiliated with any corporation of a similar name.*

---

## Contact & Support

For inquiries, architectural discussions, or security issues, please contact:

📧 **georghagn [at] tiny-frameworks.io**

---
