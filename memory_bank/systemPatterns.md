# System Patterns: gosafejson

## Architectural Patterns

### 1. Modular Design Pattern
- **Core Modules**: `any_*`, `iter_*`, `reflect_*`, `stream_*`
- **Extension Modules**: `extra/` directory with specialized codecs
- **Test Separation**: Dedicated test directories by functionality

### 2. Type System Architecture
```
JSON Value Handling:
├── any_*.go       # Dynamic type wrapper system
├── iter_*.go      # Iterator pattern for parsing
├── reflect_*.go   # Reflection-based type mapping
└── stream_*.go    # Stream processing layer
```

### 3. Configuration Pattern
- **Central Config**: `config.go` with `399 lines` - comprehensive configuration
- **Adapter Pattern**: `adapter.go` for interface compatibility
- **Pool Pattern**: `pool.go` for resource management

### 4. Error Handling Pattern
- **Type-safe Errors**: Structured error handling across modules
- **Validation**: Multiple validation layers in parsing
- **Graceful Degradation**: Fallback mechanisms for invalid data

## Code Organization Patterns

### File Naming Convention
```
Prefix-based organization:
- any_*      : Dynamic value handling (11 files)
- iter_*     : Iterator implementations (8 files)  
- reflect_*  : Reflection utilities (10 files)
- stream_*   : Stream processing (4 files)
```

### Test Structure Pattern
```
Functional test grouping:
├── any_tests/       # Dynamic value tests
├── api_tests/       # API interface tests
├── benchmarks/      # Performance tests
├── extension_tests/ # Extension functionality tests
├── misc_tests/      # General functionality tests
├── skip_tests/      # Skip logic tests
├── type_tests/      # Type system tests
└── value_tests/     # Value handling tests
```

## Design Patterns

### 1. Factory Pattern
- **Config Factory**: Multiple configuration constructors
- **Iterator Factory**: Different iterator types for different JSON structures
- **Codec Factory**: Custom encoder/decoder creation

### 2. Strategy Pattern
- **Encoding Strategies**: Multiple encoding approaches
- **Parsing Strategies**: Different parsing methods for different data types
- **Validation Strategies**: Configurable validation approaches

### 3. Adapter Pattern
- **Interface Compatibility**: `adapter.go` bridges different interfaces
- **Standard Library Integration**: Compatibility with `encoding/json`

### 4. Pool Pattern
- **Object Pooling**: `pool.go` manages resource reuse
- **Memory Optimization**: Reduces garbage collection pressure

## Extension Patterns

### Custom Codec Pattern
```go
// Pattern for extending functionality
├── binary_as_string_codec.go    # Binary data as string
├── time_as_int64_codec.go       # Time as integer
├── fuzzy_decoder.go             # Lenient parsing
└── naming_strategy.go           # Field naming flexibility
```

### Plugin Architecture
- **Extensible Design**: Easy to add new codecs
- **Configuration-driven**: Runtime behavior modification
- **Backwards Compatible**: Non-breaking extensions

## Performance Patterns

### 1. Streaming Pattern
- **Incremental Processing**: Memory-efficient large data handling
- **Lazy Evaluation**: Parse only what's needed
- **Buffer Management**: Efficient I/O operations

### 2. Caching Pattern
- **Reflection Caching**: Expensive reflection operations cached
- **Type Information**: Reuse of type metadata
- **String Interning**: String deduplication

### 3. Zero-allocation Pattern
- **Pool Reuse**: Object pooling to minimize allocations
- **In-place Operations**: Modify data without copying
- **Streaming**: Avoid loading entire documents 
