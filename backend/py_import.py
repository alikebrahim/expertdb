import sqlite3
import csv
import re
import os

# ANSI escape codes for colored terminal output
RED = "\033[31m"
YELLOW = "\033[33m"
RESET = "\033[0m"

# Connect to SQLite database
db_path = "./db/sqlite/main.db"
try:
    # Check if the file exists first
    if not os.path.exists(db_path):
        print(f"{RED}Error: Database file '{db_path}' not found{RESET}")
        print(f"Make sure the main.db file exists in the db/sqlite directory")
        exit(1)
    
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    print(f"Connected to existing database: {db_path}")
except Exception as e:
    print(f"{RED}Error connecting to database: {str(e)}{RESET}")
    exit(1)

# Load expert_areas into a lookup dictionary
# cursor.execute("SELECT id, name FROM expert_areas")
# expert_areas = {row[1]: row[0] for row in cursor.fetchall()}
# print("Loaded expert_areas:", expert_areas)
expert_areas = {}

# Function to normalize text for matching
def normalize_text(text):
    if not text:
        return ""
    # Replace multiple spaces/hyphens with single space and standardize hyphen spacing
    text = ' '.join(text.split())  # Collapse multiple spaces
    text = text.replace('-', ' - ').replace('  ', ' ')  # Ensure single space around hyphen
    return text.strip()

# Function to generate unique expert_id (EXP-<sequence>) - Not needed as we'll use CSV IDs
def generate_expert_id(cursor):
    cursor.execute("SELECT MAX(CAST(SUBSTR(expert_id, 5) AS INTEGER)) FROM experts WHERE expert_id LIKE 'EXP-%'")
    max_seq = cursor.fetchone()[0] or 0
    return f"EXP-{max_seq + 1}"

# Function to validate and parse skills
def parse_skills(skills_text):
    if not skills_text:
        return []
    # Split by commas or semicolons, strip whitespace
    skills = [s.strip() for s in re.split(r'[;,]', skills_text) if s.strip()]
    return skills if skills else []

# Function to transform CSV data to match updated experts table schema
def transform_row(row, expert_areas, cursor):
    def get_value(key, default=None):
        # Try exact match first
        for k in row.keys():
            cleaned_key = k.replace('\ufeff', '').strip()
            if cleaned_key.lower() == key.lower():
                value = row[k].strip() if row[k] else None
                return value if value else default
        
        # Try partial match for better CSV compatibility
        for k in row.keys():
            cleaned_key = k.replace('\ufeff', '').strip()
            if key.lower() in cleaned_key.lower():
                value = row[k].strip() if row[k] else None
                return value if value else default
                
        return default

    name = get_value("Name", "Unknown Expert")
    general_area_text = get_value("General Area", None)

    # Just use a hardcoded value for general_area as we don't have the expert_areas table populated yet
    general_area_id = 1

    # Use existing ID directly from CSV
    expert_id = get_value("ID", None)
    if not expert_id or expert_id in [None, '', 'NULL']:
        # Only use generate_expert_id if absolutely necessary
        expert_id = generate_expert_id(cursor)
        print(f"Generated expert_id '{expert_id}' for expert {name}")

    # Set skills to empty
    skills = []

    # Handle biography data that might be in different columns
    biography = get_value("Biography", None)
    if biography is None:
        biography = get_value("Bio", None)
    if biography is None:
        biography = get_value("Profile", None)
    if biography is None:
        biography = get_value("About", None)
    if biography is None:
        biography = get_value("Description", None)
    if biography is None:
        biography = get_value("Comments", "No biography provided")

    # Determine designation based on title or name prefix
    designation = get_value("Designation", None)
    if designation is None:
        designation = get_value("Title", None)
    if designation is None:
        # Try to extract from the "." column which might contain title
        title_prefix = get_value(".", None)
        if title_prefix:
            designation = title_prefix.strip()
        else:
            # Extract Dr./Prof. prefix from name if present
            if name.startswith("Dr.") or name.startswith("DR."):
                designation = "Dr."
            elif name.startswith("Prof.") or name.startswith("PROF."):
                designation = "Prof."
            else:
                designation = "Unknown"

    # Prepare row with all required fields
    row_data = {
        "expert_id": expert_id,
        "name": name,
        "designation": designation,
        "institution": get_value("Institution", "Unknown"),
        "is_bahraini": 1 if get_value("BH", "No").lower() in ["yes", "y", "true"] else 0,
        "nationality": "Bahraini" if get_value("BH", "No").lower() in ["yes", "y", "true"] else "Non-Bahraini",
        "is_available": 1 if get_value("Available", "No").lower() in ["yes", "y", "true"] else 0,
        "rating": get_value("Rating", "0"),
        "role": get_value("Validator/ Evaluator", "evaluator").lower(),
        "employment_type": get_value("Academic/Employer", "Unknown"),
        "general_area": general_area_id,
        "specialized_area": get_value("Specialised Area", "Unknown"),
        "is_trained": 1 if get_value("Trained", "No").lower() in ["yes", "y", "true"] else 0,
        "cv_path": get_value("CV", "./data/documents/default_cv.pdf"),
        "phone": get_value("Phone", "00000000"),
        "email": get_value("Email", "unknown@example.com"),
        "is_published": 1 if get_value("Published", "No").lower() in ["yes", "y", "true"] else 0,
        "biography": biography,
        # Remove the skills field since it doesn't exist in the table
        "approval_document_path": get_value("Approval Document", "./data/documents/default_approval.pdf"),
        "original_request_id": None,
        "updated_at": None
    }

    # Set default values for required fields without warnings
    required_fields = [
        "expert_id", "name", "designation", "institution", "is_bahraini", "is_available",
        "rating", "role", "employment_type", "general_area", "specialized_area", "is_trained",
        "cv_path", "phone", "email", "is_published", "biography", "approval_document_path"
    ]
    for field in required_fields:
        if row_data[field] is None or row_data[field] == "":
            row_data[field] = row_data[field] if row_data[field] is not None else "Unknown"

    return row_data

# Make sure data/documents directory exists for default files
# documents_dir = "./data/documents"
# if not os.path.exists(documents_dir):
#     try:
#         os.makedirs(documents_dir, exist_ok=True)
#         print(f"Created documents directory: {documents_dir}")
#     except Exception as e:
#         print(f"{RED}Error creating documents directory: {str(e)}{RESET}")

# Create default files if they don't exist
# default_cv_path = "./data/documents/default_cv.pdf"
# default_approval_path = "./data/documents/default_approval.pdf"

# for default_path in [default_cv_path, default_approval_path]:
#     if not os.path.exists(default_path):
#         try:
#             # Create an empty file
#             with open(default_path, 'w') as f:
#                 f.write("This is a placeholder document file.")
#             print(f"Created default document: {default_path}")
#         except Exception as e:
#             print(f"{RED}Error creating default document {default_path}: {str(e)}{RESET}")

# Read CSV and insert data
csv_file_path = "./experts.csv"
processed_count = 0
error_count = 0

try:
    if not os.path.exists(csv_file_path):
        print(f"{RED}Error: CSV file '{csv_file_path}' not found{RESET}")
        exit(1)
        
    with open(csv_file_path, newline='', encoding='utf-8-sig') as csvfile:
        reader = csv.DictReader(csvfile)
        if not reader.fieldnames:
            print(f"{RED}Error: CSV file has no headers or is empty{RESET}")
            exit(1)
            
        cleaned_headers = [h.replace('\ufeff', '').strip() for h in reader.fieldnames]
        print("CSV Headers found (cleaned):", cleaned_headers)
        
        # Check minimum required columns
        required_columns = ["Name", "ID"]
        found_columns = [col for col in required_columns if any(col.lower() in header.lower() for header in cleaned_headers)]
        if len(found_columns) < len(required_columns):
            missing = [col for col in required_columns if col not in found_columns]
            print(f"{RED}Error: CSV is missing required columns: {missing}{RESET}")
            print(f"Available columns: {cleaned_headers}")
            exit(1)

        # Start transaction
        conn.execute("BEGIN TRANSACTION")
        
        # Transform and collect rows
        rows = []
        for i, row in enumerate(reader, 1):
            try:
                transformed_row = transform_row(row, expert_areas, cursor)
                rows.append(transformed_row)
                processed_count += 1
                
                # Process in batches of 100 to avoid large transactions
                if len(rows) >= 100:
                    try:
                        # Batch insert into experts table
                        cursor.executemany('''
                            INSERT OR IGNORE INTO experts (
                                expert_id, name, designation, institution, is_bahraini, nationality,
                                is_available, rating, role, employment_type, general_area,
                                specialized_area, is_trained, cv_path, phone, email, is_published,
                                biography, approval_document_path, original_request_id, updated_at
                            ) VALUES (
                                :expert_id, :name, :designation, :institution, :is_bahraini, :nationality,
                                :is_available, :rating, :role, :employment_type, :general_area,
                                :specialized_area, :is_trained, :cv_path, :phone, :email, :is_published,
                                :biography, :approval_document_path, :original_request_id, :updated_at
                            )
                        ''', rows)
                        conn.commit()
                        print(f"Imported batch of {len(rows)} records (total processed: {processed_count})")
                        rows = []  # Clear the batch
                    except Exception as e:
                        conn.rollback()
                        print(f"{RED}Error inserting batch: {str(e)}{RESET}")
                        error_count += len(rows)
                        rows = []  # Clear the failed batch and continue
                
            except Exception as e:
                print(f"{RED}Error processing row {i} for expert {row.get('Name', 'Unknown')}: {str(e)}{RESET}")
                error_count += 1

        # Process remaining rows
        if rows:
            try:
                # Batch insert into experts table
                cursor.executemany('''
                    INSERT OR IGNORE INTO experts (
                        expert_id, name, designation, institution, is_bahraini, nationality,
                        is_available, rating, role, employment_type, general_area,
                        specialized_area, is_trained, cv_path, phone, email, is_published,
                        biography, approval_document_path, original_request_id, updated_at
                    ) VALUES (
                        :expert_id, :name, :designation, :institution, :is_bahraini, :nationality,
                        :is_available, :rating, :role, :employment_type, :general_area,
                        :specialized_area, :is_trained, :cv_path, :phone, :email, :is_published,
                        :biography, :approval_document_path, :original_request_id, :updated_at
                    )
                ''', rows)
                conn.commit()
                print(f"Imported final batch of {len(rows)} records")
            except Exception as e:
                conn.rollback()
                print(f"{RED}Error inserting final batch: {str(e)}{RESET}")
                error_count += len(rows)

except FileNotFoundError:
    print(f"{RED}Error: CSV file '{csv_file_path}' not found{RESET}")
    exit(1)
except Exception as e:
    print(f"{RED}Error reading CSV file: {str(e)}{RESET}")
    exit(1)

# Verify results
try:
    cursor.execute("SELECT COUNT(*) FROM experts")
    total_experts = cursor.fetchone()[0]
    print(f"\n=== Import Summary ===")
    print(f"Total records processed: {processed_count}")
    print(f"Records with errors: {error_count}")
    print(f"Total experts in database: {total_experts}")

    # Verify some sample records
    print("\nSample records from database:")
    cursor.execute("""
        SELECT expert_id, name, general_area,
               (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name,
               approval_document_path
        FROM experts
        ORDER BY RANDOM()
        LIMIT 5
    """)
    sample_records = cursor.fetchall()
    for row in sample_records:
        print(row)
        
    # Check for specific IDs from the CSV
    print("\nChecking for specific IDs from CSV:")
    cursor.execute("""
        SELECT expert_id, name, general_area,
               (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name
        FROM experts
        WHERE expert_id LIKE 'E%' 
        LIMIT 3
    """)
    id_records = cursor.fetchall()
    for row in id_records:
        print(row)
        
    # Check for generated IDs
    print("\nChecking for generated expert IDs:")
    cursor.execute("""
        SELECT expert_id, name, general_area,
               (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name
        FROM experts
        WHERE expert_id LIKE 'EXP-%' 
        LIMIT 3
    """)
    exp_records = cursor.fetchall()
    for row in exp_records:
        print(row)

except Exception as e:
    print(f"{RED}Error verifying results: {str(e)}{RESET}")

finally:
    # Close connection
    conn.close()
    print("\nDatabase connection closed")
    
print(f"\nImport process completed with {error_count} errors out of {processed_count} records processed.")
if error_count > 0:
    print(f"{YELLOW}Please check the logs above for details on errors.{RESET}")
else:
    print(f"All records processed successfully!")
