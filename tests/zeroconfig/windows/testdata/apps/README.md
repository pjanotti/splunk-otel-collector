# Tests for Splunk Zero Configuration of .NET applications hosted on Windows IIS

Tests in this directory validate that automatic instrumentation of
.NET and .NET Framework applications hosted on Windows IIS are working under
a deployment following the Zero Configuration procedures.

## Requirements to run the tests

- Build artifacts:
  - Collector MSI (download it and expand its contents to `TBD`)
- Windows OS:
  - .NET Framework
  - .NET SDK
  - NuGet command-line
  - Docker configured to run Windows containers
