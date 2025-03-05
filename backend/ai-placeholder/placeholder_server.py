from flask import Flask, request, jsonify
import os
import time
import random
import json

app = Flask(__name__)

# Get port from environment or use default
PORT = int(os.environ.get('PORT', 9000))

@app.route('/')
def home():
    return jsonify({
        "service": "ExpertDB AI Service Placeholder",
        "status": "running",
        "endpoints": [
            "/generate-profile",
            "/suggest-isced",
            "/extract-skills"
        ]
    })

@app.route('/generate-profile', methods=['POST'])
def generate_profile():
    data = request.json
    
    # Simulate processing time
    time.sleep(1)
    
    # Generate a mock profile based on input
    name = data.get('name', 'Expert')
    designation = data.get('designation', 'Professional')
    institution = data.get('institution', 'Organization')
    general_area = data.get('generalArea', 'Expertise')
    
    profile = f"{name} is a highly qualified {designation} at {institution} with extensive experience in {general_area}. "
    profile += f"They have contributed to numerous projects and have demonstrated exceptional skills in problem-solving and collaboration. "
    profile += f"Their expertise in {general_area} has been recognized through various professional achievements."
    
    return jsonify({
        "result": profile,
        "confidence_score": random.uniform(0.85, 0.98)
    })

@app.route('/suggest-isced', methods=['POST'])
def suggest_isced():
    data = request.json
    
    # Simulate processing time
    time.sleep(0.5)
    
    general_area = data.get('generalArea', '').lower()
    specialized_area = data.get('specializedArea', '').lower()
    
    # Simple mapping logic (would be replaced with actual AI in production)
    isced_mapping = {
        'computer': {'code': '06', 'name': 'Information and Communication Technologies'},
        'software': {'code': '06', 'name': 'Information and Communication Technologies'},
        'it': {'code': '06', 'name': 'Information and Communication Technologies'},
        'programming': {'code': '06', 'name': 'Information and Communication Technologies'},
        
        'engineering': {'code': '07', 'name': 'Engineering, manufacturing and construction'},
        'manufacturing': {'code': '07', 'name': 'Engineering, manufacturing and construction'},
        'construction': {'code': '07', 'name': 'Engineering, manufacturing and construction'},
        
        'business': {'code': '04', 'name': 'Business, administration and law'},
        'management': {'code': '04', 'name': 'Business, administration and law'},
        'law': {'code': '04', 'name': 'Business, administration and law'},
        'finance': {'code': '04', 'name': 'Business, administration and law'},
        
        'science': {'code': '05', 'name': 'Natural sciences, mathematics and statistics'},
        'mathematics': {'code': '05', 'name': 'Natural sciences, mathematics and statistics'},
        'physics': {'code': '05', 'name': 'Natural sciences, mathematics and statistics'},
        'chemistry': {'code': '05', 'name': 'Natural sciences, mathematics and statistics'},
        'biology': {'code': '05', 'name': 'Natural sciences, mathematics and statistics'},
        
        'education': {'code': '01', 'name': 'Education'},
        'teaching': {'code': '01', 'name': 'Education'},
        
        'art': {'code': '02', 'name': 'Arts and humanities'},
        'design': {'code': '02', 'name': 'Arts and humanities'},
        'humanities': {'code': '02', 'name': 'Arts and humanities'},
        'language': {'code': '02', 'name': 'Arts and humanities'},
        
        'health': {'code': '09', 'name': 'Health and welfare'},
        'medicine': {'code': '09', 'name': 'Health and welfare'},
        'nursing': {'code': '09', 'name': 'Health and welfare'},
        'welfare': {'code': '09', 'name': 'Health and welfare'},
    }
    
    # Check for matches in general and specialized areas
    selected_code = '00'  # Default to generic programmes
    selected_name = 'Generic programmes and qualifications'
    confidence = 0.5
    
    # Check for keywords in the general area
    for keyword, mapping in isced_mapping.items():
        if keyword in general_area:
            selected_code = mapping['code']
            selected_name = mapping['name']
            confidence = random.uniform(0.7, 0.9)
            break
    
    # Check specialized area for stronger matching
    for keyword, mapping in isced_mapping.items():
        if keyword in specialized_area:
            selected_code = mapping['code']
            selected_name = mapping['name']
            confidence = random.uniform(0.85, 0.98)
            break
    
    return jsonify({
        "result": {
            "broad_code": selected_code,
            "broad_name": selected_name,
            "confidence": confidence
        },
        "confidence_score": confidence
    })

@app.route('/extract-skills', methods=['POST'])
def extract_skills():
    # In a real implementation, this would extract skills from document text
    # For this placeholder, we'll just return mock skills
    
    # Simulate processing time
    time.sleep(1.5)
    
    # Generate mock skills
    common_skills = [
        "Communication", "Leadership", "Project Management", "Problem Solving",
        "Team Collaboration", "Critical Thinking", "Time Management"
    ]
    
    tech_skills = [
        "Python", "JavaScript", "SQL", "Data Analysis", "Machine Learning",
        "Cloud Computing", "RESTful APIs", "Docker", "Kubernetes", "DevOps",
        "React", "Node.js", "Go", "Java", "C++", "CI/CD"
    ]
    
    domain_skills = [
        "Financial Analysis", "Market Research", "Product Development",
        "Quality Assurance", "Regulatory Compliance", "Digital Marketing",
        "Customer Relationship Management", "Supply Chain Management"
    ]
    
    # Randomly select skills from each category
    num_common = random.randint(2, 4)
    num_tech = random.randint(3, 6)
    num_domain = random.randint(2, 3)
    
    selected_skills = (
        random.sample(common_skills, num_common) +
        random.sample(tech_skills, num_tech) +
        random.sample(domain_skills, num_domain)
    )
    
    return jsonify({
        "result": selected_skills,
        "confidence_score": random.uniform(0.75, 0.95)
    })

if __name__ == '__main__':
    print(f"Starting AI service placeholder on port {PORT}")
    app.run(host='0.0.0.0', port=PORT)