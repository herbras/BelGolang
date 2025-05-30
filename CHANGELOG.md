# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.6.1] - 2025-05-30

### Added
- Automated release workflow with GitHub Actions
- NPM package binaries for multiple platforms (darwin-x64, darwin-arm64, linux-x64, linux-arm64, win32-x64, win32-arm64)
- Dedicated README for NPM package distribution

### Changed
- Updated release process to use npm-package README for GitHub releases
- Improved binary distribution through automated builds

### Fixed
- Release workflow now properly uses npm-package specific documentation

## [Unreleased]
### Added
- Initial project setup
- Go modules and build configuration
- Cross-platform binary compilation 