# CSMap - Concurrent Shards HashMap in Go

CSMap is a thread-safe concurrent hashmap implementation in Go, designed to offer efficient and simultaneous access by partitioning the map into several shards. Each shard operates independently, allowing for reduced contention among concurrent goroutines.

## Features

- **Concurrent Access**: Utilizes mutexes to ensure thread-safe operations when adding, retrieving, or deleting elements.
- **Sharded Design**: Splits the hashmap into multiple shards, allowing for greater scalability and reduced lock contention compared to a single lock for the entire map.
- **Generics Support**: Built using Go's generics, allowing the storage of any type as key and value.

## Installation

To use CSMap in your project, simply copy the CSMap code into your Go module.

```go
import "github.com/dm1trypon/csmap"
```

## Usage

### Creating a New CSMap

To create a new instance of CSMap, specify the number of shards you want to use:

```go
csMap := csmap.NewCSMap[int, string](16) // Create a new CSMap with 16 shards
```

### Setting a Value

You can add or update a value associated with a specific key using the `Set` method:

```go
csMap.Set(1, "value1") // Adds or updates the value associated with key 1
```

### Getting a Value

To retrieve a value, use the `Get` method, which returns the value and a boolean indicating if the key exists:

```go
value, exists := csMap.Get(1)
if exists {
    fmt.Println("Value:", value)
}
```

### Deleting a Value

If you want to remove a key and its associated value from the map, use the `Delete` method:

```go
csMap.Delete(1) // Removes the key 1 from the map
```

## Implementation Details

### Structure

The `CSMap` consists of an array of `Shard` structures, where each `Shard` is itself a map protected by a read-write mutex. The design ensures that each shard can be accessed concurrently without locks from other shards.

### Hashing

The `hash` function computes a hash for each key to determine its corresponding shard. This implementation uses Go's unsafe package for pointer manipulation to create a hash.

### Thread Safety

- The `Set` and `Delete` methods lock the specific shard's mutex for writing.
- The `Get` method uses a read lock to safely access the values.

## Example

Here's a complete example demonstrating how to use CSMap:

```go
package main

import (
    "fmt"
	
    "github.com/dm1trypon/csmap"
)

func main() {
    csMap := csmap.NewCSMap[int, string](16)

    csMap.Set(1, "hello")
    csMap.Set(2, "world")

    if value, exists := csMap.Get(1); exists {
        fmt.Println("Key 1:", value)
    }

    csMap.Delete(1)
}
```
