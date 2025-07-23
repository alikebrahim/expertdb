#!/usr/bin/env python3
"""
Professional Background Migration Script
Populates expert_experience_entries and expert_education_entries tables from normalized bio files.
"""

import os
import re
import sqlite3
from datetime import datetime
from typing import List, Dict, Optional, Tuple

def extract_expert_id(expert_id_str: str) -> int:
    """Extract numeric ID from expert ID string (E001 -> 1, E321 -> 321)"""
    match = re.match(r'E(\d+)', expert_id_str.strip())
    if match:
        return int(match.group(1))
    raise ValueError(f"Invalid expert ID format: {expert_id_str}")

def parse_experience_entry(line: str) -> Optional[Dict[str, str]]:
    """Parse a single experience line into components"""
    line = line.strip()
    if not line or line.startswith('Experience:'):
        return None
    
    # Handle cases where dates might be missing or incomplete
    if ':' not in line:
        return None
        
    parts = line.split(':', 1)
    if len(parts) != 2:
        return None
        
    date_part = parts[0].strip()
    rest_part = parts[1].strip()
    
    # Parse date range
    start_date = ""
    end_date = ""
    is_current = False
    
    if date_part:
        if '-' in date_part:
            date_range = date_part.split('-', 1)
            start_date = date_range[0].strip()
            end_date = date_range[1].strip()
            if end_date.lower() in ['present', 'current', 'ongoing']:
                is_current = True
                end_date = ""
        else:
            start_date = date_part
    
    # Parse position, organization, country
    position = ""
    organization = ""
    country = ""
    
    if rest_part:
        # Split by comma to separate position/org from country
        parts = rest_part.split(',')
        if len(parts) >= 3:
            # Position, Organization, Country
            position = parts[0].strip()
            organization = parts[1].strip() 
            country = parts[2].strip()
        elif len(parts) == 2:
            # Position, Organization (no country specified)
            position = parts[0].strip()
            organization = parts[1].strip()
        else:
            # Just position or organization
            if rest_part:
                position = rest_part.strip()
    
    return {
        'start_date': start_date,
        'end_date': end_date,
        'is_current': is_current,
        'position': position,
        'organization': organization,
        'country': country,
        'description': ''
    }

def parse_education_entry(line: str) -> Optional[Dict[str, str]]:
    """Parse a single education line into components"""
    line = line.strip()
    if not line or line.startswith('Education:'):
        return None
    
    if ':' not in line:
        return None
        
    parts = line.split(':', 1)
    if len(parts) != 2:
        return None
        
    year_part = parts[0].strip()
    rest_part = parts[1].strip()
    
    # Parse degree, institution, country
    degree = ""
    institution = ""
    country = ""
    field_of_study = ""
    
    if rest_part:
        # Split by comma to separate degree, institution, country
        parts = rest_part.split(',')
        if len(parts) >= 3:
            # Degree, Institution, Country
            degree = parts[0].strip()
            institution = parts[1].strip()
            country = parts[2].strip()
        elif len(parts) == 2:
            # Degree, Institution (no country)
            degree = parts[0].strip()
            institution = parts[1].strip()
        else:
            # Just degree
            degree = rest_part.strip()
    
    return {
        'graduation_year': year_part,
        'degree': degree,
        'institution': institution,
        'field_of_study': field_of_study,
        'country': country,
        'description': ''
    }

def parse_bio_file(file_path: str) -> List[Dict]:
    """Parse a single bio file and return list of experts with their entries"""
    experts = []
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split by expert entries (separated by ---)
    expert_sections = content.split('---')
    
    for section in expert_sections:
        section = section.strip()
        if not section:
            continue
            
        lines = [line.strip() for line in section.split('\n') if line.strip()]
        if len(lines) < 2:
            continue
            
        expert_data = {
            'id': None,
            'name': '',
            'experience_entries': [],
            'education_entries': []
        }
        
        current_section = None
        
        for line in lines:
            if line.startswith('ID:'):
                expert_id_str = line.replace('ID:', '').strip()
                try:
                    expert_data['id'] = extract_expert_id(expert_id_str)
                except ValueError as e:
                    print(f"Warning: {e} in file {file_path}")
                    continue
                    
            elif line.startswith('Name:'):
                expert_data['name'] = line.replace('Name:', '').strip()
                
            elif line.startswith('Education:'):
                current_section = 'education'
                
            elif line.startswith('Experience:'):
                current_section = 'experience'
                
            elif current_section == 'experience':
                exp_entry = parse_experience_entry(line)
                if exp_entry:
                    expert_data['experience_entries'].append(exp_entry)
                    
            elif current_section == 'education':
                edu_entry = parse_education_entry(line)
                if edu_entry:
                    expert_data['education_entries'].append(edu_entry)
        
        if expert_data['id'] is not None:
            experts.append(expert_data)
            
    return experts

def migrate_expert_entries(db_path: str, bio_files_dir: str):
    """Main migration function"""
    
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Statistics
    total_experts = 0
    total_experience_entries = 0
    total_education_entries = 0
    errors = []
    
    print(f"Starting migration from {bio_files_dir}")
    print(f"Database: {db_path}")
    
    # Process all bio files
    bio_files = sorted([f for f in os.listdir(bio_files_dir) if f.endswith('.md')])
    
    for bio_file in bio_files:
        file_path = os.path.join(bio_files_dir, bio_file)
        print(f"\nProcessing {bio_file}...")
        
        try:
            experts = parse_bio_file(file_path)
            
            for expert in experts:
                expert_id = expert['id']
                
                # Check if expert exists in database
                cursor.execute("SELECT id FROM experts WHERE id = ?", (expert_id,))
                if not cursor.fetchone():
                    print(f"  Warning: Expert {expert_id} not found in database, skipping...")
                    continue
                
                total_experts += 1
                
                # Insert experience entries
                for exp_entry in expert['experience_entries']:
                    try:
                        cursor.execute("""
                            INSERT INTO expert_experience_entries (
                                expert_id, organization, position, start_date, end_date, 
                                is_current, country, description, created_at, updated_at
                            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                        """, (
                            expert_id,
                            exp_entry['organization'],
                            exp_entry['position'],
                            exp_entry['start_date'],
                            exp_entry['end_date'],
                            exp_entry['is_current'],
                            exp_entry['country'],
                            exp_entry['description'],
                            datetime.now().isoformat(),
                            datetime.now().isoformat()
                        ))
                        total_experience_entries += 1
                    except sqlite3.Error as e:
                        error_msg = f"Error inserting experience for expert {expert_id}: {e}"
                        errors.append(error_msg)
                        print(f"  {error_msg}")
                
                # Insert education entries  
                for edu_entry in expert['education_entries']:
                    try:
                        cursor.execute("""
                            INSERT INTO expert_education_entries (
                                expert_id, institution, degree, field_of_study, 
                                graduation_year, country, description, created_at, updated_at
                            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                        """, (
                            expert_id,
                            edu_entry['institution'],
                            edu_entry['degree'],
                            edu_entry['field_of_study'],
                            edu_entry['graduation_year'],
                            edu_entry['country'],
                            edu_entry['description'],
                            datetime.now().isoformat(),
                            datetime.now().isoformat()
                        ))
                        total_education_entries += 1
                    except sqlite3.Error as e:
                        error_msg = f"Error inserting education for expert {expert_id}: {e}"
                        errors.append(error_msg)
                        print(f"  {error_msg}")
                        
        except Exception as e:
            error_msg = f"Error processing file {bio_file}: {e}"
            errors.append(error_msg)
            print(f"  {error_msg}")
    
    # Commit all changes
    conn.commit()
    conn.close()
    
    # Print final report
    print(f"\n" + "="*50)
    print("MIGRATION COMPLETE")
    print(f"="*50)
    print(f"Total experts processed: {total_experts}")
    print(f"Experience entries inserted: {total_experience_entries}")
    print(f"Education entries inserted: {total_education_entries}")
    print(f"Errors encountered: {len(errors)}")
    
    if errors:
        print("\nERRORS:")
        for error in errors:
            print(f"  - {error}")
    
    print(f"\nMigration completed successfully!")

if __name__ == "__main__":
    # Configuration
    DB_PATH = "db/sqlite/main.db"
    BIO_FILES_DIR = "bios"
    
    # Validate paths exist
    if not os.path.exists(DB_PATH):
        print(f"Error: Database file not found at {DB_PATH}")
        exit(1)
        
    if not os.path.exists(BIO_FILES_DIR):
        print(f"Error: Bio files directory not found at {BIO_FILES_DIR}")
        exit(1)
    
    # Run migration
    migrate_expert_entries(DB_PATH, BIO_FILES_DIR)