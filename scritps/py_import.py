import sqlite3
import csv
import re
import os
import sys

# ANSI escape codes for colored terminal output
RED = "\033[31m"
YELLOW = "\033[33m"
RESET = "\033[0m"

# Check for test mode and CV import mode
TEST_MODE = "--test" in sys.argv
CV_IMPORT_MODE = "--cvs" in sys.argv
TEST_LIMIT = 10 if TEST_MODE else None

if TEST_MODE and CV_IMPORT_MODE:
    print("Running in TEST MODE for CV IMPORT - will only process first 10 CV files")
elif TEST_MODE:
    print("Running in TEST MODE - will only process first 10 records")
elif CV_IMPORT_MODE:
    print("Running in CV IMPORT MODE - will import CV documents")


# Function to normalize text for matching
def normalize_text(text):
    if not text:
        return ""
    # Replace multiple spaces/hyphens with single space and standardize hyphen spacing
    text = " ".join(text.split())  # Collapse multiple spaces
    text = text.replace("-", " - ").replace(
        "  ", " "
    )  # Ensure single space around hyphen
    return text.strip()


def lookup_general_area(area_text, expert_areas):
    """
    Lookup general area ID by text with fuzzy matching
    """
    if not area_text or area_text.strip() == "":
        return 1  # Default to "Business" if no area specified

    area_text = area_text.strip()

    # Try exact match first
    if area_text in expert_areas:
        return expert_areas[area_text]

    # Try normalized match
    normalized = normalize_text(area_text)
    if normalized in expert_areas:
        return expert_areas[normalized]

    # Try partial matching for common cases
    area_lower = area_text.lower()
    for area_name, area_id in expert_areas.items():
        if isinstance(area_name, str) and area_name.lower() in area_lower:
            return area_id
        if isinstance(area_name, str) and area_lower in area_name.lower():
            return area_id

    # Fallback mappings for common CSV variations (updated for reorganized IDs)
    fallback_mapping = {
        "business": 1,  # Business
        "education": 12,  # Education
        "engineering": 14,  # Engineering
        "information technology": 25,  # Information Technology
        "it": 25,  # Information Technology
        "science": 26,  # Science
        "medical": 35,  # Medical Science
        "law": 39,  # Law
        "health": 45,  # Health & Safety
        "aviation": 44,  # Aviation
        "art": 40,  # Art and Design
        "design": 40,  # Art and Design
    }

    for key, area_id in fallback_mapping.items():
        if key in area_lower:
            print(f"Using fallback mapping: '{area_text}' → {key} (ID: {area_id})")
            return area_id

    print(f"Warning: Could not map area '{area_text}', defaulting to Business (ID: 1)")
    return 1  # Default to "Business"


def convert_rating(rating_text):
    """
    Convert CSV rating text to database integer - all default to 0
    """
    return 0  # All ratings default to 0 as requested


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

# Get or create SYSTEM user for audit trail
system_user_id = None
try:
    cursor.execute("SELECT id FROM users WHERE name = ? AND email = ?", ("SYSTEM", "system@expertdb.internal"))
    result = cursor.fetchone()
    if result:
        system_user_id = result[0]
        print(f"Found existing SYSTEM user with ID: {system_user_id}")
    else:
        # Create special SYSTEM user for import operations
        # Use "system" as role to distinguish from regular users
        cursor.execute("""
            INSERT INTO users (name, email, password_hash, role, is_active, created_at)
            VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
        """, ("SYSTEM", "system@expertdb.internal", "SYSTEM-USER-NO-LOGIN", "system", 0))
        system_user_id = cursor.lastrowid
        conn.commit()
        print(f"Created special SYSTEM user with ID: {system_user_id} for import operations")
        print(f"  - Role: system (special type)")
        print(f"  - Status: inactive (cannot login)")
        print(f"  - Purpose: audit trail for automated operations")
except Exception as e:
    print(f"{YELLOW}Warning: Could not create/find SYSTEM user: {str(e)}{RESET}")
    print(f"{YELLOW}Will attempt to use existing admin user as fallback{RESET}")
    # Try to find first admin user as fallback
    try:
        cursor.execute("SELECT id FROM users WHERE role IN ('admin', 'super_user') ORDER BY id LIMIT 1")
        fallback_result = cursor.fetchone()
        if fallback_result:
            system_user_id = fallback_result[0]
            print(f"Using fallback admin user ID: {system_user_id}")
        else:
            system_user_id = 1  # Final fallback
            print(f"Using final fallback user ID: {system_user_id}")
    except:
        system_user_id = 1  # Final fallback

# Load expert_areas into a lookup dictionary with normalization
cursor.execute("SELECT id, name FROM expert_areas")
expert_areas = {}
for row in cursor.fetchall():
    area_id, area_name = row
    # Store both exact and normalized versions for flexible matching
    expert_areas[area_name] = area_id
    expert_areas[normalize_text(area_name)] = area_id
print(f"Loaded {len(expert_areas) // 2} expert_areas with normalized lookup")


# Function to validate and parse skills
def parse_skills(skills_text):
    if not skills_text:
        return []
    # Split by commas or semicolons, strip whitespace
    skills = [s.strip() for s in re.split(r"[;,]", skills_text) if s.strip()]
    return skills if skills else []


def validate_csv_data(csv_file_path):
    """
    Pre-validate CSV data for common issues including malformed rows
    """
    issues = []

    with open(csv_file_path, newline="", encoding="utf-8-sig") as csvfile:
        reader = csv.DictReader(csvfile)

        # Check for duplicate IDs and malformed rows
        ids_seen = set()
        for i, row in enumerate(reader, 2):  # Start from row 2 (after header)
            # Check for None keys (indicates malformed CSV structure)
            none_keys = [k for k in row.keys() if k is None]
            if none_keys:
                # Try to get expert name for better error reporting
                expert_name = "Unknown"
                for k, v in row.items():
                    if k and "name" in str(k).lower():
                        expert_name = v if v else "Unknown"
                        break
                issues.append(
                    f"Row {i} ({expert_name}): Malformed structure - has {len(none_keys)} None keys"
                )

            csv_id = None
            for k in row.keys():
                if k is None:
                    continue
                if k.replace("\ufeff", "").strip().lower() == "id":
                    csv_id = row[k].strip() if row[k] else None
                    break

            if not csv_id:
                issues.append(f"Row {i}: Missing ID")
            elif csv_id in ids_seen:
                issues.append(f"Row {i}: Duplicate ID '{csv_id}'")
            else:
                ids_seen.add(csv_id)

    return issues


def import_cv_documents(cursor, test_limit=None):
    """
    Import CV documents from filesystem into expert_documents table
    Expected file pattern: ID{####}.pdf (e.g., ID0001.pdf, ID0002.pdf)
    """
    cv_dir = "./data/documents/cvs"
    processed_count = 0
    error_count = 0
    error_log = []

    # Create directory if it doesn't exist
    os.makedirs(cv_dir, exist_ok=True)

    # Check if directory exists after creation
    if not os.path.exists(cv_dir):
        print(f"{RED}CV directory could not be created: {cv_dir}{RESET}")
        return 0, 0

    # Find all PDF files matching ID pattern
    cv_files = []
    pattern = re.compile(r"^ID(\d{4})\.pdf$", re.IGNORECASE)

    for filename in os.listdir(cv_dir):
        match = pattern.match(filename)
        if match:
            expert_id = int(match.group(1))
            file_path = os.path.join(cv_dir, filename)
            cv_files.append((expert_id, filename, file_path))

    cv_files.sort(key=lambda x: x[0])  # Sort by expert ID

    if not cv_files:
        print(
            f"{YELLOW}No CV files found matching pattern ID####.pdf in {cv_dir}{RESET}"
        )
        return 0, 0

    print(f"Found {len(cv_files)} CV files to process")

    # Process CV files
    for expert_id, filename, file_path in cv_files:
        if test_limit and processed_count >= test_limit:
            print(f"Test mode: Stopping after {test_limit} CV files")
            break

        try:
            # Check if expert exists
            cursor.execute("SELECT id FROM experts WHERE id = ?", (expert_id,))
            if not cursor.fetchone():
                error_msg = f"Expert ID {expert_id} not found for file {filename}"
                print(f"{RED}{error_msg}{RESET}")
                error_log.append(error_msg)
                error_count += 1
                continue

            # Check if CV document already exists
            cursor.execute(
                """
                SELECT id FROM expert_documents 
                WHERE expert_id = ? AND document_type = 'cv'
            """,
                (expert_id,),
            )

            if cursor.fetchone():
                error_msg = f"CV document already exists for expert ID {expert_id}, skipping {filename}"
                print(f"{YELLOW}{error_msg}{RESET}")
                error_log.append(error_msg)
                continue

            # Get file stats
            file_stats = os.stat(file_path)
            file_size = file_stats.st_size

            # Insert into expert_documents table
            cursor.execute(
                """
                INSERT INTO expert_documents (
                    expert_id, document_type, filename, file_path, 
                    content_type, file_size, upload_date
                ) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
            """,
                (expert_id, "cv", filename, file_path, "application/pdf", file_size),
            )

            processed_count += 1
            print(
                f"Imported CV for expert ID {expert_id}: {filename} ({file_size} bytes)"
            )

        except Exception as e:
            error_msg = f"Error processing {filename}: {str(e)}"
            print(f"{RED}{error_msg}{RESET}")
            error_log.append(error_msg)
            error_count += 1

    # Generate summary
    print(f"\n=== CV Import Summary ===")
    print(f"Total CV files processed: {processed_count}")
    print(f"Files with errors: {error_count}")

    # Write error log if there are errors
    if error_log:
        error_log_path = "./cv_import_errors.log"
        with open(error_log_path, "w") as f:
            f.write("CV Import Error Log\n")
            f.write("==================\n\n")
            for error in error_log:
                f.write(f"{error}\n")
        print(f"{YELLOW}Error log written to: {error_log_path}{RESET}")

    return processed_count, error_count


# Function to transform CSV data to match updated experts table schema
def transform_row(row, expert_areas, cursor, system_user_id):
    # Check for malformed row with None keys
    none_keys = [k for k in row.keys() if k is None]
    if none_keys:
        expert_name = "Unknown"
        for k, v in row.items():
            if k and "name" in str(k).lower():
                expert_name = v if v else "Unknown"
                break
        print(
            f"{YELLOW}Warning: Malformed row detected for expert '{expert_name}' - has {len(none_keys)} None keys. Row data: {dict(row)}{RESET}"
        )

    # Validate required fields early
    csv_id = None
    for k in row.keys():
        if k is None:
            continue
        cleaned_key = k.replace("\ufeff", "").strip()
        if cleaned_key.lower() == "id":
            csv_id = row[k].strip() if row[k] else None
            break

    if not csv_id:
        print(f"Skipping row: Missing ID")
        return None

    def get_value(key, default=None):
        # Try exact match first
        for k in row.keys():
            # Skip None keys that can occur with malformed CSV rows
            if k is None:
                continue
            cleaned_key = k.replace("\ufeff", "").strip()
            if cleaned_key.lower() == key.lower():
                value = row[k]
                if value is not None:
                    value = value.strip() if value else None
                return value if value else default

        # Try partial match for better CSV compatibility
        for k in row.keys():
            # Skip None keys that can occur with malformed CSV rows
            if k is None:
                continue
            cleaned_key = k.replace("\ufeff", "").strip()
            if key.lower() in cleaned_key.lower():
                value = row[k]
                if value is not None:
                    value = value.strip() if value else None
                return value if value else default

        return default

    name = get_value("Name", "Unknown Expert")
    general_area_text = get_value("General Area", None)

    # Lookup general area ID using the loaded expert_areas
    general_area_id = lookup_general_area(general_area_text, expert_areas)

    # Extract ID and map to database sequence
    csv_id = get_value("ID", None)
    if not csv_id or csv_id in [None, "", "NULL"]:
        print(f"Error: Missing ID for expert {name}")
        return None

    # Extract numeric part for database ID
    import re

    numeric_match = re.search(r"(\d+)", csv_id)
    if not numeric_match:
        print(f"Error: Invalid ID format '{csv_id}' for expert {name}")
        return None

    sequence_id = int(numeric_match.group(1))

    # Set skills to empty
    skills = []

    # Biography field removed - not needed for database schema

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
        "id": sequence_id,  # Add explicit ID for database
        "name": name,
        "designation": designation,
        "affiliation": get_value("Institution", "Unknown"),
        "is_bahraini": 1
        if get_value("BH", "No").lower() in ["yes", "y", "true"]
        else 0,
        "is_available": 1
        if get_value("Available", "No").lower() in ["yes", "y", "true"]
        else 0,
        "rating": convert_rating(get_value("Rating", "No Rating")),
        "role": get_value("Validator/ Evaluator", "evaluator").lower(),
        "employment_type": get_value("Academic/Employer", "Unknown"),
        "general_area": general_area_id,
        "specialized_area": normalize_specialized_areas_to_ids(
            get_value("Specialised Area", ""), general_area_text, cursor
        ),
        "is_trained": 1
        if get_value("Trained", "No").lower() in ["yes", "y", "true"]
        else 0,
        "cv_document_id": None,  # Documents managed via expert_documents table
        "phone": get_value("Phone", "00000000"),
        "email": get_value("Email", "unknown@example.com"),
        "is_published": 1
        if get_value("Website", "No").lower() in ["yes", "y", "true"]
        else 0,
        # Remove the skills field since it doesn't exist in the table
        "approval_document_id": None,  # Documents managed via expert_documents table
        "original_request_id": None,
        "updated_at": None,
        "last_edited_by": system_user_id,  # Mark as edited by SYSTEM user
        "last_edited_at": None,  # Will be set to CURRENT_TIMESTAMP in SQL
    }

    # Set default values for required fields without warnings
    required_fields = [
        "name",
        "designation",
        "affiliation",
        "is_bahraini",
        "is_available",
        "rating",
        "role",
        "employment_type",
        "general_area",
        "specialized_area",
        "is_trained",
        "phone",
        "email",
        "is_published",
    ]
    for field in required_fields:
        if row_data[field] is None or row_data[field] == "":
            row_data[field] = (
                row_data[field] if row_data[field] is not None else "Unknown"
            )

    return row_data


def normalize_specialized_areas_to_ids(areas_text, general_area_text, cursor):
    """
    Process and normalize specialized areas from CSV, returning comma-separated IDs

    Args:
        areas_text (str): Raw specialized areas from CSV (may contain multiple areas)
        general_area_text (str): Expert's general area to exclude from specialized areas
        cursor: Database cursor for inserting into specialized_areas table

    Returns:
        str: Comma-separated specialized area IDs (e.g., "1,4,6")
    """
    if not areas_text or areas_text.strip() == "":
        return ""

    # Normalization mapping dictionary
    normalization_map = {
        "IT": "Information Technology",
        "Business Adminstration": "Business Administration",
        "Human Resource": "Human Resources",
        "Banking and Finance": ["Banking", "Finance"],
        "Accounting and Finance": ["Accounting", "Finance"],
        "Marketing and Management": ["Marketing", "Management"],
        "Process Instrumentation and Control Engineering": "Process Engineering",
        "Electrical and Electronics Engineering": [
            "Electrical Engineering",
            "Electronics Engineering",
        ],
        "Electrical & Electronics Engineering": [
            "Electrical Engineering",
            "Electronics Engineering",
        ],
        "Business Management/Business Adminstration": [
            "Business Management",
            "Business Administration",
        ],
        "SAP Analyst": "Information Technology",
        "Team Management": "Management",
        "Supervisory": "Management",
        "Article writer": "Writing",
        "Comparative Jurisprudence": "Jurisprudence",
        "Shari'ah": "Islamic Law",
        "Consultant": "Consulting",
        "Surgical and Medical Instrument": "Medical Instrumentation",
        "enviormental engineering": "Environmental Engineering",
        "Multimedia Technology": "Multimedia",
        "Artificial Inteligence": "Artificial Intelligence",
        "Acturial Science": "Actuarial Science",
        "Anti-Money Laundering and Compliance": "Anti-Money Laundering",
        "Quality Assurance/English": ["Quality Assurance", "English"],
        "Education Quality Assurance": "Education Quality Assurance",
        "Interior Architecture / Interior Design": [
            "Interior Architecture",
            "Interior Design",
        ],
        "Fashion Technology/ Fashion Design by Computer/ Luxury Fashion Management": [
            "Fashion Technology",
            "Fashion Design",
            "Luxury Management",
        ],
    }

    # Split on "/" and process each area
    areas = [area.strip() for area in areas_text.split("/") if area.strip()]
    normalized_areas = []

    for area in areas:
        # Apply normalization mapping
        if area in normalization_map:
            normalized = normalization_map[area]
            # Handle multi-area normalizations
            if isinstance(normalized, list):
                normalized_areas.extend(normalized)
            else:
                normalized_areas.append(normalized)
        else:
            normalized_areas.append(area)

    # Remove duplicates and areas that match general area
    unique_areas = []
    general_area_lower = general_area_text.lower() if general_area_text else ""

    for area in normalized_areas:
        area_lower = area.lower()
        # Skip if area matches general area or is already in unique list
        if area not in unique_areas and area_lower != general_area_lower:
            # Also skip if area is a subset of general area (e.g., skip "Engineering" if general is "Engineering")
            if not (general_area_lower and area_lower in general_area_lower):
                unique_areas.append(area)

    # Insert areas into specialized_areas table and collect IDs
    area_ids = []
    for area in unique_areas:
        # Check if area already exists
        cursor.execute("SELECT id FROM specialized_areas WHERE name = ?", (area,))
        result = cursor.fetchone()

        if result:
            area_ids.append(str(result[0]))
        else:
            # Insert new area
            cursor.execute(
                "INSERT INTO specialized_areas (name, created_at) VALUES (?, CURRENT_TIMESTAMP)",
                (area,),
            )
            area_id = cursor.lastrowid
            area_ids.append(str(area_id))

    # Return comma-separated IDs
    return ",".join(area_ids) if area_ids else ""


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

# Main execution logic
if CV_IMPORT_MODE:
    # CV Import Mode
    try:
        processed_count, error_count = import_cv_documents(cursor, TEST_LIMIT)

        # Commit any pending transactions
        conn.commit()

        print(
            f"\nCV import process completed with {error_count} errors out of {processed_count} files processed."
        )
        if error_count > 0:
            print(f"{YELLOW}Please check the error log for details.{RESET}")
        else:
            print(f"All CV files processed successfully!")

    except Exception as e:
        print(f"{RED}Error during CV import: {str(e)}{RESET}")
        conn.rollback()
    finally:
        conn.close()
        print("\nDatabase connection closed")

    exit(0)

# Expert Data Import Mode (default)
csv_file_path = "./files/experts.csv"
processed_count = 0
error_count = 0
error_log = []

try:
    if not os.path.exists(csv_file_path):
        print(f"{RED}Error: CSV file '{csv_file_path}' not found{RESET}")
        exit(1)

    # Pre-validate CSV data
    print("Validating CSV data...")
    validation_issues = validate_csv_data(csv_file_path)
    if validation_issues:
        print(f"{YELLOW}Found {len(validation_issues)} validation issues:{RESET}")
        for issue in validation_issues:
            print(f"  - {issue}")
        print("Continuing with import (issues logged)...")
    else:
        print("CSV validation passed.")

    with open(csv_file_path, newline="", encoding="utf-8-sig") as csvfile:
        reader = csv.DictReader(csvfile)
        if not reader.fieldnames:
            print(f"{RED}Error: CSV file has no headers or is empty{RESET}")
            exit(1)

        cleaned_headers = [h.replace("\ufeff", "").strip() for h in reader.fieldnames]
        print("CSV Headers found (cleaned):", cleaned_headers)

        # Check minimum required columns
        required_columns = ["Name", "ID"]
        found_columns = [
            col
            for col in required_columns
            if any(col.lower() in header.lower() for header in cleaned_headers)
        ]
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
            if TEST_MODE and processed_count >= TEST_LIMIT:
                print(f"Test mode: Stopping after {TEST_LIMIT} records")
                break

            try:
                transformed_row = transform_row(row, expert_areas, cursor, system_user_id)
                if transformed_row:  # Only append if transformation succeeded
                    rows.append(transformed_row)
                    processed_count += 1
                else:
                    error_log.append(
                        f"Row {i}: Failed to transform row - missing required data"
                    )
                    error_count += 1

                # Process in batches of 100 to avoid large transactions
                if len(rows) >= 100:
                    try:
                        # Batch insert into experts table
                        cursor.executemany(
                            """
                            INSERT OR REPLACE INTO experts (
                                id, name, designation, affiliation, is_bahraini,
                                is_available, rating, role, employment_type, general_area,
                                specialized_area, is_trained, cv_document_id, phone, email, is_published,
                                approval_document_id, original_request_id, updated_at, last_edited_by, last_edited_at
                            ) VALUES (
                                :id, :name, :designation, :affiliation, :is_bahraini,
                                :is_available, :rating, :role, :employment_type, :general_area,
                                :specialized_area, :is_trained, :cv_document_id, :phone, :email, :is_published,
                                :approval_document_id, :original_request_id, :updated_at, :last_edited_by, CURRENT_TIMESTAMP
                            )
                        """,
                            rows,
                        )
                        conn.commit()
                        print(
                            f"Imported batch of {len(rows)} records (total processed: {processed_count})"
                        )
                        rows = []  # Clear the batch
                    except Exception as e:
                        conn.rollback()
                        print(f"{RED}Error inserting batch: {str(e)}{RESET}")
                        error_count += len(rows)
                        rows = []  # Clear the failed batch and continue

            except Exception as e:
                error_msg = f"Row {i}: Error processing expert {row.get('Name', 'Unknown')}: {str(e)}"
                print(f"{RED}{error_msg}{RESET}")
                error_log.append(error_msg)
                error_count += 1

        # Process remaining rows
        if rows:
            try:
                # Batch insert into experts table
                cursor.executemany(
                    """
                    INSERT OR REPLACE INTO experts (
                        id, name, designation, affiliation, is_bahraini,
                        is_available, rating, role, employment_type, general_area,
                        specialized_area, is_trained, cv_document_id, phone, email, is_published,
                        approval_document_id, original_request_id, updated_at, last_edited_by, last_edited_at
                    ) VALUES (
                        :id, :name, :designation, :affiliation, :is_bahraini,
                        :is_available, :rating, :role, :employment_type, :general_area,
                        :specialized_area, :is_trained, :cv_document_id, :phone, :email, :is_published,
                        :approval_document_id, :original_request_id, :updated_at, :last_edited_by, CURRENT_TIMESTAMP
                    )
                """,
                    rows,
                )
                conn.commit()
                print(f"Imported final batch of {len(rows)} records")

            except Exception as e:
                conn.rollback()
                print(f"{RED}Error inserting final batch: {str(e)}{RESET}")
                error_count += len(rows)

        # Reset auto-increment sequence to continue from last imported ID
        cursor.execute("SELECT MAX(id) FROM experts")
        max_id = cursor.fetchone()[0] or 0
        cursor.execute(
            f"UPDATE sqlite_sequence SET seq = {max_id} WHERE name = 'experts'"
        )
        conn.commit()
        print(f"Reset auto-increment sequence to {max_id}")

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
        SELECT id, name, general_area,
               (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name,
               approval_document_id
        FROM experts
        ORDER BY RANDOM()
        LIMIT 5
    """)
    sample_records = cursor.fetchall()
    for row in sample_records:
        print(row)

    # Check first few records by ID
    print("\nFirst few imported records:")
    cursor.execute("""
        SELECT id, name, general_area,
               (SELECT name FROM expert_areas WHERE id = experts.general_area) AS area_name
        FROM experts
        ORDER BY id
        LIMIT 5
    """)
    first_records = cursor.fetchall()
    for row in first_records:
        print(row)

except Exception as e:
    print(f"{RED}Error verifying results: {str(e)}{RESET}")

finally:
    # Close connection
    conn.close()
    print("\nDatabase connection closed")

print(
    f"\nImport process completed with {error_count} errors out of {processed_count} records processed."
)

# Write error log if there are errors
if error_log:
    error_log_path = "./import_errors.log"
    with open(error_log_path, "w") as f:
        f.write("Import Error Log\n")
        f.write("================\n\n")
        for error in error_log:
            f.write(f"{error}\n")
    print(f"{YELLOW}Error log written to: {error_log_path}{RESET}")
    print(f"{YELLOW}Please review errors and fix CSV data if needed.{RESET}")
else:
    print(f"All records processed successfully!")
