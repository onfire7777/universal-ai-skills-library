# VHS Tape Templates Reference

Reusable .tape templates for common terminal recording scenarios.
Each template is copy-paste-ready with inline comments explaining each section.

---

## 1. Quickstart Template

A concise "install and run" demo. Ideal for README hero images and getting-started guides.
Target duration: under 15 seconds.

```tape
# Quickstart Demo
# Shows: package installation → first command execution
# Duration: ~12 seconds

Output quickstart.gif

# --- Terminal settings ---
Set FontSize 16
Set Width 80
Set Height 20
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 60ms

# --- Recording ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `mycli` | Your CLI tool name |
| `npm install -g mycli` | Your install command (brew, pip, cargo, etc.) |
| `mycli hello --name world` | Your flagship command |
| `quickstart.gif` | Desired output filename (.gif, .mp4, .webm) |

---

## 2. Feature Demo Template

Showcase a specific CLI feature with flags and options.
Shows help text first, then demonstrates the feature with realistic data.

```tape
# Feature Demo
# Shows: help output → feature execution with flags → result
# Duration: ~20 seconds

Output feature-demo.gif

# --- Terminal settings ---
Set FontSize 14
Set Width 100
Set Height 28
Set Theme "Catppuccin Mocha"
Set Padding 16
Set TypingSpeed 50ms

# --- Recording ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `mycli convert` | Your subcommand |
| `--format json --pretty` | Your feature flags |
| `input.csv`, `result.json` | Realistic filenames for your domain |
| `--stats` | Optional verification command |

---

## 3. Before-After Template

Two-phase recording showing improvement (performance, output quality, workflow).
Uses Hide/Show to silently set up state, then reveals the visible session.

```tape
# Before-After Comparison
# Shows: "before" state → transition → "after" state
# Duration: ~25 seconds

Output before-after.gif

# --- Terminal settings ---
Set FontSize 14
Set Width 100
Set Height 28
Set Theme "Catppuccin Mocha"
Set Padding 16
Set TypingSpeed 50ms

# --- Silent setup: prepare the "before" environment ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `config.old.yml` / `config.new.yml` | Your before/after artifacts |
| `mycli build --config` | Your build or processing command |
| `time` prefix | Use if showing performance improvement |
| `Hide` / `Show` blocks | Add more setup steps as needed |

---

## 4. Interactive Session Template

Demonstrate interactive prompts, menus, or wizards.
Shows a user navigating choices, entering values, and completing a flow.

```tape
# Interactive Session Demo
# Shows: wizard launch → user selections → completion
# Duration: ~30 seconds

Output interactive.gif

# --- Terminal settings ---
Set FontSize 14
Set Width 100
Set Height 30
Set Theme "Catppuccin Mocha"
Set Padding 16
Set TypingSpeed 70ms

# --- Recording ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `mycli init` | Your interactive command |
| `Down` / `Up` | Arrow key navigation for menus |
| `Tab` | Tab completion where applicable |
| Prompt responses | Match your CLI's actual prompts |
| `TypingSpeed 70ms` | Slower typing feels more natural for interactive |

### Key VHS Commands for Interactivity

```tape
# Navigation
Up                  # Arrow up
Down                # Arrow down
Left                # Arrow left
Right               # Arrow right
Tab                 # Tab key
Enter               # Enter/Return key
Space               # Space bar (toggle checkboxes)
Backspace           # Delete character

# Special keys
Escape              # Escape key
Ctrl+C              # Interrupt
Ctrl+D              # EOF / Exit
```

---

## 5. Error Handling Template

Show how a CLI gracefully handles errors: invalid input, missing files, network failures.
Demonstrates error messages and recovery paths.

```tape
# Error Handling Demo
# Shows: invalid input → missing file → recovery → success
# Duration: ~30 seconds

Output error-handling.gif

# --- Terminal settings ---
Set FontSize 14
Set Width 100
Set Height 28
Set Theme "Catppuccin Mocha"
Set Padding 16
Set TypingSpeed 50ms

# --- Recording ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `mycli process` | Your command that can demonstrate error paths |
| Error scenarios | Match your CLI's actual error messages |
| `--format xml` | An intentionally unsupported flag value |
| `proccess` | A realistic typo of your command |

### Tips for Error Demos

- Show the error message clearly, then pause so the viewer can read it
- Progress from simple errors to more complex ones
- Always end with a successful run to show recovery
- Use `Sleep` generously between scenarios for readability

---

## 6. Multi-Step Workflow Template

Complex workflow with multiple commands building on each other.
Shows a realistic terminal session where each step depends on the previous.

```tape
# Multi-Step Workflow Demo
# Shows: init → configure → execute → verify → deploy
# Duration: ~40 seconds

Output workflow.gif

# --- Terminal settings ---
Set FontSize 14
Set Width 110
Set Height 30
Set Theme "Catppuccin Mocha"
Set Padding 16
Set TypingSpeed 50ms

# --- Silent setup: prepare a clean workspace ---
// ...
```

### Customization Notes

| Placeholder | Replace with |
|---|---|
| `mycli init`, `config`, `add`, etc. | Your CLI's subcommands |
| `--template starter` | Your project templates |
| Resource names (`maindb`, `sessions`) | Realistic resource identifiers |
| `--env staging` | Your deployment targets |

---

## Common Settings Reference

### Theme Options

Popular themes that work well in recordings:

```tape
Set Theme "Catppuccin Mocha"     # Dark, warm, high contrast
Set Theme "Dracula"              # Dark, purple-accented
Set Theme "Monokai"              # Dark, classic developer theme
Set Theme "One Dark"             # Dark, Atom-inspired
Set Theme "Tokyo Night"          # Dark, blue-toned
Set Theme "Catppuccin Latte"     # Light, warm tones
Set Theme "GitHub Light"         # Light, clean
```

### Output Formats

```tape
Output demo.gif                  # Animated GIF (most compatible)
Output demo.mp4                  # MP4 video (smaller file size)
Output demo.webm                 # WebM video (web-optimized)
Output frames/                   # PNG frame sequence
```

### Window Chrome

```tape
Set WindowBar Colorful           # macOS-style colored dots
Set WindowBar Rings              # Ring-style buttons
Set WindowBar Blocks             # Block-style buttons
Set WindowBarSize 40             # Height of the title bar
```

### Cursor Styles

```tape
Set CursorBlink true             # Blinking cursor
Set CursorBlink false            # Static cursor
```

### Timing Presets

```tape
# Fast demo (experienced user feel)
Set TypingSpeed 30ms

# Normal demo (natural reading pace)
Set TypingSpeed 50ms

# Slow demo (tutorial, step-by-step)
Set TypingSpeed 80ms

# Pause durations
Sleep 0.5                        # Brief pause between keystrokes
Sleep 1                          # Short pause (transition)
Sleep 2                          # Medium pause (read short output)
Sleep 3                          # Long pause (read longer output)
Sleep 5                          # Extended pause (complex output)
```

### Size Presets

```tape
# Compact (badges, quick demos)
Set Width 80
Set Height 20
Set FontSize 16

# Standard (feature demos, tutorials)
Set Width 100
Set Height 28
Set FontSize 14

# Wide (log output, multi-column layouts)
Set Width 120
Set Height 32
Set FontSize 13

// ...
```
