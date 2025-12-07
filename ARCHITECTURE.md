# Architecture Documentation

## Overview

This project uses a **Processor Pattern** architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Commands                           │
│                   (convert, removepoints)                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    App (Facade)                             │
│         Orchestrates the conversion process                 │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│            CoordinatesProcessor (Main Processor)            │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────┐  │
│  │   Reader     │ ──→  │  Transform   │ ──→  │  Writer  │  │
│  │  (Reading)   │      │  (Processing)│      │ (Writing)│  │
│  └──────────────┘      └──────────────┘      └──────────┘  │
└─────────────────────────────────────────────────────────────┘
         ▲                                           ▲
         │                                           │
         │                                           │
    ┌────┴─────────────────────┐         ┌──────────┴─────────┐
    │                          │         │                    │
    ▼                          ▼         ▼                    ▼
┌──────────────┐        ┌─────────────┐ ┌─────────────┐
│ExcelReader   │        │   Models    │ │GeoJsonWriter│
│  - Reads     │        │  - CordsData│ │  - Writes   │
│    Excel     │        │  - Data     │ │    GeoJSON  │
│    files     │        │    classes  │ │  - Manages  │
└──────────────┘        └─────────────┘ │    features │
                                        └─────────────┘
```

## Directory Structure

```
internal/domain/
├── models/                    # Data models
│   └── cords_data.go         # CordsData structure (latitude, longitude, etc)
│
├── readers/                   # Input readers (Reader interface implementations)
│   ├── reader.go             # Reader interface definition
│   └── excel_reader.go       # ExcelReader - reads from Excel files
│
├── writers/                   # Output writers (Writer interface implementations)
│   ├── writer.go             # Writer interface definition
│   └── geojson_writer.go     # GeojsonWriter - writes to GeoJSON format
│
├── processors/               # Business logic processors
│   └── coordinates_processor.go  # CoordinatesProcessor - orchestrates Read→Process→Write
│
└── app/
    └── app.go               # App facade - initializes components and calls processor
```

## Core Components

### 1. **Reader Interface** (`readers/reader.go`)
Defines the contract for reading coordinate data from any source.

```go
type Reader interface {
    Read() (*[]models.CordsData, error)
    Close() error
}
```

**Implementations:**
- `ExcelReader` - Reads coordinates from Excel files with configurable sheet, columns, and row start

### 2. **Writer Interface** (`writers/writer.go`)
Defines the contract for writing coordinate data to any format.

```go
type Writer interface {
    Write(data *[]models.CordsData, color string) error
    Save(path string) error
    Close() error
}
```

**Implementations:**
- `GeojsonWriter` - Writes coordinates to GeoJSON FeatureCollection format

### 3. **CoordinatesProcessor** (`processors/coordinates_processor.go`)
Main orchestrator that combines Reader and Writer:

1. Reads data using the Reader
2. Processes it (validation, transformation)
3. Writes it using the Writer

The processor is **format-agnostic** - it doesn't care if input is Excel or CSV, or if output is GeoJSON or KML.

### 4. **Models** (`models/cords_data.go`)
Data structures that represent coordinates:
- `CordsData` - Contains latitude/longitude, caption, description

### 5. **App** (`app/app.go`)
High-level facade that:
- Loads configuration
- Creates Reader and Writer instances
- Creates and runs the Processor
- Handles cleanup

## Data Flow

```
┌─────────────────────────┐
│   Excel File            │
│  (Sheet1, Columns A-C)  │
└────────────┬────────────┘
             │
             ▼
      ┌─────────────┐
      │ExcelReader  │
      │  .Read()    │
      └────────────┬┘
                   │
                   ▼
        ┌────────────────────┐
        │  CordsData Array   │
        │  [                 │
        │    {lat, lon, ...}│
        │    {lat, lon, ...}│
        │  ]                 │
        └────────────────────┘
                   │
                   ▼
      ┌─────────────────────┐
      │CoordinatesProcessor │
      │  .Process()         │
      └──────────┬──────────┘
                 │
                 ▼
        ┌──────────────────┐
        │GeojsonWriter     │
        │ .Write()         │
        │ .Save()          │
        └──────────┬───────┘
                   │
                   ▼
        ┌──────────────────┐
        │  GeoJSON File    │
        │  (FeatureCol.)   │
        │ [                │
        │   Features...    │
        │ ]                │
        └──────────────────┘
```

## Key Design Patterns

### 1. **Processor Pattern**
The main flow is: **Read → Process → Write**
Each step is independent and can be replaced.

### 2. **Interface Segregation**
- Reader and Writer are separate interfaces
- Easy to add new readers (CSV, JSON, Database)
- Easy to add new writers (KML, GeoJSON Lines, etc)

### 3. **Dependency Injection**
- Processor receives Reader and Writer as dependencies
- App creates instances and passes them to Processor
- No hard-coded dependencies

### 4. **Facade Pattern**
- App acts as a facade
- Clients (CLI commands) interact with App, not individual components
- Simplifies the public API

## Adding New Components

### To add a new Reader (e.g., CSVReader):

1. Create `readers/csv_reader.go`
2. Implement the `Reader` interface:
   ```go
   type CSVReader struct { ... }
   func (r *CSVReader) Read() (*[]models.CordsData, error) { ... }
   func (r *CSVReader) Close() error { ... }
   ```
3. Use it in Processor:
   ```go
   reader := readers.NewCSVReader("data.csv")
   processor := processors.NewCoordinatesProcessor(reader, writer)
   ```

### To add a new Writer (e.g., KMLWriter):

1. Create `writers/kml_writer.go`
2. Implement the `Writer` interface:
   ```go
   type KMLWriter struct { ... }
   func (w *KMLWriter) Write(data *[]models.CordsData, color string) error { ... }
   func (w *KMLWriter) Save(path string) error { ... }
   func (w *KMLWriter) Close() error { ... }
   ```
3. Use it in Processor:
   ```go
   writer := writers.NewKMLWriter("template.kml")
   processor := processors.NewCoordinatesProcessor(reader, writer)
   ```

## Configuration

Configuration is loaded from YAML files and passed to App:
- Excel file path, sheet name, column mappings
- GeoJSON input/output paths
- Appearance settings (marker color, etc)

See `config.example.yaml` for structure.

## Error Handling

- Each component validates its inputs
- ExcelReader validates sheet and column names on creation
- Errors are wrapped with context using `fmt.Errorf`
- App properly closes all resources in `defer` blocks