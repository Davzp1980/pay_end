package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type input struct {
	name     string
	password string
}

func Test_CreateAdmin(t *testing.T) {
	//isAdmin:=true
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		mock    func()
		input   input
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("test", "password", true).WillReturnRows(rows)
			},
			input: input{
				name:     "test",
				password: "password",
			},
			want: "Admin 1 created",
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("", "password", true).WillReturnRows(rows)
			},
			input: input{
				name:     "",
				password: "password",
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := CreateAdmin(db, tt.input.name, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
