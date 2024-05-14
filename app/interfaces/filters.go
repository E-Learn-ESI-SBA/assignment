package interfaces

type AssignmentFilter struct {
	Module		*string 		`json:"module_id,omitempty"`
	Teacher		*string			`json:"teacher_id,omitempty"`
}


func (a *AssignmentFilter) newAssignmentFilter(module string, teacher string){
	a.Module = &module
	a.Teacher = &teacher
}