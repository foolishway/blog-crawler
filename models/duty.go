package models

type Duty struct {
	Id             int
	OnDuty         int
	Name, PhoneNum string
	EmployeesNum   string `json:"employees_num"`
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

func UpdateDuty(d Duty) error {
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

func GetNextDuty() (nextDuty []Duty, err error) {
	var reset = func() {
		//update nextDuty
		if len(nextDuty) > 0 {
			ids := make([]int, 0)
			for _, nd := range nextDuty {
				ids = append(ids, nd.Id)
			}
			err = db.Table("duty").Where("id IN (?)", ids).Update("on_duty", 1).Error
			//reset preDuty
			err = db.Table("duty").Where("id NOT IN (?)", ids).Update("on_duty", 0).Error
		}
	}
	preDuty := Duty{}
	err = db.Table("duty").Where("on_duty = 1").Last(&preDuty).Error
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		} else {
			err = db.Table("duty").Where("on_duty = 0").Limit(2).Find(&nextDuty).Error
			reset()
			return
		}
	}
	err = db.Table("duty").Where("on_duty = 0 AND id > ?", preDuty.Id).Limit(2).Find(&nextDuty).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0)
	if len(nextDuty) > 0 {
		for _, d := range nextDuty {
			ids = append(ids, d.Id)
		}
	}
	//log.Printf("ids:%v", ids)
	reDuty := make([]Duty, 0)
	if len(nextDuty) < 2 {
		if len(ids) > 0 {
			err = db.Table("duty").Where("on_duty = 0 AND id NOT IN (?)", ids).Limit(2 - len(nextDuty)).Find(&reDuty).Error
		} else {
			err = db.Table("duty").Where("on_duty = 0").Limit(2 - len(nextDuty)).Find(&reDuty).Error
		}
		if err != nil {
			return nil, err
		}
		for _, d := range reDuty {
			nextDuty = append(nextDuty, d)
		}
		//log.Printf("reDuty:%#v", reDuty)
		if len(nextDuty) < 2 {
			err = db.Table("duty").Limit(2 - len(nextDuty)).Find(&reDuty).Error
			for _, d := range reDuty {
				nextDuty = append(nextDuty, d)
			}
		}
	}
	reset()
	return
}
