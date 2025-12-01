# Installing a Module as a Dependency

To install a module into another module:

```bash
dagger install --mod <target-module> <source-module>
```

## Installing into the Repo Module

```bash
dagger install --mod ./.dagger ./terraform-docs
```

## Installing External Dependencies

```bash
dagger install --mod ./my-module github.com/example/some-dependency
```
