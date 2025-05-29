#!/bin/bash

# Output CSV file
output_file="python_files.csv"

# Write header to CSV
echo "id,folder_name,file_name,file_content" > "$output_file"

# Initialize ID counter
id=1

# Loop through each subfolder
for folder in */ ; do
  [ -d "$folder" ] || continue
  folder_name=$(basename "$folder")

  # Loop through each .py file in the subfolder
  for file in "$folder"*.py; do
    [ -f "$file" ] || continue
    file_name=$(basename "$file")

    # Read and escape content
    file_content=$(<"$file")
    file_content="${file_content//\"/\"\"}"  # Escape quotes

    # Write to CSV
    echo "$id,\"$folder_name\",\"$file_name\",\"$file_content\"" >> "$output_file"

    # Increment ID
    ((id++))
  done
done

echo "CSV file created: $output_file"

