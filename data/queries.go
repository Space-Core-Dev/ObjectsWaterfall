package data

const CreateTable = `CREATE TABLE IF NOT EXISTS %s (
id INTEGER PRIMARY KEY,
json_data TEXT
)`

const CreateWorkerSettingsTable = `CREATE TABLE IF NOT EXISTS Worker_Settings (
id INTEGER PRIMARY KEY,
Table_Name TEXT UNIQUE,
Timer INTEGER,
Request_Delay INTEGER,
Random BOOLEAN,
Writes_Number_To_Send INTEGER,
Total_To_Send INTEGER,
Stop_When_Table_Ends BOOLEAN
)`

const InsertData = `INSERT INTO %s (json_data) values (?)`

const GetJson = `SELECT json_data FROM %s LIMIT ?, ?`

const Count = `Select count () from %s`

const Tables = `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`

const InsertWorkerSettings = `INSERT INTO Worker_Settings (
		Worker_id
		Table_name, 
		Timer, 
		Request_Delay, 
		Random, 
		Writes_Number_To_Send, 
		Total_To_Send, 
		Stop_When_Table_Ends
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

const GetWorkerSettings = `SELECT (Table_Name, 
		Timer, 
		Request_Delay, 
		Random, 
		Writes_Number_To_Send, 
		Total_To_Send, 
		Stop_When_Table_Ends) 
	FROM Worker_Settings 
	WHERE Table_Name = ?`
