package interfaces


type AssignmentFilter struct {
	ModuleId		*string 		`json:"module_id,omitempty"`
	TeacherId		*string			`json:"teacher_id,omitempty"`
}


func (a *AssignmentFilter) newAssignmentFilter(moduleId string, teacherId string){
	a.ModuleId = &moduleId
	a.TeacherId = &teacherId
}