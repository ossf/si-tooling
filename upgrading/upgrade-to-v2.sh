#!/bin/bash
VI_INSIGHTS=$1
PROJECT_NAME=$2

V1_SCHEMA_URL=$3
# Check if the user provided a v1 schema URL
if [ -z "$V1_SCHEMA_URL" ]; then
    V1_SCHEMA_URL='https://raw.githubusercontent.com/ossf/security-insights-spec/main/schema/v1-schema.cue'
    # V1_SCHEMA_URL='https://raw.githubusercontent.com/trumant/security-insights-spec/refs/heads/convert-v1-to-v2/v1-schema.cue'
fi

'https://raw.githubusercontent.com/ossf/security-insights-spec/main/schema/v1-schema.cue'
# check if the user provided a v2 schema URL
V2_SCHEMA_URL=$4
if [ -z "$V2_SCHEMA_URL" ]; then
    V2_SCHEMA_URL='https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/schema.cue'
    # V2_SCHEMA_URL='https://raw.githubusercontent.com/trumant/security-insights-spec/refs/heads/convert-v1-to-v2/schema.cue'
fi

if [ -f "v1-schema.cue" ]; then
    echo "Found local v1-schema.cue, skipping download"
else
    echo "Downloading v1 schema"
    curl "$V1_SCHEMA_URL" -o v1-schema.cue --silent
fi

# Check if insights-v1.yml exists locally first
if [ -f "insights-v1.yml" ]; then
    echo "Found local insights-v1.yml, skipping download"
else
    echo "Downloading v1 insights"
    curl "$VI_INSIGHTS" -o insights-v1.yml --silent
fi

# double-check that it actually conforms to the v1 schema
echo "Validating v1 insights against the v1 schema"

# convert the v1 data to v2 data
if ! cue vet insights-v1.yml v1-schema.cue -d '#v1'; then
    echo "Error: v1 insights file failed schema validation. Please correct the errors in v1 data before attempting another conversion"
    exit 1
fi

if [ -f "schema.cue" ]; then
    echo "Found local schema.cue, skipping v2 schema download"
else
    echo "Downloading v2 schema"
    curl -L "$V2_SCHEMA_URL" -o schema.cue --silent
fi

echo "Converting v1 data to v2 insights and saving to security-insights.yml"
cue export .:security_insights_spec -l input: insights-v1.yml -e output --out yaml -t project="$PROJECT_NAME" > security-insights.yml

echo "Validating v2 insights against the v2 schema"
cue vet security-insights.yml schema.cue -d '#v2'

echo ""
echo ""
echo "Thank you for using the v2 upgrade script."
echo ""
echo "The v1 insights data in $VI_INSIGHTS has been converted to the v2 schema and saved to security-insights.yml. You must review the file and make any necessary adjustments to resolve TODOs before using it."