# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## [1.3.0] - 2024-09-04

### Changed
- New `NotLastByte` padding method, that is a faster version of arbitrary tail byte padding.
- Renamed package `tests` to the more appropriate `benchmarks`.

## [1.2.0] - 2024-04-28

### Changed
- New `PadLastBlock` function that is much more efficient than Pad.

## [1.1.0] - 2024-04-24

### Changed
- Padding and unpadding is now constant-time (at least, as much, as possible).

## [1.0.2] - 2024-04-20

### Changed
- Moved example in package comment to example test file.

## [1.0.1] - 2024-04-20

### Fixed
- Fixed misplaced documentation comment of `Pad` function.

## [1.0.0] - 2024-04-20

### Added
- Initial release.
