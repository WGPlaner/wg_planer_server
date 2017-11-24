# Contribution Guidelines

## Introduction

This document explains how to contribute changes to the WGPlaner (Server) project.

## Testing

Before sending code out for review, run all the tests for the
whole tree to make sure the changes don't break other usage.

## Styleguide

For imports you should use the following format (_without_ the comments)
```go
import (
  // stdlib
  "encoding/json"
  "fmt"

  // local packages
  "github.com/wgplaner/wg_planer_server/models"
  "github.com/wgplaner/wg_planer_server/wgplaner"

  // external packages
  "github.com/foo/bar"
)
```
