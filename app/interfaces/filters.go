package interfaces

type AssignmentFilter struct {
	Module		*string 		`json:"module_id,omitempty"`
	Teacher		*int		`json:"teacher_id,omitempty"`
}


func (a *AssignmentFilter) newAssignmentFilter(module string, teacher int){
	a.Module = &module
	a.Teacher = &teacher
}