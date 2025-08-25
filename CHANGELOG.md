# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of SafeMap
- Thread-safe map implementation using goroutines and channels
- Generic support for any comparable key type and any value type
- Iterator support for keys and key-value pairs using Go 1.24+ iter package
- Comprehensive API with Set, Get, Delete, Exist, Keys, All, Length, and GetMap methods
- Complete documentation and examples
- Performance benchmarking examples
- Concurrent usage examples

### Features

- **Thread Safety**: All operations are safe for concurrent use
- **Type Safety**: Full generic type support with compile-time type checking
- **Iterator Support**: Native support for range iteration over keys and values
- **Channel-based Architecture**: Uses Go channels for internal synchronization
- **Zero Dependencies**: No external dependencies except for testing
- **Panic Safety**: Clear error messages for common mistakes

### Documentation

- Comprehensive README with usage examples
- Detailed API documentation
- Performance comparison with sync.Map and mutex-based maps
- Multiple example programs demonstrating different use cases
- Benchmarking tools for performance analysis

## Design Decisions

### Why Channels Over Mutexes?

- **Simplicity**: No need to manage complex locking scenarios
- **Deadlock Prevention**: Channel-based design eliminates deadlock possibilities
- **Go Idioms**: Follows Go's philosophy of "Don't communicate by sharing memory; share memory by communicating"
- **Sequential Consistency**: All operations are processed in a single goroutine, ensuring strong consistency

### Why Generic Types?

- **Type Safety**: Compile-time type checking prevents runtime type assertion errors
- **Performance**: No boxing/unboxing overhead compared to interface{} based solutions
- **Developer Experience**: Better IDE support and code completion

### API Design Principles

- **Familiarity**: Method names and behavior similar to standard Go maps
- **Consistency**: All methods follow similar patterns and error handling
- **Safety**: Clear panic messages for misuse
- **Extensibility**: Iterator support for future enhancements
