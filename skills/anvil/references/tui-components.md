# TUI Components

**Purpose:** Practical TUI component patterns for prompts, tables, spinners, and progress views across Node.js, Python, Go, and Rust.
**Read when:** Adding terminal interactivity, rich output, or a full-screen TUI to a CLI.

## Contents

- Language/Library Matrix
- Progress Indicators
- Selection Menus
- Table Display
- Rust CLI Full Pattern

## Language/Library Matrix

| Language | Interactive Prompts | Rich Output | Full TUI |
|----------|---------------------|-------------|----------|
| **Node.js** | inquirer, prompts | chalk, ora, cli-table3 | ink, blessed |
| **Python** | click, questionary | rich, colorama | textual, urwid |
| **Go** | survey, promptui | color, tablewriter | bubbletea, tview |
| **Rust** | dialoguer, inquire | colored, prettytable | ratatui, crossterm |

---

## Progress Indicators

### Node.js (ora)

```typescript
import ora from 'ora';

async function withSpinner<T>(task: () => Promise<T>, message: string): Promise<T> {
  const spinner = ora(message).start();
  try {
    const result = await task();
    spinner.succeed();
    return result;
  } catch (error) {
    spinner.fail();
    throw error;
  }
}
```

### Python (rich)

```python
from rich.progress import Progress, SpinnerColumn, TextColumn

def with_progress(tasks: list[tuple[str, Callable]]):
    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
    ) as progress:
        for description, task in tasks:
            task_id = progress.add_task(description)
            task()
            progress.update(task_id, completed=True)
```

### Go (bubbletea)

```go
type spinnerModel struct {
    spinner spinner.Model
    message string
    done    bool
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    return m, nil
}
```

### Rust (indicatif)

```rust
use indicatif::{ProgressBar, ProgressStyle};
use std::time::Duration;

fn with_spinner<T, F>(message: &str, task: F) -> T
where F: FnOnce() -> T {
    let spinner = ProgressBar::new_spinner();
    spinner.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.cyan} {msg}").unwrap()
    );
    spinner.set_message(message.to_string());
    spinner.enable_steady_tick(Duration::from_millis(100));
    let result = task();
    spinner.finish_with_message(format!("✓ {}", message));
    result
}
```

---

## Selection Menus

### Node.js (inquirer)

```typescript
import inquirer from 'inquirer';

async function selectOption<T extends string>(
  message: string,
  choices: { name: string; value: T }[]
): Promise<T> {
  const { selection } = await inquirer.prompt([
    { type: 'list', name: 'selection', message, choices },
  ]);
  return selection;
}
```

### Python (questionary)

```python
import questionary

def select_option(message: str, choices: list[str]) -> str:
    return questionary.select(message, choices=choices, use_shortcuts=True).ask()
```

### Rust (dialoguer)

```rust
use dialoguer::{theme::ColorfulTheme, Select, Input, Confirm};

fn interactive_setup() -> Result<Config, Box<dyn std::error::Error>> {
    let name: String = Input::with_theme(&ColorfulTheme::default())
        .with_prompt("Project name")
        .default("my-project".into())
        .interact_text()?;

    let template = Select::with_theme(&ColorfulTheme::default())
        .with_prompt("Select template")
        .items(&["minimal", "full", "custom"])
        .default(0)
        .interact()?;

    Ok(Config { name, template })
}
```

---

## Table Display

### Node.js (cli-table3)

```typescript
import Table from 'cli-table3';

function displayTable(headers: string[], rows: string[][]): void {
  const table = new Table({ head: headers, style: { head: ['cyan'] } });
  rows.forEach(row => table.push(row));
  console.log(table.toString());
}
```

### Python (rich)

```python
from rich.console import Console
from rich.table import Table

def display_table(title: str, columns: list[str], rows: list[list[str]]):
    console = Console()
    table = Table(title=title)
    for col in columns:
        table.add_column(col)
    for row in rows:
        table.add_row(*row)
    console.print(table)
```

### Rust (tabled)

```rust
use tabled::{Table, Tabled};

#[derive(Tabled)]
struct Row {
    name: String,
    status: String,
    count: u32,
}

fn display_table(rows: Vec<Row>) {
    let table = Table::new(rows).to_string();
    println!("{}", table);
}
```

---

## Rust CLI Full Pattern (Clap)

```rust
use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "myapp", version, about)]
struct Cli {
    #[arg(short, long, action = clap::ArgAction::Count)]
    verbose: u8,

    #[arg(long)]
    json: bool,

    #[arg(long)]
    no_color: bool,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize a new project
    Init {
        #[arg(short, long)]
        name: Option<String>,
    },
    /// Build the project
    Build {
        #[arg(long)]
        watch: bool,
    },
}

fn main() {
    let cli = Cli::parse();
    match cli.command {
        Commands::Init { name } => init_project(name),
        Commands::Build { watch } => build_project(watch),
    }
}
```
