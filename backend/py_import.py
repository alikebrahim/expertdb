import sqlite3
import csv

# ANSI escape codes for red text in terminal
RED = "\033[31m"
RESET = "\033[0m"

# Connect to SQLite database
db_path = "./db/sqlite/expertdb.sqlite"
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Load expert_areas into a lookup dictionary
cursor.execute("SELECT id, name FROM expert_areas")
expert_areas = {row[1]: row[0] for row in cursor.fetchall()}
print("Loaded expert_areas:", expert_areas)

# Function to normalize text for matching
def normalize_text(text):
    if not text:
        return ""
    # Replace multiple spaces/hyphens with single space and standardize hyphen spacing
    text = ' '.join(text.split())  # Collapse multiple spaces
    text = text.replace('-', ' - ').replace('  ', ' ')  # Ensure single space around hyphen
    return text.strip()

# Function to transform CSV data to match table schema
def transform_row(row, expert_areas):
    def get_value(key):
        for k in row.keys():
            cleaned_key = k.replace('\ufeff', '').strip()
            if cleaned_key.lower() == key.lower():
                return row[k].strip() if row[k] else None
        raise KeyError(f"Column '{key}' not found in CSV")

    name = get_value("Name") or "Unknown Expert"
    general_area_text = get_value("General Area")
    
    # Handle missing or empty general_area_text
    if not general_area_text:
        # Default to "Unknown" and log
        cursor.execute("INSERT OR IGNORE INTO expert_areas (name) VALUES ('Unknown')")
        cursor.execute("SELECT id FROM expert_areas WHERE name = 'Unknown'")
        general_area_id = cursor.fetchone()[0]
        expert_areas["Unknown"] = general_area_id
        print(f"{RED}Warning: 'General Area' is missing or empty for expert {name}, using 'Unknown' (ID {general_area_id}){RESET}")
    else:
        # Normalize CSV text and expert_areas for matching
        normalized_text = normalize_text(general_area_text)
        general_area_id = expert_areas.get(normalized_text)

        # Handle unmatched cases
        if general_area_id is None:
            # Try additional normalization for known cases
            if "Science" in normalized_text and "Mathematics" in normalized_text:
                normalized_text = "Science - Mathematics"
                general_area_id = expert_areas.get(normalized_text)
            
            # If still no exact match, try partial match
            if general_area_id is None:
                for area_name, area_id in expert_areas.items():
                    normalized_area = normalize_text(area_name)
                    if normalized_text in normalized_area:
                        general_area_id = area_id
                        print(f"Matched '{general_area_text}' to '{area_name}' (ID {area_id}) for expert {name}")
                        break
            
            # If no match found, use Unknown
            if general_area_id is None:
                cursor.execute("INSERT OR IGNORE INTO expert_areas (name) VALUES ('Unknown')")
                cursor.execute("SELECT id FROM expert_areas WHERE name = 'Unknown'")
                general_area_id = cursor.fetchone()[0]
                expert_areas["Unknown"] = general_area_id
                print(f"{RED}Warning: '{general_area_text}' not found in expert_areas for expert {name}, using 'Unknown' (ID {general_area_id}){RESET}")
        else:
            print(f"Matched '{general_area_text}' to '{normalized_text}' (ID {general_area_id}) for expert {name}")

    return {
        "expert_id": get_value("ID"),
        "name": name,
        "designation": get_value("Designation"),
        "institution": get_value("Institution"),
        "is_bahraini": 1 if get_value("BH") == "Yes" else (0 if get_value("BH") == "No" else None),
        "nationality": "Bahraini" if get_value("BH") == "Yes" else ("Non-Bahraini" if get_value("BH") == "No" else "Unknown"),
        "is_available": 1 if get_value("Available") == "Yes" else (0 if get_value("Available") == "No" else None),
        "rating": get_value("Rating"),
        "role": get_value("Validator/ Evaluator"),
        "employment_type": get_value("Academic/Employer"),
        "general_area": general_area_id,
        "specialized_area": get_value("Specialised Area"),
        "is_trained": 1 if get_value("Trained") == "Yes" else (0 if get_value("Trained") == "No" else None),
        "cv_path": get_value("CV") if get_value("CV") else None,
        "phone": get_value("Phone") if get_value("Phone") else None,
        "email": get_value("Email") if get_value("Email") else None,
        "is_published": 1 if get_value("Published") == "Yes" else (0 if get_value("Published") == "No" else None),
        "biography": None,
        "original_request_id": None,
        "updated_at": None
    }

# Read CSV and insert data
csv_file_path = "./experts.csv"
with open(csv_file_path, newline='', encoding='utf-8-sig') as csvfile:
    reader = csv.DictReader(csvfile)
    cleaned_headers = [h.replace('\ufeff', '').strip() for h in reader.fieldnames]
    print("CSV Headers found (cleaned):", cleaned_headers)
    
    # Transform and collect rows
    rows = [transform_row(row, expert_areas) for row in reader]

    # Batch insert into experts table
    cursor.executemany('''
        INSERT OR IGNORE INTO experts (
            expert_id, name, designation, institution, is_bahraini, nationality, 
            is_available, rating, role, employment_type, general_area, 
            specialized_area, is_trained, cv_path, phone, email, is_published,
            biography, original_request_id, updated_at
        ) VALUES (
            :expert_id, :name, :designation, :institution, :is_bahraini, :nationality,
            :is_available, :rating, :role, :employment_type, :general_area,
            :specialized_area, :is_trained, :cv_path, :phone, :email, :is_published,
            :biography, :original_request_id, :updated_at
        )
    ''', rows)

# Commit changes and verify
conn.commit()
cursor.execute("SELECT COUNT(*) FROM experts")
print(f"Total rows in experts table after import: {cursor.fetchone()[0]}")

# Verify a few records, including known fail cases
cursor.execute("""
    SELECT expert_id, name, general_area, 
           (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name 
    FROM experts 
    WHERE expert_id IN ('E020', 'E059', 'E105', 'E112', 'E137', 'E211', 'E240', 'E341')
""")
print("\nVerified fail case records:")
for row in cursor.fetchall():
    print(row)

# Close connection
conn.close()
