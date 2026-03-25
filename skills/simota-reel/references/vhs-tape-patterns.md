# VHS Tape Patterns Reference

Comprehensive reference for VHS (charmbracelet/vhs) — a terminal recording tool that uses `.tape` files to script and produce GIF/MP4/WebM recordings of terminal sessions.

---

## VHS Overview

VHS converts `.tape` script files into terminal recordings. Each `.tape` file is a declarative script describing keystrokes, settings, and timing that VHS replays in a virtual terminal, capturing the output as video or image.

```bash
# Record a tape file
vhs record > demo.tape    # Interactive recording to tape
vhs demo.tape             # Play back and render output
vhs < demo.tape           # Pipe tape content directly
```

### Tape File Structure

A `.tape` file follows this general order:

1. **Settings** — Output format, shell, dimensions, theme, fonts
2. **Requirements** — Commands that must exist on PATH
3. **Environment** — Environment variables for the session
4. **Actions** — Keystrokes, sleeps, waits, screenshots

Lines starting with `#` are comments. Blank lines are ignored.

---

## Commands Reference

### Type — Enter Text

Type text into the terminal. Supports quoted strings with escape sequences.

```tape
Type "echo 'Hello, World!'"
Type "ls -la"
Type "git status"
```

**Typing speed override** — Append `@<duration>` to control per-character speed:

```tape
Type@50ms "fast typing"
Type@500ms "slow dramatic typing"
Type@0s "instant typing"
```

### Enter — Press Enter Key

Send the Enter/Return key to execute commands.

```tape
Type "echo hello"
Enter
```

**Repeat shorthand** — Specify a count:

```tape
Enter 3    # Press Enter 3 times
```

### Sleep — Pause Execution

Pause for a specified duration. Useful for pacing and letting output render.

```tape
Sleep 1        # Sleep 1 second
Sleep 500ms    # Sleep 500 milliseconds
Sleep 2s       # Sleep 2 seconds
Sleep 0.5      # Sleep 0.5 seconds
```

### Ctrl+Key — Control Characters

Send control key combinations.

```tape
Ctrl+C         # Interrupt / SIGINT
Ctrl+D         # EOF / exit
Ctrl+L         # Clear screen
Ctrl+A         # Move to beginning of line
Ctrl+E         # Move to end of line
Ctrl+R         # Reverse search (bash/zsh)
Ctrl+U         # Clear line before cursor
Ctrl+K         # Clear line after cursor
Ctrl+W         # Delete word before cursor
Ctrl+Z         # Suspend process
```

### Alt+Key — Alt/Meta Characters

Send Alt (Meta) key combinations.

```tape
Alt+B          # Move back one word
Alt+F          # Move forward one word
Alt+D          # Delete word after cursor
Alt+Backspace  # Delete word before cursor
```

### Backspace — Delete Characters

Press the Backspace key to delete characters.

```tape
Backspace      # Press once
Backspace 5    # Press 5 times
```

### Tab — Tab Key / Autocomplete

Press the Tab key, often used for shell autocompletion.

```tape
Type "git ch"
Tab            # Trigger autocomplete
Tab 2          # Press Tab twice (cycle completions)
```

### Space — Space Key

Press the Space key.

```tape
Space          # Press once
Space 3        # Press 3 times
```

### Arrow Keys — Navigation

Navigate with arrow keys. All support repeat counts.

```tape
Up             # Previous command in history
Up 3           # Go up 3 times
Down           # Next command in history
Down 2         # Go down 2 times
Left           # Move cursor left
Left 5         # Move cursor left 5 characters
Right          # Move cursor right
Right 3        # Move cursor right 3 characters
```

### Page Navigation

```tape
PageUp         # Scroll up one page
PageDown       # Scroll down one page
```

### Home / End

```tape
Home           # Move to beginning of line
End            # Move to end of line
```

### Escape

```tape
Escape         # Press Escape key
```

### Hide / Show — Control Recording Visibility

Hide terminal output from the recording. Useful for setup commands that should not appear in the final output.

```tape
Hide
Type "export SECRET=hidden_setup"
Enter
Sleep 1
Show

# From here, everything is visible in the recording
Type "echo 'Now recording!'"
Enter
```

### Wait — Wait for Output

Wait until specific text appears on screen before proceeding. Critical for commands with unpredictable execution times.

```tape
# Wait for a specific string on screen
Type "npm install"
Enter
Wait+Screen "added"        # Wait until "added" appears
Sleep 1

# Wait with regex pattern
Wait+Screen /\$\s*/         # Wait for shell prompt
```

### Source — Include Another Tape File

Include and execute commands from another `.tape` file. Enables modular, reusable tape compositions.

```tape
Source "setup.tape"          # Run setup commands
Source "common/header.tape"  # Relative path inclusion
```

### Require — Assert Command Availability

Ensure a command exists on PATH before proceeding. VHS exits with an error if the command is not found.

```tape
Require echo
Require git
Require node
Require docker
```

### Env — Set Environment Variables

Set environment variables for the recording session.

```tape
Env EDITOR "vim"
Env TERM "xterm-256color"
Env PS1 "$ "
Env NO_COLOR "1"
Env PATH "/usr/local/bin:/usr/bin:/bin"
```

### Screenshot — Capture Still Image

Take a screenshot at the current point in the recording.

```tape
Type "neofetch"
Enter
Sleep 2
Screenshot "neofetch.png"
```

---

## Settings Reference

Settings configure the recording environment. Place them at the top of the `.tape` file before any commands.

### Output — Recording Format

Specify the output file. The extension determines format.

```tape
Output demo.gif             # Animated GIF (default)
Output demo.mp4             # MP4 video
Output demo.webm            # WebM video
Output demo.png             # PNG screenshot (last frame)
Output frames/              # Directory of PNG frames
```

Multiple outputs can be specified:

```tape
Output demo.gif
Output demo.mp4
```

### Shell — Terminal Shell

```tape
Set Shell "bash"            # Default
Set Shell "zsh"
Set Shell "fish"
Set Shell "sh"
Set Shell "powershell"
```

### Font Settings

```tape
Set FontSize 14             # Default: 22
Set FontFamily "JetBrains Mono"  # Default: monospace system font
Set LetterSpacing 0         # Default: 0 (pixels)
Set LineHeight 1.2          # Default: 1.0
```

### Dimensions

Width and Height are specified in character columns and rows respectively.

```tape
Set Width 80                # Default: 80 (characters)
Set Height 24               # Default: 24 (rows)
```

### Theme

Use a built-in theme or reference a custom theme file.

```tape
Set Theme "Monokai"
Set Theme "Dracula"
Set Theme "Catppuccin Mocha"
Set Theme "GitHub Dark"
Set Theme "Solarized Dark"
Set Theme "Nord"
Set Theme "One Dark"
Set Theme "Rosé Pine"
Set Theme "Tokyo Night"
```

Custom theme from JSON:

```tape
Set Theme "path/to/theme.json"
```

### Typing Speed

Control the default speed of `Type` commands.

```tape
Set TypingSpeed 50ms        # Default: 50ms per character
Set TypingSpeed 0           # Instant typing
Set TypingSpeed 100ms       # Slow, deliberate typing
```

### Padding and Margin

```tape
Set Padding 0               # Default: 0 (pixels around terminal content)
Set Margin 0                # Default: 0 (pixels outside window frame)
Set MarginFill "#674EFF"    # Background color behind margin
Set MarginFill "path/to/image.png"  # Background image behind margin
```

### Framerate and Playback

```tape
Set Framerate 50            # Default: 50 (frames per second)
Set PlaybackSpeed 1         # Default: 1.0 (multiplier; 2 = double speed)
Set LoopOffset 0            # Default: 0 (percentage; controls GIF loop start)
```

### Cursor

```tape
Set CursorBlink true        # Default: true
```

### Window Chrome

```tape
Set WindowBar Colorful      # Options: Colorful, ColorfulRight, Rings, RingsRight
Set WindowBarSize 40        # Default: 40 (pixels)
Set BorderRadius 8          # Default: 0 (pixels, rounded corners)
```

---

## Scene Composition Patterns

Structure recordings into logical scenes for clarity and storytelling.

### Pattern 1: Opening Scene — Context Setup

Establish context by showing the environment, directory, or state before the main action. Use `Hide`/`Show` for invisible pre-configuration.

```tape
# --- Pre-configuration (hidden) ---
Hide
Type "cd /tmp/demo-project"
Enter
Type "git checkout main && git pull"
Enter
Sleep 2
Show

# --- Opening: show context ---
Type "pwd"
Enter
Sleep 1
Type "ls -la"
Enter
// ...
```

### Pattern 2: Action Scene — Main Demonstration

The core of the recording. Pace keystrokes and pauses to keep the viewer engaged.

```tape
# --- Action: demonstrate the feature ---
Sleep 1
Type "git checkout -b feature/new-widget"
Enter
Sleep 1

Type "cat src/widget.ts"
Enter
Sleep 3

Type "npm test"
Enter
Wait+Screen "Tests:"
Sleep 2
```

### Pattern 3: Result Scene — Show Outcome

End with the result or confirmation of the action. Allow extra sleep time so the viewer can absorb the output.

```tape
# --- Result: confirm success ---
Type "git log --oneline -5"
Enter
Sleep 3

Type "echo 'Done!'"
Enter
Sleep 5    # Hold the final frame longer for GIF loops
```

### Pattern 4: Full Scene Template

A complete tape combining all three scenes.

```tape
# Settings
Output demo.gif
Set Shell "bash"
Set FontSize 16
Set Width 100
Set Height 30
Set Theme "Dracula"
Set TypingSpeed 75ms
Set Padding 10
Set WindowBar Colorful
Set BorderRadius 8

# Requirements
Require git
Require node
// ...
```

### Transition Patterns Between Scenes

#### Clear Screen Transition

```tape
# End of scene 1
Sleep 2
Ctrl+L         # Clear screen
Sleep 0.5
# Begin scene 2
```

#### Comment Banner Transition

```tape
Sleep 1
Type "# --- Step 2: Deploy ---"
Enter
Sleep 1
```

#### Prompt Reset Transition

```tape
Sleep 2
Enter           # Extra blank line
Enter
Sleep 0.5
```

### Modular Scene Composition with Source

Break large recordings into reusable scene files.

```tape
# main.tape
Output full-demo.gif
Set Theme "Catppuccin Mocha"
Set Width 100
Set Height 30

Source "scenes/setup.tape"
Source "scenes/build.tape"
Source "scenes/test.tape"
Source "scenes/deploy.tape"

Sleep 5
```

```tape
# scenes/setup.tape
Hide
Type "cd ~/project"
Enter
Sleep 1
Show
Type "git status"
Enter
Sleep 2
```

### Timing Guidelines

| Context                 | Recommended Sleep |
|-------------------------|-------------------|
| After showing a command | 0.5s - 1s        |
| After command output    | 1s - 3s          |
| After complex output    | 3s - 5s          |
| Final frame hold        | 3s - 5s          |
| Scene transition        | 0.5s - 1s        |
| Typing speed (default)  | 50ms - 100ms     |
| Typing speed (emphasis) | 100ms - 200ms    |
| Typing speed (setup)    | 0ms - 30ms       |
