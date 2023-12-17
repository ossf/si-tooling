# Command-line Generator and Validator

A Python command-line tool to help maintainters, developers, and contributors to generate or validate the SECURITY INSIGHTS yaml file.

## Usage

[Docker](https://www.docker.com/) needs to be installed. Build the container image:

```
docker build -t sec-insights .
```

If you want to validate a `SECURITY-INSIGHTS.yml`, run the following command by replacing `$/your/path:` with the path where are the `SECUIRTY-INSIGHTS.yml` of the project and `security-insights-schema.yml` which can be found [here](https://github.com/ossf/security-insights-spec/blob/main/security-insights-schema.yaml):

```
docker run -v $/your/path:/security-insights-validation -it sec-insights si-validator verify $/your/path/security-insights-validation/SECURITY-INSIGHTS.yml $/your/path/security-insights-validation/security-insights-schema.yaml
```

If you want to create a new `SECURITY-INSIGHTS.yml` by complying the YAML schema, launch this command:

```
docker run -v $/your/path:/security-insights-validation -it sec-insights si-validator create /security-insights-validation/SECURITY-INSIGHTS.yml /security-insights-validation/security-insights-schema.yaml
```

and fill out the questionnaire by following the wizard. The questions labeled with `(optional)` are not mandatory and can be skipped.

## Bugs

If you find any bugs, please open an issue or submit a pull request.

### Known Bugs

- [ ] Value type and format are not printed in the wizard
- [ ] The script accepts just `True` or `False` for boolean values (case sensitive)
- [ ] The script supports just single-line comments

## Security

If you find a security vulnerability, please report it via [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability).



