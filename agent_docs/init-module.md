# Creating a New Module

From repository root:

```bash
dagger init --sdk=go --name=<module-name> <module-name>
```

Example:

```bash
dagger init --sdk=go --name=basics basics
```

## Import Path

Always use:

```go
import "dagger/<module-name>/internal/dagger"
```
