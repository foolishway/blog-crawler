package models

type Duty struct {
	Id                           int
	OnDuty                       int
	Name, PhoneNum, EmployeesNum string
}

func GetAllDuty() []Duty {
	duty := make([]Duty, 0)
	db.Table("duty").Find(&duty)
	return duty
}

func GetDutyById(DutyId int) Duty {
	duty := Duty{}
	db.Table("duty").Where("id = ?", DutyId).Find(&duty)
	return duty
}

func UpdateDutyById(d Duty) error {
	duty := Duty{}
	db.Table("duty").Where("id = ?", d.Id).First(&duty)
	duty.PhoneNum = d.PhoneNum
	duty.EmployeesNum = d.EmployeesNum
	duty.Name = d.Name
	return db.Table("duty").Save(&duty).Error
}

func DelDutyById(d Duty) error {
	return db.Table("duty").Delete(d).Error
}

func InsertDuty(duty Duty) error {
	return db.Table("duty").Create(&duty).Error
}
