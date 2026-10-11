package go_resume_schema

import (
	"os"
	"testing"
)

func TestDecodeResumeJSON(t *testing.T) {
	data, err := os.ReadFile("resume-schema/sample.resume.json")
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	resume, err := DecodeResume(data)
	if err != nil {
		t.Fatalf("DecodeResume: %v", err)
	}
	if resume.Basics == nil || resume.Basics.Name == nil || *resume.Basics.Name != "Richard Hendriks" {
		t.Fatalf("basics: %+v", resume.Basics)
	}
	if resume.Meta.Version == "" {
		t.Fatal("missing meta.version")
	}
}

func TestDecodeResumeYAMLUnquotedDates(t *testing.T) {
	const raw = `
meta:
  version: "1.0.0"
basics:
  name: Ada Lovelace
work:
  - employer: Analytical Engine
    position: Programmer
    startDate: 2020-05-05
    endDate: "2021-01"
`
	resume, err := DecodeResume([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeResume: %v", err)
	}
	if resume.Basics == nil || resume.Basics.Name == nil || *resume.Basics.Name != "Ada Lovelace" {
		t.Fatalf("basics: %+v", resume.Basics)
	}
	if len(resume.Work) != 1 || resume.Work[0].StartDate == nil || string(*resume.Work[0].StartDate) != "2020-05-05" {
		t.Fatalf("work: %+v", resume.Work)
	}
	if resume.Work[0].EndDate == nil || string(*resume.Work[0].EndDate) != "2021-01" {
		t.Fatalf("endDate: %+v", resume.Work[0].EndDate)
	}
}

func TestDecodeResumeEmpty(t *testing.T) {
	_, err := DecodeResume([]byte("   "))
	if err == nil {
		t.Fatal("expected error")
	}
}
