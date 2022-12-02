package anyframer

// Column ...
type Column struct {
	Selector   string       `json:"selector,omitempty"`
	Alias      string       `json:"alias,omitempty"`
	Format     ColumnFormat `json:"format,omitempty"`
	TimeFormat string       `json:"timeFormat,omitempty"`
}

// ColumnFormat ...
type ColumnFormat string

const (
	// ColumnFormatString ...
	ColumnFormatString ColumnFormat = "string"
	// ColumnFormatNumber ...
	ColumnFormatNumber ColumnFormat = "number"
	// ColumnFormatBoolean ...
	ColumnFormatBoolean ColumnFormat = "boolean"
	// ColumnFormatTimeStamp ...
	ColumnFormatTimeStamp ColumnFormat = "timestamp"
	// ColumnFormatUnixMsecTimeStamp ...
	ColumnFormatUnixMsecTimeStamp ColumnFormat = "timestamp_epoch"
	// ColumnFormatUnixSecTimeStamp ...
	ColumnFormatUnixSecTimeStamp ColumnFormat = "timestamp_epoch_s"
)
