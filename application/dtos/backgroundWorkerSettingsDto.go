package dtos

type BackgroundWorkerSettingsDto struct {
	TableName          string  `json:"tableName"`
	Timer              float64 `json:"timer"`
	RequestDelay       int     `json:"requestDellay"`
	Random             bool    `json:"random"`
	WritesNumberToSend int     `json:"writesNumberToSend"`
	TotalToSend        int64   `json:"totalToSend"`
	StopWhenTableEnds  bool    `json:"stopWhenTableEnds"`
}
