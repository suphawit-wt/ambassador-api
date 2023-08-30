package models

type Link struct {
	Id       uint      `json:"id"`
	Code     string    `json:"code"`
	UserId   uint      `json:"user_id"`
	User     User      `json:"user" gorm:"foreignKey:UserId"`
	Products []Product `json:"products" gorm:"many2many:link_products;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Orders   []Order   `json:"orders,omitempty" gorm:"-"`
}

type CreateLinkRequest struct {
	Products []int `json:"products_id"`
}
