# envlayer

Hierarchical environment variable manager that merges `.env` files based on runtime context.

## Installation

```bash
go get github.com/yourname/envlayer
```

## Usage

`envlayer` loads and merges `.env` files in order of priority — base → environment-specific → local overrides.

**File resolution order** (later files take precedence):
```
.env → .env.<environment> → .env.local → .env.<environment>.local
```

**Example:**

```go
package main

import (
    "fmt"
    "github.com/yourname/envlayer"
)

func main() {
    err := envlayer.Load(envlayer.Options{
        Environment: "production", // reads .env, .env.production, .env.local, etc.
        AutoExpand:  true,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(os.Getenv("DATABASE_URL"))
}
```

**CLI usage:**

```bash
envlayer run --env staging -- ./myapp
```

## Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Environment` | `string` | `""` | Runtime environment name |
| `AutoExpand` | `bool` | `false` | Expand variable references |
| `Override` | `bool` | `true` | Override existing env vars |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)