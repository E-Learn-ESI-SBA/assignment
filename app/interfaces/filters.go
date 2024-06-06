package interfaces


type AssignmentFilter struct {
	ModuleId		*string 		`json:"module_id,omitempty"`
	TeacherId		*string			`json:"teacher_id,omitempty"`
	Year		    *string			`json:"year,omitempty"`
}


func (a *AssignmentFilter) newAssignmentFilter(moduleId string, teacherId string, year string){
	a.ModuleId = &moduleId
	a.TeacherId = &teacherId
	a.Year = &year
}