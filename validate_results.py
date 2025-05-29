import csv
import requests

# Input and output file names
input_csv = "python_files.csv"
output_csv = "scanned_python_files.csv"

#input_csv = "cn_python_files.csv"
#output_csv = "cn_scanned_python_files.csv"

# Open the input and output CSV files
i = 0
with open(input_csv, newline='', encoding='utf-8') as infile, \
     open(output_csv, mode='w', newline='', encoding='utf-8') as outfile:
    i += 1
    reader = csv.DictReader(infile)
    fieldnames = reader.fieldnames + ['success', 'iterations', 'resulting_code']
    writer = csv.DictWriter(outfile, fieldnames=fieldnames)
    writer.writeheader()

    for row in reader:
        payload = {
            'filename': row['file_name'],
            'content': row['file_content']
        }

        try:
            response = requests.post("http://localhost:8080/scan", json=payload)
            response.raise_for_status()
            data = response.json()

            # Add response fields to the row
            row['success'] = data.get('success', False)
            row['iterations'] = data.get('iterations', 0)
            row['resulting_code'] = data.get('resulting_code', '').replace('"', '""')

        except Exception as e:
            # Handle request errors or missing data
            row['success'] = False
            row['iterations'] = 0
            row['resulting_code'] = f"Error: {e}"

        writer.writerow(row)
        print(f"Linha {i} finalizada!")

    

print(f"New CSV created: {output_csv}")
