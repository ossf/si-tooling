# Command-line Generator and Validator

> [!WARNING]
> This was designed as part of the v1 release, and does not have any feature support for v2.
>
> The recommended method for validating v2 is to use the Cue schema included in the release.

A Python command-line tool to help maintainters, developers, and contributors to generate or validate the SECURITY INSIGHTS yaml file.

## Usage

[Docker](https://www.docker.com/) needs to be installed. Build the container image:

```
docker build -t sec-insights .
```

If you want to validate a `SECURITY-INSIGHTS.yml`, run the following command in the folder of the project's `SECURITY-INSIGHTS.yml`:

```
docker run -v $PWD:/tmp -it sec-insights verify SECURITY-INSIGHTS.yml
```

If you want to create a new `SECURITY-INSIGHTS.yml` by complying the YAML schema, launch this command:

```
docker run -v $PWD:/tmp -it sec-insights create
```

and fill out the questionnaire by following the wizard. The questions labeled with `(optional)` are not mandatory and can be skipped.

## Bugs

If you find any bugs, please open an issue or submit a pull request.

### Known Bugs

- [ ] Value type and format are not printed in the wizard
- [ ] The script supports just single-line comments

## Security

If you find a security vulnerability, please report it via [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability).



