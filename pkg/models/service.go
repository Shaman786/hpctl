package models

type ServiceResponse struct {
	Services []ServiceItem `json:"services"`
}

type ServiceItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`     // e.g. SSDVIN-4G
	Domain       string `json:"domain"`   // e.g. kvm9298
	Category     string `json:"category"` // e.g. SSD VPS ( IN )
	Price        string `json:"total"`
	BillingCycle string `json:"billingcycle"`
	NextDue      string `json:"next_due"`
	Status       string `json:"status"` // e.g. Cancelled, Active
}

type ServiceDetailResponse struct {
	Service ServiceDetail `json:"service"` // The API usually wraps it in a root key
}

type ServiceDetail struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Domain           string `json:"domain"`
	DedicatedIP      string `json:"dedicatedip"` // specific to single view
	Status           string `json:"status"`
	RegistrationDate string `json:"regdate"`
	NextDueDate      string `json:"nextduedate"`
	BillingCycle     string `json:"billingcycle"`
	Amount           string `json:"recurringamount"`
	PaymentMethod    string `json:"paymentmethodname"`
}
