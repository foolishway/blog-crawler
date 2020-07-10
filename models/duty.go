package models

type Duty struct {
	Id                                    int
	Name, PhoneNum, EmployeesNum, DutyDay string
}

func GetAllDuty() []Duty {
	//defer db.Close()
	duty := make([]Duty, 0)
	db.Table("duty").Find(&duty)
	//fmt.Println(articles)
	return duty
}

func GetDutyById(DutyId int) Duty {
	return Duty{}
}

func UpdateDutyById(DutyId int) error {
	return nil
}

func DelDutyById(DutyId int) error {
	return nil
}

func InsertDuty(duty Duty) error {
	return nil
}
