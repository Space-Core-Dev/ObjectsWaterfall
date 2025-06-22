package models

type BackgroundWorkerSettings struct {
	TableName          string           `json:"tableName"`
	Timer              int64            `json:"timer"`
	RequestDellay      int64            `json:"requestDellay"`
	Random             bool             `json:"random"`
	WritesNumberToSend int64            `json:"writesNumberToSend"`
	TotalToSend        int64            `json:"totalToSend"`
	StopWhenTableEnds  bool             `json:"stopWhenTableEnds"`
	ConsumerSettings   ConsumerSettings `json:"consumerSettings"`
}

type ConsumerSettings struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// func NewBackgroundWorkerSettings(timer, writesNumberToSend, totalToSend, requestDellay int64, random, stpoWhenTableEnds bool, tableName string) BackgroundWorkerSettings {
// 	return BackgroundWorkerSettings{
// 		tableName:          tableName,
// 		timer:              timer,
// 		requestDellay:      requestDellay,
// 		random:             random,
// 		writesNumberToSend: writesNumberToSend,
// 		totalToSend:        totalToSend,
// 		stpoWhenTableEnds:  stpoWhenTableEnds,
// 	}
// }

// func (b BackgroundWorkerSettings) TableName() string {
// 	return b.tableName
// }

// func (b BackgroundWorkerSettings) Timer() int64 {
// 	return b.timer
// }

// func (b BackgroundWorkerSettings) RequestDellay() int64 {
// 	return b.requestDellay
// }

// func (b BackgroundWorkerSettings) IsRandomSending() bool {
// 	return b.random
// }

// func (b BackgroundWorkerSettings) WritesNumberToSend() int64 {
// 	return b.writesNumberToSend
// }

// func (b BackgroundWorkerSettings) TotalToSend() int64 {
// 	return b.totalToSend
// }

// func (b BackgroundWorkerSettings) ShoudStopWhenTableEnds() bool {
// 	return b.stpoWhenTableEnds
// }
