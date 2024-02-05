package repository

import (
	"clean/pkg/utils/models"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetUserDetails(t *testing.T) {
	tests := []struct {
		name    string
		args    int
		stub    func(mockSQL sqlmock.Sqlmock)
		want    models.UsersProfileDetails
		wantErr error
	}{
		{
			name: "success",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := "select users.firstname,users.lastname,users.email,users.phone from users  where users.id = ?"
				mockSQL.ExpectQuery(expectQuery).WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "email", "phone"}).AddRow("preya", "v", "preya@gmail.com", "7902689612"))
			},
			want: models.UsersProfileDetails{
				Firstname: "preya",
				Lastname:  "v",
				Email:     "preya@gmail.com",
				Phone:     "7902689612",
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := `select users.firstname,users.lastname,users.email,users.phone from users  where users.id = ?`
				mockSQL.ExpectQuery(expectQuery).WillReturnError(errors.New("error"))
			},
			want:    models.UsersProfileDetails{},
			wantErr: errors.New("could not get user details"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.stub(mockSQL)
			u := NewuserRepository(gormDB)

			result, err := u.UserDetails(tt.args)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.SignupDetail
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.SignupDetailResponse
		wantErr    error
	}{
		{
			name: "success signup user",
			args: args{
				input: models.SignupDetail{
					Firstname: "preya",
					Lastname:  "v",
					Email:     "preya@gmail.com",
					Password:  "12345",
					Phone:     "7902689612",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(firstname, lastname, email, password, phone\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, firstname, lastname, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("preya", "v", "preya@gmail.com", "12345", "7902689612").
					WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
						AddRow(1, "preya", "v", "preya@gmail.com", "7902689612"))
			},
			want: models.SignupDetailResponse{
				Id:        1,
				Firstname: "preya",
				Lastname:  "v",
				Email:     "preya@gmail.com",
				Phone:     "7902689612",
			},
			wantErr: nil,
		},
		{
			name: "error signup user",
			args: args{
				input: models.SignupDetail{},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(firstname, lastname, email, password, phone\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, firstname, lastname, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("", "", "", "", "").
					WillReturnError(errors.New("email should be unique"))
			},

			want:    models.SignupDetailResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSQL)
			u := NewuserRepository(gormDB)
			got, err := u.UserSignup(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}