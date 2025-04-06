# Upgrading from insights v1 to v2

If a project has an existing version 1 file, ex: https://raw.githubusercontent.com/in-toto/go-witness/refs/heads/main/SECURITY-INSIGHTS.yml that file can be used as the basis to create a version 2 file.

## Automated conversion process

### Dependencies

The script depends on `curl` and `cue`

### Generate the new `security-insights.yml`

Using the example insights file from `go-witness`

```sh
./upgrade-to-v2.sh \
    "https://raw.githubusercontent.com/in-toto/go-witness/refs/heads/main/SECURITY-INSIGHTS.yml" \
    "go-witness"
```

You'll see output like:

```
Downloading v1 schema
Downloading v1 insights
Validating v1 insights against the v1 schema
"security-testing".0."tool-version": conflicting values 2 and string (mismatched types int and string):
    ./insights-v1.yml:49:17
    ./v1-schema.cue:15:2
    ./v1-schema.cue:150:25
    ./v1-schema.cue:174:21
Error: v1 insights file failed schema validation. Please correct the errors in v1 data before attempting another conversion
```

Edit `insights-v1.yml` to correctly specify a value of `"2"` for all `security-testing.*.tool-version` values and then re-run the script.

```sh
./upgrade-to-v2.sh \
    "https://raw.githubusercontent.com/in-toto/go-witness/refs/heads/main/SECURITY-INSIGHTS.yml" \
    "go-witness"
```

You'll see output like:

```
Found local v1-schema.cue, skipping download
Found local insights-v1.yml, skipping download
Validating v1 insights against the v1 schema
Found local schema.cue, skipping v2 schema download
Converting v1 data to v2 insights and saving to security-insights.yml
Validating v2 insights against the v2 schema


Thank you for using the v2 upgrade script.

The v1 insights data in https://raw.githubusercontent.com/in-toto/go-witness/refs/heads/main/SECURITY-INSIGHTS.yml has been converted to the v2 schema and saved to security-insights.yml. You must review the file and make any necessary adjustments to resolve TODOs before using it.
```

### Validate conversion and address TODOs

The conversion process will use default values as needed if it cannot determine the v1 source or there is no mapping from a v1 field to v2 field. You MUST review tne v2 output and replace all occurences of `TODO` with the relevant data for your project.