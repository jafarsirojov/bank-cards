package cards

const cardsDDL = `
CREATE TABLE IF NOT EXISTS cards (
	id BIGSERIAL,
	number TEXT NOT NULL UNIQUE,
    name TEXT,
   	balance INTEGER NOT NULL,
	owner_id INTEGER NOT NULL,
    blocked BOOLEAN DEFAULT FALSE
);`

const initialInsertCard = `INSERT INTO cards(id, number, name, balance, owner_id) VALUES (0, $1, 'Bank Count',  0, 0);`

const selectAllCards  =`SELECT id, number, name, balance, owner_id FROM cards WHERE blocked = FALSE;`

const selectCardById  =`
	SELECT id, number, name, balance, owner_id 
	FROM cards 
	WHERE blocked = FALSE and id = $1;`
const selectCardsByIdAndUserId = `
	SELECT idCard, number, name, balance, owner_id 
	FROM cards 
	WHERE blocked = FALSE and idCard = $1 and owner_id = $2;`
const selectCardsByOwnerId  = `
	SELECT id, number, name, balance, owner_id 
	FROM cards 
	WHERE blocked = FALSE and owner_id = $1;`
const selectIdLimit1  =`SELECT id FROM cards ORDER BY id DESC LIMIT 1`

const insertCard  = `INSERT INTO cards(number, name, balance, owner_id) VALUES ($1, $2, $3, $4)`