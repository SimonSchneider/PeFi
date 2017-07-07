package model

type (
	Loan struct {
		Transaction
		PaybackIds []int64 `json:"payback_ids"`
	}
)

//func (l *Loan) Table() (s []string) {
//s = l.Transaction.Table()
//paybackIds := []string{}
//for _, id := range l.PaybackIds {
//paybackIds = append(paybackIds, strconv.Itoa(int(id)))
//}
//s = append(s, strings.Join(paybackIds, ","))
//return s
//}
