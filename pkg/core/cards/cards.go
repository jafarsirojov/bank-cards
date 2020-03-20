package cards

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
)

const initNumberCard = 2021600000000000

type Service struct {
	pool *pgxpool.Pool // dependency
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	//conn, err := service.pool.Acquire(context.Background())
	//defer conn.Release()
	// TODO: panic

	_, err := service.pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS cards (
	id BIGSERIAL,
	number TEXT NOT NULL UNIQUE,
    name TEXT,
   	balance INTEGER NOT NULL,
	owner_id INTEGER NOT NULL,
    blocked BOOLEAN DEFAULT FALSE
);
`)
	log.Print(err)
	numberBankCount := strconv.Itoa(initNumberCard)
	_, err = service.pool.Exec(context.Background(), `
INSERT INTO cards(id, number, name, balance, owner_id) VALUES (0, $1, 'Bank Count',  0, 0);
`, numberBankCount)
	log.Print("Has Bank Count")

}

// created timestamp
// modified timestamp

// CRUD
func (service *Service) All() (models []Cards, err error) {
	rows, err := service.pool.Query(context.Background(), `
SELECT id, number, name, balance, owner_id FROM cards WHERE blocked = FALSE;
`)
	if err != nil {
		return nil, fmt.Errorf("can't get cards from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := Cards{}
		err = rows.Scan(&user.Id, &user.Number, &user.Name, &user.Balance, &user.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't get cards from db: %w", err)
		}
		models = append(models, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get cards from db: %w", err)
	}
	return models, nil
}

func (service *Service) ById(id int) (model []Cards, err error) {
	cards := Cards{}
	_ = service.pool.QueryRow(context.Background(), `
SELECT id, number, name, balance, owner_id FROM cards WHERE blocked = FALSE and id = $1;
`, id).Scan(&cards.Id, &cards.Number, &cards.Name, &cards.Balance, &cards.OwnerID)
	model = append(model, cards)
	return model, nil
}

func (service *Service) ViewCardsByOwnerId(id int) (model []Cards, err error) {
	user := Cards{}
	tx, err := service.pool.Begin(context.Background())
	if err != nil {
		log.Printf("can't begin view cards by id owner: %d", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}
		err = tx.Commit(context.Background())
	}()
	err = tx.QueryRow(context.Background(), `
SELECT id, number, name, balance, owner_id FROM cards WHERE blocked = FALSE and owner_id = $1;
`, id).Scan(&user.Id, &user.Number, &user.Name, &user.Balance, &user.OwnerID)

	model = append(model, user)
	return model, nil
}

func (service *Service) AddCard(model Cards) {
	selectDescIdFromCard := 0
	var numberCard int
	err := service.pool.QueryRow(context.Background(), `SELECT id FROM cards ORDER BY id DESC LIMIT 1`).Scan(&selectDescIdFromCard)
	if err != nil {
		log.Print()
	}
	numberCard = selectDescIdFromCard + 1 + initNumberCard
	model.Number = strconv.Itoa(numberCard)
	_, err = service.pool.Exec(context.Background(), `INSERT INTO cards(number, name, balance, owner_id) VALUES ($1, $2, $3, $4)`, model.Number, model.Name, model.Balance, model.OwnerID)
	if err != nil {
		log.Print("can't exec insert ")
	}
}

func (service *Service) BlockById(id int) (err error) {
	tx, err := service.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}
		err = tx.Commit(context.Background())
	}()
	_, err = tx.Exec(context.Background(), `UPDATE cards SET blocked=TRUE WHERE blocked = FALSE and id=$1`, id)
	if err != nil {
		log.Print("can't exec update blocked ", err)
		return err
	}
	return nil
}

func (service *Service) UnBlockedById(id int) (err error) {
	tx, err := service.pool.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}
		err = tx.Commit(context.Background())
	}()
	_, err = tx.Exec(context.Background(), `UPDATE cards SET blocked=FALSE WHERE blocked = TRUE and id=$1`, id)
	if err != nil {
		log.Print("can't exec unblocked", err)
		return err
	}
	return nil
}

func (service *Service) TransferMoneyCardToCard(idCardSender int, numberCardRecipient string, count int64) (err error) {
	tx, err := service.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}
		err = tx.Commit(context.Background())
	}()
	var oldBalanceCardSender int64
	var oldBalanceCardRecipient int64
	err = tx.QueryRow(context.Background(), `SELECT  balance FROM cards WHERE blocked = FALSE and id = $1;`, idCardSender).Scan(&oldBalanceCardSender)
	if err != nil {
		log.Print(err)
		return err
	}
	err = tx.QueryRow(context.Background(), `SELECT  balance FROM cards WHERE blocked = FALSE and number = $1;`, numberCardRecipient).Scan(&oldBalanceCardRecipient)
	if err != nil {
		log.Print(err)
		return err
	}
	var newBalanceCardSender = oldBalanceCardSender - count
	var newBalanceCardRecipient = oldBalanceCardRecipient + count
	_, err = tx.Exec(context.Background(), `UPDATE cards SET balance=$1 WHERE blocked = FALSE and id = $2`, newBalanceCardSender, idCardSender)
	if err != nil {
		log.Printf("can't exec update transfer money: %d", err)
		return err
	}
	_, err = tx.Exec(context.Background(), `UPDATE cards SET balance=$1 WHERE blocked = FALSE and number = $2`, newBalanceCardRecipient, numberCardRecipient)
	if err != nil {
		log.Printf("can't exec update transfer money: %d", err)
		return err
	}
	log.Print("transfer ok")
	return nil
}

type Cards struct {
	Id      int
	Number  string
	Name    string
	Balance int64
	OwnerID int64
}
type ModelTransferMoneyCardToCard struct {
	IdCardSender        int
	NumberCardRecipient string
	Count               int64
}
