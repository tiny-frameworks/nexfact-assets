<sup>🌍 **Sprache:** 🇬🇧 [Englisch →](README.md)</sub>

[![Codeberg Release](https://img.shields.io/codeberg/v/release/tiny-frameworks/nexfact-assets?logo=codeberg&logoColor=white)](https://codeberg.org/tiny-frameworks/nexfact-assets)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

| | |
|---|---|
|![GSF-NexFact-Logo](nexfact.png)| ***GSF-NexFact-Assets***<br>Ein dediziertes zentrales Testressourcen Repository für das **GSF-NexFact** e-invoicing ecosystem.|

<sup>***GSF*** steht für ***Go Small Frameworks*** — minimalistische Tools für robuste Applicationen.</sup>

---

### Repository

Das Projekt wird auf Codeberg entwickelt und gewartet.

**Canonical repository:**

> https://codeberg.org/tiny-frameworks/nexfact-assets

GitHub ist ein automatisierter Mirror des Codeberg Repositories

**GitHub mirror:**

> https://github.com/tiny-frameworks/nexfact-assets

---

## GSF-NexFact-Assets

`GSF-Nexfact-Assets` ist das zentrale Test-Asset-Repository für das **GSF-nexfact** E-Invoicing-Ökosystem. Es enthält gemeinsam genutzte Testumgebungen (`env`), Referenz-Ressourcen (Schemas, Templates, Mock-Profile) sowie Beispiel-Rechnungsdateien, die von `nexfact`, `nexfact-server` und verwandten Microservices verwendet werden.

Durch die Auslagerung großer binärer Testdateien in dieses Repository bleiben die Kern-Repositories (`nexfact`, `nexfact-server`, `nexutils`) schlank und schnell zu klonen.

Dieses Projekt ist Teil der **tiny-frameworks** Initiative.

---

## Verzeichnisstruktur

```text
nexfact-assets/
├── setup.go             # Erstellt isolierte Arbeitsverzeichnisse & freie dynamische Ports
├── env/                 # nexfact Umgebungs-Templates & Schemas
└── testdata/
    └── resources/       # Test-Eingabedateien (PDFs, XMLs, JSON-Payloads, Bilder)
        ├── test_generate_fx.json
        ├── test_combine_invoice.pdf
        ├── test_combine_invoice.xml
        └── test_attachment.png

```

---

## Testausführung & Isolations-Muster

Um sicherzustellen, dass Tests deterministisch und frei von Seiteneffekten bleiben, **dürfen Tests niemals direkt in dieses Asset-Repository schreiben**.

Stattdessen folgen die Tests dem **TempDir-Isolationsmuster**:

```text
┌───────────────────────────┐
│      nexfact-assets       │ (Schreibgeschützte Master-Templates & Fixtures)
└─────────────┬─────────────┘
              │ 1. Kopieren bei Teststart
              ▼
┌───────────────────────────┐
│        tempDir            │ (Isoliertes temporäres Arbeitsverzeichnis)
│  ├── env/                 │
│  └── resources/           │
└─────────────┬─────────────┘
              │ 2. Verarbeiten, Ausgaben schreiben & Prüfen
              ▼
┌───────────────────────────┐
│     os.RemoveAll()        │ (Manuelles/Verzögertes Aufräumen durch Aufrufer)
└───────────────────────────┘

```

---

## Verwendung in Tests

```go
import "codeberg.org/tiny-frameworks/nexfact-assets"

func TestMyServer(t *testing.T) {
    setupMap, err := assets.SetupTestEnv("native")
    if err != nil {
        t.Fatalf("Fehler beim Erstellen der Testumgebung: %v", err)
    }

    envRoot := setupMap["envRoot"]
    resourcePath := setupMap["resourcePath"]
    tempDir := setupMap["tempDir"]
    defer os.RemoveAll(tempDir)

    // ... Tests ausführen ...
}

```

---

## Lokale Entwicklung

Dieses Repository ist für die gemeinsame Nutzung mit anderen `tiny-frameworks`-Modulen (`nexfact`, `nexfact-server`, `nexutils`) ausgelegt.

Für die lokale Entwicklung empfiehlt es sich, die Repositories nebeneinander in einem gemeinsamen Workspace-Verzeichnis abzulegen.

---

## Organisation & Standards

* **Copyright:** © 2026 Georg Hagn.
* **Repository:** `codeberg.org/tiny-frameworks/nexfact-assets`
* **Lizenz:** Apache License, Version 2.0.

*GSF-NexFact-Assets ist ein unabhängiges Open-Source-Projekt und steht in keiner Verbindung zu Unternehmen mit ähnlichem Namen.*

---

## Kontakt & Support

Für Anfragen, architektonische Diskussionen oder Sicherheitsfragen wende dich bitte an:

📧 **georghagn [at] tiny-frameworks.io**

---
