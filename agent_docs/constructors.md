# Sharing Configuration Across Multiple Functions

Use the `func New()` constructor pattern when multiple functions need the same parameters.

## Requirements

- **Public fields only** - capitalize field names (e.g., `ApiUrl`, not `apiUrl`)
- **One constructor per module**
- **Return module pointer** - return `*YourModule`

## Example

```go
func New(
    // +default="https://api.example.com"
    apiUrl string,
    // +optional
    apiKey string,
) *MyModule {
    return &MyModule{
        ApiUrl: apiUrl,
        ApiKey: apiKey,
    }
}

type MyModule struct {
    ApiUrl string
    ApiKey string
}

func (m *MyModule) GetUser(id string) (*User, error) {
    // m.ApiUrl and m.ApiKey available here
}

func (m *MyModule) CreateUser(name string) (*User, error) {
    // Same fields available here
}
```
