## NOTES:
- Plan frontend dev of each feature one by one i.e database browing, expert_creation...etc
- Create a mermaid.js for workflows
- Create wireframe for service frontend
- Remove complexity: is there a need for context package?
- Enhance script to sumulate workflows
  - create db entries of all user types
  - simulate expert_creation ==> Apprval/ Rejection ==> expert creation
- Future: Implement email notification service:
  - New expert request (mailto:admin)
  - approved expert request (mailto:user)
  - rejected expert request (mailto:user)
  - updated expert request (mailto:admin)
- Ensure all nullable values are handled properly in the backend

## Types:
A. Users: super, admin, planner (replaced scheduler), user
B. Documents: cv, approval_document
C. Expert_Request: pending, rejected, approved
D. Application: Qualification Placement (QP), Institutional Listing (IL)
E. Stats: annual growth (year over year, since last year), nationality representation (Bahrainin to no none-bahraini), engagement (number of engagements for each expert by type:QP and/or IL), 

## Functions and Features:
- Facilitate adding experts to database
- Facilitate browsing expert database w/ filtering and sorting
- Facilitate planning phases
- Provide statistics on the database

## Workflows
1. Expert Creation:
  a. User creates expert_request
    - Fill form
    - Attach documents
  b. Admin received expert_request:
    1. Approve request ==> create expert entry in experts table
    2. Reject request ==> user receives rejected request for amendment (c)
    3. Update request ==> approve request ==> create expert entry in experts table


2. Phase Planninng:
  a. Admin creates Phase
  b. Amdin creates Applications under Phase
    - Application Types: Qualification Placement (QP) and Institutional Listing (IL)
  c. Admin assings Application/s (singl or batch) to Planner
  d. Planner receives Application/s:
    d1. Planner assigns Expert-1 and Expert-2 to each application
    d2. Planner submits Application/s to admin for review
  e. Admin reviews planner suggestions:
    e1. Approved plan ==> create engagements
    e2. Reject plan ==> 
      - ==> return to user for amnedments
      - ==> amend and accept

