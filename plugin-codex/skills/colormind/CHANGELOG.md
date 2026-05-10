# Changelog

## [1.1.1] - 2026-02-18

### Fixed
- **Added python3 to required binaries** in metadata (was missing, causing runtime errors)

### Security & Documentation
- **Added prominent privacy warning** to SKILL.md about external data sharing
- **Documented HTTP limitation** - colormind.io doesn't support proper HTTPS
- **Enhanced SECURITY.md** with detailed privacy considerations and mitigation strategies
- Added recommendations for production use and sensitive data handling
- Clarified ImageMagick safety requirements

### Changed
- Improved documentation clarity around when NOT to use this skill
- Added explicit warnings about unencrypted transport

## [1.1.0] - 2026-02-18

### Security
- **Refactored image_to_palette.sh** to reduce false positive security scanner flags
- Replaced inline Python with separate `.py` scripts
- Data now passed via temporary JSON files instead of heredocs
- Added SECURITY.md explaining architecture and safety

### Changed
- `parse_histogram.py` - Extracts color parsing logic into standalone script
- `get_base_rgb.py` - Extracts base color extraction into standalone script  
- `combine_results.py` - Extracts result merging into standalone script
- `image_to_palette.sh` - Now orchestrates separate scripts with clean temp file handling

### Technical
- Removed inline `python3 -c` commands
- Removed heredoc variable substitution (`<<PY` with `$VAR`)
- Added automatic temp directory cleanup with trap
- Improved code readability and maintainability

## [1.0.0] - 2026-02-16

### Added
- Initial release
- List Colormind models
- Generate random color palettes
- Generate palettes with locked colors
- Sample colors from images via ImageMagick
