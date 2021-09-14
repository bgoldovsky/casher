package operations

import (
	"database/sql"
	"testing"

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
	_, err := s.db.Query("delete from operations")
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
	err := s.store.Create(&models.Operation{
		UserID:  10000000,
		Subject: "Таверна Fish & Chips",
		Amount:  150000,
		Type:    models.Withdraw,
		Message: "Отметил приезд",
	})
	if err != nil {
		s.T().Fatal(err)
	}

	res, err := s.db.Query(`select count(*) from operations where user_id=10000000 and subject='Таверна Fish & Chips' and amount=150000`)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *storeSuite) TestGet() {
	_, err := s.db.Query(`insert into operations (user_id, subject, amount, type, message) values(10000000,'Таверна Fish & Chips', 150000,2, 'Отметил приезд')`)
	if err != nil {
		s.T().Fatal(err)
	}

	paginator, err := s.store.Get(10000000, 1, 1)
	if err != nil {
		s.T().Fatal(err)
	}

	count := len(paginator.Operations)
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}

	var exp = &models.Operation{
		UserID:  10000000,
		Subject: "Таверна Fish & Chips",
		Amount:  150000,
		Type:    models.Withdraw,
		Message: "Отметил приезд",
	}

	act := paginator.Operations[0]

	if act.UserID != exp.UserID {
		s.T().Errorf("expected %v, got %v", exp.UserID, act.UserID)
	}

	if act.Subject != exp.Subject {
		s.T().Errorf("expected %v, got %v", exp.Subject, act.Subject)
	}

	if act.Amount != exp.Amount {
		s.T().Errorf("expected %v, got %v", exp.Amount, act.Amount)
	}

	if act.Type != exp.Type {
		s.T().Errorf("expected %v, got %v", exp.Type, act.Type)
	}

	if act.Message != exp.Message {
		s.T().Errorf("expected %v, got %v", exp.Message, act.Message)
	}
}
