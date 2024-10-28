package handlers

type Asset struct {
    DealerID    string  `json:"dealerId"`
    MSISDN      string  `json:"msisdn"`
    MPIN        string  `json:"mpin"`
    Balance     float64 `json:"balance"`
    Status      string  `json:"status,omitempty"`
    TransAmount float64 `json:"transAmount,omitempty"`
    TransType   string  `json:"transType,omitempty"`
    Remarks     string  `json:"remarks,omitempty"`
}