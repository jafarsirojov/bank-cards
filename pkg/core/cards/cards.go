package cards

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"time"
)

const initNumberCard = 2021600000000000

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	_, err := service.pool.Exec(context.Background(), cardsDDL)
	log.Print(err)
	numberBankCount := strconv.Itoa(initNumberCard)
	_, err = service.pool.Exec(context.Background(), initialInsertCard, numberBankCount)
	log.Print("Has Bank Count")

}

func (service *Service) All() (models []Cards, err error) {
	rows, err := service.pool.Query(context.Background(), selectAllCards)
	if err != nil {
		return nil, fmt.Errorf("can't get cards from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := Cards{}
		err = rows.Scan(
			&user.Id,
			&user.Number,
			&user.Name,
			&user.Balance,
			&user.OwnerID,
		)
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
	err = service.pool.QueryRow(context.Background(), selectCardById, id).Scan(
		&cards.Id,
		&cards.Number,
		&cards.Name,
		&cards.Balance,
		&cards.OwnerID,
	)
	if err != nil {
		log.Printf("can't select cards by id: %d", err)
		return nil, err
	}
	model = append(model, cards)
	return model, nil
}

func (service *Service) ByIdUserCard(idCard int, ownerID int) (model []Cards, err error) {
	cards := Cards{}
	err = service.pool.QueryRow(context.Background(), selectCardsByIdAndUserId, idCard, ownerID).Scan(
		&cards.Id,
		&cards.Number,
		&cards.Name,
		&cards.Balance,
		&cards.OwnerID,
	)
	if err != nil {
		log.Printf("can't select cards by idCard: %d", err)
		log.Print("client not owner card & note ")
		return nil, err
	}
	model = append(model, cards)
	return model, nil
}

func (service *Service) ViewCardsByOwnerId(id int) (models []Cards, err error) {
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


	rows, err := service.pool.Query(context.Background(), selectCardsByOwnerId, id)
	if err != nil {
		return nil, fmt.Errorf("can't get cards from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := Cards{}
		err = rows.Scan(
			&user.Id,
			&user.Number,
			&user.Name,
			&user.Balance,
			&user.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("can't get cards from db: %w", err)
		}
		models = append(models, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get cards from db: %w", err)
	}
	return models, nil
	//err = tx.QueryRow(context.Background(), selectCardsByOwnerId, id).Scan(
	//	&user.Id,
	//	&user.Number,
	//	&user.Name,
	//	&user.Balance,
	//	&user.OwnerID,
	//)
	//
	//model = append(model, user)
	//return model, nil
}

func (service *Service) AddCard(model Cards) (err error) {
	selectDescIdFromCard := 0
	var numberCard int
	err = service.pool.QueryRow(context.Background(), selectIdLimit1).Scan(&selectDescIdFromCard)
	if err != nil {
		log.Print("select id cards desc limit 1")
		return err
	}
	numberCard = selectDescIdFromCard + 1 + initNumberCard
	model.Number = strconv.Itoa(numberCard)
	_, err = service.pool.Exec(context.Background(), insertCard, model.Number, model.Name, model.Balance, model.OwnerID)
	if err != nil {
		log.Print("can't exec insert ")
		return err
	}
	return nil
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

func (service *Service) UserBlockCardById(userID int, model ModelBlockCard) (err error) {
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
	_, err = tx.Exec(context.Background(), `
	UPDATE cards 
	SET blocked=TRUE 
	WHERE blocked = FALSE and id=$1 and owner_id=$2`,
		model.Id,
		userID,
	)
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

func (service *Service) TransferMoneyCardToCard(ownerID int, model ModelTransferMoneyCardToCard) (err error, modelHistoryTranSender, modelHistoryTranRecipient ModelOperationsLog) {
	tx, err := service.pool.Begin(context.Background())
	if err != nil {
		log.Print("can't begin tx")
		return err, modelHistoryIsNil, modelHistoryIsNil
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
	var numberCardSender string
	var senderID int64
	var recipientID int64
	err = tx.QueryRow(context.Background(),
		`SELECT  balance, number, owner_id
		FROM cards 
		WHERE blocked = FALSE and id = $1 and owner_id = $2;`, model.IdCardSender,ownerID).Scan(&oldBalanceCardSender, &numberCardSender, &senderID)
	if err != nil {
	log.Printf("can't select sender to db: %s",err)
		return err, modelHistoryIsNil, modelHistoryIsNil
	}
	err = tx.QueryRow(context.Background(), `
	SELECT  balance, owner_id 
	FROM cards 
	WHERE blocked = FALSE and number = $1;`,
		model.NumberCardRecipient).Scan(&oldBalanceCardRecipient, &recipientID)
	if err != nil {
		log.Printf("can't select recipient to db: %s",err)
		return err, modelHistoryIsNil, modelHistoryIsNil
	}
	var newBalanceCardSender = oldBalanceCardSender - model.Count
	var newBalanceCardRecipient = oldBalanceCardRecipient + model.Count
	_, err = tx.Exec(context.Background(),
		`UPDATE cards 
		SET balance=$1 
		WHERE blocked = FALSE and id = $2 and owner_id = $3`,
		newBalanceCardSender,
		model.IdCardSender /*model.IdCardSender*/,
		ownerID,
	)
	if err != nil {
		log.Printf("can't exec update transfer money: %d", err)
		return err, modelHistoryIsNil, modelHistoryIsNil
	}
	_, err = tx.Exec(context.Background(),
		`UPDATE cards 
		SET balance=$1 
		WHERE blocked = FALSE and number = $2`,
		newBalanceCardRecipient,
		model.NumberCardRecipient,
	)
	if err != nil {
		log.Printf("can't exec update transfer money: %d", err)
		return err, modelHistoryIsNil, modelHistoryIsNil
	}
	log.Print("transfer ok")
	modelHistoryTranSender = ModelOperationsLog{
		Id:              0,
		Name:            "Transfer_money",
		Number:          numberCardSender,
		RecipientSender: "sender",
		Count:           model.Count,
		BalanceOld:      oldBalanceCardSender,
		BalanceNew:      newBalanceCardSender,
		Time:            time.Now().Unix(),
		OwnerID:         senderID,
	}
	modelHistoryTranRecipient = ModelOperationsLog{
		Id:              0,
		Name:            "Transfer_money",
		Number:          model.NumberCardRecipient,
		RecipientSender: "recipient",
		Count:           model.Count,
		BalanceOld:      oldBalanceCardRecipient,
		BalanceNew:      newBalanceCardRecipient,
		Time:            time.Now().Unix(),
		OwnerID:         recipientID,
	}

	log.Print("save model history to dto")
	return
}

type Cards struct {
	Id      int    `json:"id"`
	Number  string `json:"number"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
	OwnerID int64  `json:"owner_id"`
}
type ModelTransferMoneyCardToCard struct {
	IdCardSender        int    `json:"id_card_sender"`
	NumberCardRecipient string `json:"number_card_recipient"`
	Count               int64  `json:"count"`
}

type ModelBlockCard struct {
	Id     int    `json:"id"`
	Number string `json:"number"`
}

type ModelOperationsLog struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	RecipientSender string `json:"recipientsender"`
	Count           int64  `json:"count"`
	BalanceOld      int64  `json:"balanceold"`
	BalanceNew      int64  `json:"balancenew"`
	Time            int64  `json:"time"`
	OwnerID         int64  `json:"ownerid"`
}

var modelHistoryIsNil = ModelOperationsLog{
	Id:              0,
	Name:            "",
	Number:          "",
	RecipientSender: "",
	Count:           0,
	BalanceOld:      0,
	BalanceNew:      0,
	Time:            0,
	OwnerID:         0,
}

type TokenResponse struct {
	Token string `json:"token"`
}
