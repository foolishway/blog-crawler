package models

import (
	"log"
	"testing"
)

func TestInsertDuty(t *testing.T) {
	duty := Duty{Name: "刘传明", EmployeesNum: "161265", PhoneNum: "18210094182"}
	err := InsertDuty(duty)
	if err != nil {
		log.Fatalf("Inser error %s", err)
	}
}

func TestGetDutyById(t *testing.T) {
	id := 2
	duty := GetDutyById(id)
	if duty.Name != "流川枫" {
		t.Fatalf("id:%d name want 流川枫 but %s", id, duty.Name)
	}
}

func TestUpdateDutyById(t *testing.T) {
	duty := Duty{Id: 2, Name: "流川枫", EmployeesNum: "106316", PhoneNum: "13523123213"}
	err := UpdateDuty(duty)
	if err != nil {
		t.Fatalf("update error %v", err)
	}
	duty = GetDutyById(2)
	if duty.Name != "流川枫" {
		t.Fatalf("update error id:%d name want 刘传明 but %s", 2, duty.Name)
	}
}

func TestDelDutyById(t *testing.T) {
	id := 2
	duty := GetDutyById(id)
	err := DelDutyById(duty)
	if err != nil {
		t.Fatalf("delete error %v", err)
	}
	duty = GetDutyById(id)
	if duty.Name != "" {
		t.Fatalf("delete error want name nil, but %s", duty.Name)
	}
}

func TestGetNextDuty(t *testing.T) {
	want := "张威,刘传明,"
	nextDuty, err := GetNextDuty()
	if err != nil {
		t.Fatalf("GetNextDuty error %v", err)
	}

	if len(nextDuty) != 2 {
		t.Fatalf("GetNextDuty error %v, next duty != 2, duty info %v", err, nextDuty)
	}
	names := ""
	for _, nd := range nextDuty {
		names += nd.Name + ","
	}
	if names != want {
		t.Fatalf("GetNextDuty error %v, want %s but %s", err, want, names)
	}
}
