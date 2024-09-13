# cargobump

Rust tool to declaratively bump dependencies using cargo.

# Usage

The idea is that there are some `packages` that should be applied to the upstream
Cargo.lock file. You can specify these via `--packages` flag, or via
`--bump-file`.

## Specifying Dependencies to be patched

You can specify the patches that should be applied two ways. They are mutually
exclusive, so you can only specify one of them at the time.

### --packages flag

You can specify patches via `--packages` flag by encoding them
(similarly to gobump) in the following format:

```shell
--packages="<name@version[@scope[@type]]> <name...>"
```



### --bump-file flag

You can specify a yaml file that contains the patches, which is the preferred
way, because it's less error prone, and allows for inline comments to keep track
of which patches are for which CVEs.

An example yaml file looks like this:
```yaml
patches:
  # CVE-2023-34062
  - name: tokio
    version: 1.0.39
  # CVE-2023-5072
  - name: chrono
    version: "20231013"
```

