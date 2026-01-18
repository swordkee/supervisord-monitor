# Agent Guidelines for Supervisord Monitor

## Project Overview
Go + Vue application for monitoring multiple Supervisord servers. Backend uses Gin framework, frontend uses Vue 3 with Vite. Static assets are embedded into a single executable.

## Build Commands

### Full Build
```bash
# All platforms
make build              # Build complete app (frontend + backend)
make frontend           # Build frontend only
make backend            # Build backend only

# Platform-specific
make windows            # Build Windows executable (.exe)
make linux              # Build Linux executable
make darwin             # Build macOS executable
```

### Development
```bash
# Start both frontend (port 5173) and backend (port 8080)
make dev

# Frontend only
cd frontend && npm run dev

# Backend only
make run                # Equivalent to: go run main.go
```

### Install Dependencies
```bash
make install            # Installs Go deps and frontend node_modules
```

### Clean
```bash
make clean              # Removes build artifacts, node_modules, and dist
```

## Linting & Testing

**Note:** This project currently does not have automated tests or linting configured.

When adding tests:
- Go tests: Place alongside source files with `_test.go` suffix
- Run single test: `go test -v -run TestName ./path/to/package`
- Run all tests: `go test ./...`

Consider adding:
- `gofmt` for Go code formatting
- ESLint/Prettier for frontend code
- Test coverage reporting

## Code Style Guidelines

### Go (Backend)

#### Imports
- Use blank lines between groups: stdlib, third-party, local packages
- Keep imports sorted alphabetically within groups
- Example:
```go
import (
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"

    "supervisord-monitor/config"
    "supervisord-monitor/services"
)
```

#### Naming Conventions
- **Packages**: lowercase, single word (e.g., `handlers`, `services`, `config`)
- **Exported functions**: PascalCase (e.g., `GetDashboard`, `SetupRoutes`)
- **Unexported functions**: camelCase (e.g., `getProcessName`, `formatProcessDescription`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE for cross-package use
- **Variables**: camelCase for locals, PascalCase for exported fields

#### Error Handling
- Always check errors, never ignore them
- Return errors from functions, don't panic in application code
- Use `fmt.Errorf` with `%w` for error wrapping
- Example:
```go
client, err := services.NewSupervisorClient(server)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

#### Struct Tags
- Use `mapstructure` tags for YAML config structs
- Use `json` tags for API response structs
- Keep tags on same line as field for readability

#### HTTP Handlers
- Accept `*gin.Context` as parameter
- Return early on errors with appropriate status codes
- Use `c.JSON()` for JSON responses with status codes
- Example statuses: 400 (bad request), 500 (internal error), 200 (success)

#### Configuration
- Use Viper for config management
- Provide sensible defaults via `viper.SetDefault()`
- Support multiple config file locations
- Use `mapstructure` tags for YAML mapping
- HTTP Auth: Use `crypto/subtle.ConstantTimeCompare` for secure password comparison to prevent timing attacks

### Vue (Frontend)

#### Component Structure
- Use Vue 3 Composition API with `<script setup>`
- Define props with types
- Use `emit()` for parent communication
- Example:
```vue
<script setup>
import { ref, computed } from 'vue'
import axios from 'axios'

const props = defineProps({
  server: { type: Object, required: true },
  showHost: { type: Boolean, default: false }
})

const emit = defineEmits(['refresh'])
</script>
```

#### Naming Conventions
- **Components**: PascalCase (e.g., `ServerCard.vue`)
- **Files**: kebab-case for utilities, PascalCase for components
- **Variables**: camelCase
- **Constants**: UPPER_SNAKE_CASE or PascalCase

#### API Calls
- Use Axios for HTTP requests
- Wrap in async/await with try/catch
- Handle errors gracefully (alerts or UI feedback)
- Use absolute paths for API endpoints (e.g., `/api/dashboard`)

#### Styling
- Use Bootstrap 2.3 classes from CDN
- Scoped styles in `<style scoped>` blocks
- Keep component-specific styles within components
- Global styles in `src/style.css`

#### State Management
- Use `ref()` for reactive primitives
- Use `computed()` for derived state
- Use `onMounted()` for side effects
- Use `onUnmounted()` for cleanup

#### Template Best Practices
- Use `v-if` for conditional rendering
- Use `v-for` with `:key` attribute
- Prefer `@click.prevent` over `href="#"` for actions
- Use Bootstrap button classes: `btn btn-mini btn-success`, etc.

## File Organization

### Backend Structure
```
supervisord-monitor-go/
├── config/          # Config loading and structures
├── handlers/        # HTTP request handlers (controllers)
├── services/        # Business logic and external service clients
├── main.go          # Application entry point
└── embed.go         # Embedded frontend assets
```

### Frontend Structure
```
frontend/
├── src/
│   ├── components/  # Reusable Vue components
│   ├── App.vue      # Root component
│   ├── main.js      # Entry point
│   └── style.css    # Global styles
├── index.html       # HTML template
└── vite.config.js   # Build configuration
```

## Adding New Features

### Backend API Endpoint
1. Add handler function in `handlers/handlers.go`
2. Register route in `SetupRoutes()` function
3. Add service method in `services/supervisord.go` if needed
4. Return JSON with appropriate HTTP status codes

### Frontend Component
1. Create `.vue` file in `frontend/src/components/`
2. Import and use in parent component
3. Add props for data passing
4. Emit events for parent communication
5. Use Bootstrap classes for consistent styling

### Configuration Options
1. Add field to `Config` struct in `config/config.go`
2. Add default value via `viper.SetDefault()`
3. Add to example `config.yaml`
4. Update README.md documentation

## Important Notes

- Frontend is embedded in Go binary - rebuild backend after frontend changes
- Use `log.Printf()` for debugging, remove before committing
- Process naming: `group:name` format when group differs from name
- Server authentication: check `Username != ""` for auth presence
- Always handle errors from supervisord client calls
- Maintain backward compatibility with original PHP version API
