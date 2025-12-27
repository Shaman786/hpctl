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
