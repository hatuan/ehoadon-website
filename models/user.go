package models

// User represents the user model
type User struct {
	ID                  *int64     `db:"id" json:",string"`
	Name                string     `db:"name"`
	Password            string     `json:"-" db:"password"`
	Salt                string     `json:"-" db:"salt"`
	Comment             string     `db:"comment"`
	FullName            string     `db:"full_name"`
	PasswordAnswer      string     `json:"-" db:"password_answer"`
	PasswordQuestion    string     `json:"-" db:"password_question"`
	Email               string     `json:"email"`
	CreatedDate         *Timestamp `json:"omitempty" db:"created_date"`
	IsActivated         bool       `json:"is_activated" db:"is_activated"`
	IsLockedOut         bool       `json:"is_locked_out" db:"is_locked_out"`
	LastLockedOutDate   *Timestamp `json:"omitempty" db:"last_locked_out_date"`
	LastLockedOutReason string     `db:"last_locked_out_reason"`
	LastLoginDate       *Timestamp `json:"omitempty" db:"last_login_date"`
	LastLoginIP         string     `db:"last_login_ip"`
	LastModifiedDate    *Timestamp `json:"omitempty" db:"last_modified_date"`
	ClientID            int64      `db:"client_id" json:",string"`
	OrganizationID      int64      `db:"organization_id" json:",string"`
	CultureUIID         string     `db:"culture_ui_id"`
}
