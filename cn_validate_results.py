import pandas as pd
import os
import subprocess
import csv
import io

# Load the dataset
df = pd.read_json(
    "hf://datasets/CyberNative/Code_Vulnerability_Security_DPO/secure_programming_dpo.json",
    lines=True
)

# Filter for Python code (case insensitive)
filtered_df = df[df["lang"].str.lower() == "python"]

print(len(filtered_df))

# Reset index to ensure a clean sequence for filenames
filtered_df = filtered_df.reset_index(drop=True)

# Write to CSV
with open("cn_python_files.csv", mode="w", newline='', encoding="utf-8") as csvfile:
    writer = csv.DictWriter(csvfile, fieldnames=["file_name", "file_content"])
    writer.writeheader()
    
    for idx, row in filtered_df.iterrows():
        writer.writerow({
            "file_name": f"{idx + 1}",
            "file_content": row["rejected"]
        })