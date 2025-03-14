Dear Clade, kindly take note of the below instructions:
I have amended both the CreateExpertRequest and Expert structs in types.go. Both structs should be reflective of each other, but for now will be satisfied with a watererd down CreateExpertRequest
below is the edited request struct. 
// NOTE: Dear Claude, kindly amend below CreateExpertRequest struct according to the instructions
type CreateExpertRequest struct {
	Name           string   `json:"name"`
	Affiliation    string   `json:"affiliation"`
	PrimaryContact string   `json:"primaryContact"`
	ContactType    string   `json:"contactType"` // "email" or "phone"
	Skills         []string `json:"skills"`
	Role           string   `json:"role"`
	EmploymentType string   `json:"employmentType"`
	GeneralArea    string   `json:"generalArea"`
	CVPath         string   `json:"cvPath"`
	Biography      string   `json:"biography"`
	IsBahraini     bool     `json:"isBahraini"`
	Availability   string   `json:"availability"` // Availability means still active or not. A way to mark as inactive. to be removed
}
 I've also reflected the Biography field to Expert struct.

# Instructions:
- Review the amended structs in types.go and reflect the necessary changes where relevant (e.g. struct initializations, relevant methods, database schema, handler request parsing...etc)
- Amend the migration files, ensuring that the code will have no problem creating the db and populating it with the provided csv.
- Pay attention and do not take an action you are not sure about it's full context.
