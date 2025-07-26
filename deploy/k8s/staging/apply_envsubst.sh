#!/usr/bin/env bash
# This script loops through each file in the template folder,
# applies envsubst, and writes the output to the base folder.

TEMPLATE_DIR="template"
OUTPUT_DIR="base"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

for file in "$TEMPLATE_DIR"/*; do
    filename=$(basename "$file")
    envsubst < "$file" > "$OUTPUT_DIR/$filename"
    echo "Processed $file -> $OUTPUT_DIR/$filename"
done