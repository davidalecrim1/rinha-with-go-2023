package domain

import (
	"testing"
)

func TestPerson_NewPerson(t *testing.T) {
	type fields struct {
		nickname string
		name     string
		dob      string
		stack    []string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid person",
			fields: fields{
				nickname: "john",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			person, err := NewPerson(tt.fields.nickname, tt.fields.name, tt.fields.dob, tt.fields.stack)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(person.ID) != 36 {
				t.Errorf("NewPerson() error = invalid uuid, wantErr %v", tt.wantErr)
				return
			}
		})
	}
}

func TestPerson_Validate(t *testing.T) {
	type fields struct {
		nickname string
		name     string
		dob      string
		stack    []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid person",
			fields: fields{
				nickname: "john",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: false,
		},
		{
			name: "Invalid long nickname",
			fields: fields{
				nickname: "this is a really long nickname with 49 characters",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: true,
		},
		{
			name: "Invalid long name",
			fields: fields{
				nickname: "john",
				name:     "John Doe with a really long name that is over a 113 characters that why it is invalid, also this needs to be long",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: true,
		},
		{
			name: "Invalid empty nickname",
			fields: fields{
				nickname: "",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: true,
		},
		{
			name: "Invalid empty name",
			fields: fields{
				nickname: "johndoe",
				name:     "",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python"},
			},
			wantErr: true,
		},
		{
			name: "Invalid DOB",
			fields: fields{
				nickname: "johndoe",
				name:     "John Doe",
				dob:      "this is a invalid date",
				stack:    []string{"Go", "Python"},
			},
			wantErr: true,
		},
		{
			name: "Valid empty stack",
			fields: fields{
				nickname: "johndoe",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{},
			},
			wantErr: false,
		},
		{
			name: "Invalid stack item",
			fields: fields{
				nickname: "johndoe",
				name:     "John Doe",
				dob:      "1990-01-01",
				stack:    []string{"Go", "Python", "This is a really long stack name with 51 characters"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				Nickname: tt.fields.nickname,
				Name:     tt.fields.name,
				Dob:      tt.fields.dob,
				Stack:    tt.fields.stack,
			}

			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Person.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
