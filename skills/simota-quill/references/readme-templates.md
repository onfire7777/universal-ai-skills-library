# README Scaffolding Templates

Templates for different project types to ensure consistent documentation.

Purpose: Read this when Quill must create or repair README structure for a library, application, or CLI project.

Contents:
- `Library/Package README`: package-oriented installation, API, and config template
- `Application README`: app onboarding, project structure, scripts, and deployment template
- `CLI Tool README`: command-focused installation, usage, config, and examples template

## Library/Package README

```markdown
# Package Name

Brief description of what this package does.

## Installation

\`\`\`bash
npm install package-name
# or
yarn add package-name
\`\`\`

## Quick Start

\`\`\`typescript
import { mainFunction } from 'package-name';

const result = mainFunction({ option: 'value' });
\`\`\`

## API Reference

### `mainFunction(options)`

Description of the main function.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `option` | `string` | - | Required option |
| `timeout` | `number` | `5000` | Optional timeout in ms |

**Returns**: `ResultType` - Description of return value

**Example**:
\`\`\`typescript
const result = mainFunction({ option: 'value', timeout: 10000 });
\`\`\`

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `PACKAGE_API_KEY` | API key for service | - |
| `PACKAGE_TIMEOUT` | Request timeout | `5000` |

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup.

## License

MIT
```

## Application README

```markdown
# Application Name

Brief description of the application.

## Prerequisites

- Node.js >= 18
- PostgreSQL >= 14
- Redis >= 6

## Getting Started

### 1. Clone and Install

\`\`\`bash
git clone https://github.com/org/repo.git
cd repo
npm install
\`\`\`

### 2. Environment Setup

\`\`\`bash
cp .env.example .env
# Edit .env with your values
\`\`\`

### 3. Database Setup

\`\`\`bash
npm run db:migrate
npm run db:seed  # Optional: seed test data
\`\`\`

### 4. Run Development Server

\`\`\`bash
npm run dev
# Open http://localhost:3000
\`\`\`

## Project Structure

\`\`\`
src/
├── api/          # API routes
├── components/   # React components
├── lib/          # Shared utilities
├── pages/        # Page components
└── types/        # TypeScript types
\`\`\`

## Available Scripts

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run test` | Run tests |
| `npm run lint` | Run linter |

## Deployment

See [docs/deployment.md](./docs/deployment.md) for deployment instructions.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

MIT
```

## CLI Tool README

```markdown
# CLI Tool Name

Brief description of the CLI tool.

## Installation

\`\`\`bash
npm install -g cli-tool-name
# or
npx cli-tool-name
\`\`\`

## Usage

\`\`\`bash
cli-tool <command> [options]
\`\`\`

## Commands

### `init`

Initialize a new project.

\`\`\`bash
cli-tool init [project-name]

Options:
  --template <name>  Use a specific template
  --force            Overwrite existing files
\`\`\`

### `build`

Build the project.

\`\`\`bash
cli-tool build [options]

Options:
  --watch    Watch for changes
  --minify   Minify output
\`\`\`

## Configuration

Create `cli-tool.config.js` in your project root:

\`\`\`javascript
module.exports = {
  input: './src',
  output: './dist',
  plugins: [],
};
\`\`\`

## Examples

### Basic Usage

\`\`\`bash
cli-tool init my-project
cd my-project
cli-tool build
\`\`\`

### With Options

\`\`\`bash
cli-tool build --watch --minify
\`\`\`

## License

MIT
```
