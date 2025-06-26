package data

const CreateTable = `CREATE TABLE IF NOT EXISTS %s (
id INTEGER PRIMARY KEY,
json_data TEXT
)`

const InsertData = `INSERT INTO %s (json_data) values (?)`

const GetJson = `SELECT json_data FROM %s LIMIT ?, ?`

const Count = `Select count () from %s`
