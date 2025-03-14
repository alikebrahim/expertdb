import sqlite3
import csv

# Connect to SQLite database
db_path = "./db/sqlite/expertdb.sqlite"  # Maintained as per your setup
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Function to transform CSV data to match table schema
def transform_row(row):
    # Helper function to get value by key, case-insensitive and stripped
    def get_value(key):
        for k in row.keys():
            # Remove BOM and strip whitespace, compare case-insensitively
            cleaned_key = k.replace('\ufeff', '').strip()
            if cleaned_key.lower() == key.lower():
                return row[k].strip() if row[k] else None
        raise KeyError(f"Column '{key}' not found in CSV")

    return {
        "expert_id": get_value("ID"),
        "name": get_value("Name"),
        "designation": get_value("Designation"),
        "institution": get_value("Institution"),
        "is_bahraini": 1 if get_value("BH") == "Yes" else (0 if get_value("BH") == "No" else None),
        "nationality": "Bahraini" if get_value("BH") == "Yes" else ("Non-Bahraini" if get_value("BH") == "No" else "Unknown"),
        "is_available": 1 if get_value("Available") == "Yes" else (0 if get_value("Available") == "No" else None),
        "rating": get_value("Rating"),
        "role": get_value("Validator/ Evaluator"),
        "employment_type": get_value("Academic/Employer"),
        "general_area": get_value("General Area"),
        "specialized_area": get_value("Specialised Area"),
        "is_trained": 1 if get_value("Trained") == "Yes" else (0 if get_value("Trained") == "No" else None),
        "cv_path": get_value("CV") if get_value("CV") else None,
        "phone": get_value("Phone") if get_value("Phone") else None,
        "email": get_value("Email") if get_value("Email") else None,
        "is_published": 1 if get_value("Published") == "Yes" else (0 if get_value("Published") == "No" else None),
        "biography": None,
        "isced_level_id": None,
        "isced_field_id": None,
        "original_request_id": None,
        "updated_at": None
    }

# Read CSV and insert data
csv_file_path = "./experts.csv"  # Maintained as per your setup
with open(csv_file_path, newline='', encoding='utf-8-sig') as csvfile:  # 'utf-8-sig' skips BOM
    reader = csv.DictReader(csvfile)
    
    # Print headers for debugging (cleaned)
    cleaned_headers = [h.replace('\ufeff', '').strip() for h in reader.fieldnames]
    print("CSV Headers found (cleaned):", cleaned_headers)
    
    # Transform and collect rows
    rows = [transform_row(row) for row in reader]

    # Batch insert for efficiency
    cursor.executemany('''
        INSERT OR IGNORE INTO experts (
            expert_id, name, designation, institution, is_bahraini, nationality, 
            is_available, rating, role, employment_type, general_area, 
            specialized_area, is_trained, cv_path, phone, email, is_published,
            biography, isced_level_id, isced_field_id, original_request_id, updated_at
        ) VALUES (
            :expert_id, :name, :designation, :institution, :is_bahraini, :nationality,
            :is_available, :rating, :role, :employment_type, :general_area,
            :specialized_area, :is_trained, :cv_path, :phone, :email, :is_published,
            :biography, :isced_level_id, :isced_field_id, :original_request_id, :updated_at
        )
    ''', rows)

# Commit changes and verify
conn.commit()
cursor.execute("SELECT COUNT(*) FROM experts")
print(f"Total rows in experts table after import: {cursor.fetchone()[0]}")

# Close connection
conn.close()
