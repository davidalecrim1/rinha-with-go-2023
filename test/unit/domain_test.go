package domain

import (
	"go-rinha-de-backend-2023/internal/domain"
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
			person, err := domain.NewPerson(tt.fields.nickname, tt.fields.name, tt.fields.dob, tt.fields.stack)

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
			p := &domain.Person{
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

func BenchmarkPerson_Validate(b *testing.B) {
	person := domain.Person{
		ID:       "1",
		Nickname: "john",
		Name:     "John Doe",
		Dob:      "1990-01-01",
		Stack:    []string{"Go", "Python"},
	}

	for i := 0; i < b.N; i++ {
		person.Validate()
	}

}

func BenchmarkPerson_NewPerson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = domain.NewPerson(
			"john",
			"John Doe",
			"1990-01-01",
			[]string{"Go", "Python"},
		)
	}
}

func BenchmarkPeople_WithAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var people []domain.Person
		for j := 0; j < 50; j++ {
			people = append(people, domain.Person{ID: "34aafceb-325d-48a6-b16b-9de90d3578a8"})
		}
	}
}

// This option is 76x times faster the allocating on demand as above
// Trade off for more memory usage and faster execution
func BenchmarkPeople_WithPreallocatedSliceWith50Valid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		people := make([]domain.Person, 50)
		validCount := 0
		for j := 0; j < 50; j++ {
			people[validCount] = domain.Person{ID: "34aafceb-325d-48a6-b16b-9de90d3578a8"}
			validCount++
		}
		people = people[:validCount]
	}
}

// This is only 2x or 3x times faster then allocating on demand as the above one
// Trade off for more memory usage and faster execution
func BenchmarkPeople_WithPreallocatedSliceAndOnly15Created(b *testing.B) {
	var people []domain.Person
	for i := 0; i < b.N; i++ {
		people = make([]domain.Person, 50)
		validCount := 0
		for j := 0; j < 15; j++ {
			people[validCount] = domain.Person{ID: "34aafceb-325d-48a6-b16b-9de90d3578a8"}
			validCount++
		}
		people = people[:validCount]
	}
}

// This is 1x time faster then allocating the slice on demand
// Trade off for more memory usage and faster execution
func BenchmarkPeople_WithPreallocatedSliceAndOnly15CreatedAndValidated(b *testing.B) {
	var people []domain.Person
	for i := 0; i < b.N; i++ {
		people = make([]domain.Person, 50)
		for j := 0; j < 15; j++ {
			people[j] = domain.Person{ID: "34aafceb-325d-48a6-b16b-9de90d3578a8"}
			j++
		}

		validCount := 0
		for _, person := range people {
			if person.ID == "" {
				people[validCount] = person
				validCount++
			}
		}
		people = people[:validCount] // Trim the slice to remove invalid entries
	}
}
