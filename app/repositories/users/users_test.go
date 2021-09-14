package users

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bgoldovsky/casher/app/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type storeSuite struct {
	suite.Suite
	store *repository
	db    *sql.DB
}

func (s *storeSuite) SetupSuite() {
	connString := "dbname=casher sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &repository{db: db}
}

func (s *storeSuite) SetupTest() {
	_, err := s.db.Query("delete from operations; delete from users;")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *storeSuite) TearDownSuite() {
	_ = s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(storeSuite)
	suite.Run(t, s)
}

func (s *storeSuite) TestCreate() {
	var (
		_, err = s.store.Create(&models.User{
			Login:    "jondoe",
			Password: "qwerty",
			Name:     "Jon Doe",
			Birth:    time.Now(),
		})
	)
	if err != nil {
		s.T().Fatal(err)
	}

	res, err := s.db.Query(`select count(*) from users where login='jondoe' and name='Jon Doe'`)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err = res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *storeSuite) TestGet() {
	_, err := s.db.Query(`insert into users (id, login, password, name, birth) values(10000000, 'jondoe','qwerty', 'Jon Doe', now())`)
	if err != nil {
		s.T().Fatal(err)
	}

	act, err := s.store.Get(10000000)
	if err != nil {
		s.T().Fatal(err)
	}

	var exp = &models.User{
		ID:       10000000,
		Login:    "jondoe",
		Password: "qwerty",
		Name:     "Jon Doe",
		Birth:    time.Now(),
	}

	if act.ID != exp.ID {
		s.T().Errorf("expected %v, got %v", exp.ID, act.ID)
	}

	if act.Login != exp.Login {
		s.T().Errorf("expected %v, got %v", exp.Login, act.Login)
	}

	if act.Password != exp.Password {
		s.T().Errorf("expected %v, got %v", exp.Password, act.Password)
	}

	if act.Name != exp.Name {
		s.T().Errorf("expected %v, got %v", exp.Name, act.Name)
	}
}

func (s *storeSuite) TestLogin() {
	_, err := s.db.Query(`insert into users (id, login, password, name, birth) values(10000000, 'jondoe','qwerty', 'Jon Doe', now())`)
	if err != nil {
		s.T().Fatal(err)
	}

	act, err := s.store.Auth("jondoe")
	if err != nil {
		s.T().Fatal(err)
	}

	var exp = &models.User{
		ID:       10000000,
		Login:    "jondoe",
		Password: "qwerty",
		Name:     "Jon Doe",
		Birth:    time.Now(),
	}

	if act.ID != exp.ID {
		s.T().Errorf("expected %v, got %v", exp.ID, act.ID)
	}

	if act.Login != exp.Login {
		s.T().Errorf("expected %v, got %v", exp.Login, act.Login)
	}

	if act.Password != exp.Password {
		s.T().Errorf("expected %v, got %v", exp.Password, act.Password)
	}

	if act.Name != exp.Name {
		s.T().Errorf("expected %v, got %v", exp.Name, act.Name)
	}
}
