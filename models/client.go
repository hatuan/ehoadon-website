package models

import (
	"database/sql"
	"erpvietnam/ehoadon-website/crypto"
	"erpvietnam/ehoadon-website/log"
	"erpvietnam/ehoadon-website/settings"
	"erpvietnam/ehoadon-website/utils"
	"fmt"
	"html/template"
	"net"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/ssh"

	"errors"

	_ "github.com/lib/pq"
)

// Client represents the client model
type Client struct {
	ClientID                    *int64          `db:"id" json:",string"`
	Code                        string          `db:"code"`
	Description                 string          `db:"description"`
	IsActivated                 bool            `db:"is_activated"`
	ActivatedCode               string          `db:"activated_code"`
	CultureID                   string          `db:"culture_id"`
	AmountDecimalPlaces         int16           `db:"amount_decimal_places"`
	AmountRoundingPrecision     decimal.Decimal `db:"amount_rounding_precision" json:",string"`
	UnitAmountDecimalPlaces     int16           `db:"unit_amount_decimal_places"`
	UnitAmountRoundingPrecision decimal.Decimal `db:"unit_amount_rounding_precision" json:",string"`
	CurrencyLCYId               int64           `db:"currency_lcy_id" json:",string"`
	VatNumber                   string          `db:"vat_number"`
	GroupUnitCode               string          `db:"group_unit_code"`
	VatMethodCode               string          `db:"vat_method_code"`
	ProvinceCode                string          `db:"province_code"`
	DistrictsCode               string          `db:"districts_code"`
	Address                     string          `db:"address"`
	AddressTransition           string          `db:"address_transition"`
	Telephone                   string          `db:"telephone"`
	Email                       string          `db:"email"`
	Fax                         string          `db:"fax"`
	Website                     string          `db:"website"`
	RepresentativeName          string          `db:"representative_name"`
	RepresentativePosition      string          `db:"representative_position"`
	ContactName                 string          `db:"contact_name"`
	Mobile                      string          `db:"mobile"`
	BankAccount                 string          `db:"bank_account"`
	BankName                    string          `db:"bank_name"`
	TaxAuthoritiesID            *int64          `db:"tax_authorities_id"`
	Version                     int16           `db:"version"`
	RecCreatedByID              int64           `db:"rec_created_by" json:",string"`
	RecCreated                  *Timestamp      `db:"rec_created_at"`
	RecModifiedByID             int64           `db:"rec_modified_by" json:",string"`
	RecModified                 *Timestamp      `db:"rec_modified_at"`
}

// ErrOrganizationsIsEmpty is thrown when do not found any Organization.
var ErrOrganizationsIsEmpty = errors.New("Organizations is empty")

// ErrClientActiveCodeExpired is thrown when Client with active code expired.
var ErrClientActiveCodeExpired = errors.New("Client Active Code Expired")

// ErrClientActiveCodeNotFound is thrown when do not found any Client with active code.
var ErrClientActiveCodeNotFound = errors.New("Client Active Code not found")

// ErrClientCreateNewCompanyFail is thrown when create new company fail
var ErrClientCreateNewCompanyFail = errors.New("Client create new company fail")

// ErrClientNotFound is thrown when do not found any Client.
var ErrClientNotFound = errors.New("Client not found")

// ErrClientValidate indicates there was validate error
var ErrClientValidate = errors.New("Client has validate error")

// ErrClientCodeNotSpecified indicates there was no code given by the user
var ErrClientCodeNotSpecified = errors.New("Client's code not specified")

// ErrClientCodeDuplicate indicates there was duplicate of code given by the user
var ErrClientCodeDuplicate = errors.New("Client's code is duplicate")

// ErrClientDescriptionNotSpecified indicates there was no name given by the user
var ErrClientDescriptionNotSpecified = errors.New("Client's Description not specified")

// ErrClientAddressNotSpecified indicates there was no name given by the user
var ErrClientAddressNotSpecified = errors.New("Client's Address not specified")

// ErrClientFatal indicates there was fatal error
var ErrClientFatal = errors.New("Client has fatal error")

// Validate checks to make sure there are no invalid fields in a submitted
func (c *Client) Validate() map[string]InterfaceArray {
	validationErrors := make(map[string]InterfaceArray)

	if c.Description == "" {
		validationErrors["Description"] = append(validationErrors["Description"], ErrClientDescriptionNotSpecified.Error())
	}

	if c.Address == "" {
		validationErrors["Address"] = append(validationErrors["Address"], ErrClientAddressNotSpecified.Error())
	}

	if c.Code == "" {
		validationErrors["Code"] = append(validationErrors["Code"], ErrClientCodeNotSpecified.Error())
	}

	if c.Code != "" {
		var otherID string
		ID := int64(0)
		if c.ClientID != nil {
			ID = *c.ClientID
		}
		err := DB.Get(&otherID, "SELECT id FROM client WHERE code = $1 AND id != $2", c.Code, ID)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			validationErrors["Fatal"] = append(validationErrors["Fatal"], ErrClientFatal.Error())
		}
		if otherID != "" && err != sql.ErrNoRows {
			validationErrors["Code"] = append(validationErrors["Code"], ErrClientCodeDuplicate.Error())
		}
	}

	return validationErrors
}

func (c *Client) Get(id int64) error {

	err := DB.QueryRowx("SELECT client.* "+
		" FROM client "+
		" WHERE client.id=$1 ", id).StructScan(c)
	if err == sql.ErrNoRows {
		return ErrClientNotFound
	} else if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *Client) GetByCode(code string) error {

	code = strings.ToUpper(code)

	err := DB.QueryRowx("SELECT client.* "+
		" FROM client "+
		" WHERE client.code=$1 ", code).StructScan(c)
	if err == sql.ErrNoRows {
		return ErrClientNotFound
	} else if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *Client) GetByActiveCode(activeCode string) error {

	err := DB.QueryRowx("SELECT client.* "+
		" FROM client "+
		" WHERE client.activated_code=$1 ", activeCode).StructScan(c)
	if err == sql.ErrNoRows {
		return ErrClientActiveCodeNotFound
	} else if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *Client) Active(activeCode string) TransactionalInformation {

	err := DB.QueryRowx("SELECT client.* "+
		" FROM client "+
		" WHERE is_activated = false AND client.activated_code=$1", activeCode).StructScan(c)

	if err == sql.ErrNoRows {
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientActiveCodeNotFound.Error()}, ReturnError: []error{ErrClientActiveCodeNotFound}}
	} else if err != nil {
		log.Error(err)
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{err.Error()}, ReturnError: []error{err}}
	}

	//check active code expired
	//durationExpired := time.Duration(24) * time.Hour
	//if c.RecCreated.Add(durationExpired).Before(time.Now()) {
	//	return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientActiveCodeExpired.Error()}, ReturnError: []error{ErrClientActiveCodeExpired}}
	//}

	objectName := fmt.Sprintf("%s", strings.ToLower(template.HTMLEscapeString(c.Code)))

	success := make(chan bool)
	go createNewCompany(objectName, success)

	result := <-success

	if !result {
		log.Error(ErrClientCreateNewCompanyFail)
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientCreateNewCompanyFail.Error()}, ReturnError: []error{ErrClientCreateNewCompanyFail}}
	}
	return TransactionalInformation{ReturnStatus: true, ReturnMessage: []string{"Updated/Created successfully"}}
}

func (c *Client) Update() TransactionalInformation {
	if validateErrs := c.Validate(); len(validateErrs) != 0 {
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientValidate.Error()}, ReturnError: []error{ErrClientValidate}, ValidationErrors: validateErrs}
	}

	stmt, _ := DB.PrepareNamed("INSERT INTO client as client(id, " +
		" code, " +
		" description, " +
		" is_activated, " +
		" activated_code, " +
		" culture_id, " +
		" amount_decimal_places, " +
		" amount_rounding_precision, " +
		" unit_amount_decimal_places, " +
		" unit_amount_rounding_precision, " +
		" currency_lcy_id, " +
		" vat_number, " +
		" group_unit_code, " +
		" vat_method_code, " +
		" province_code, " +
		" districts_code, " +
		" address, " +
		" address_transition, " +
		" telephone, " +
		" email, " +
		" fax, " +
		" website, " +
		" representative_name, " +
		" representative_position, " +
		" contact_name, " +
		" mobile, " +
		" bank_account, " +
		" bank_name, " +
		" tax_authorities_id, " +
		" version, " +
		" rec_created_by, " +
		" rec_modified_by, " +
		" rec_created_at, " +
		" rec_modified_at) " +
		" VALUES ( COALESCE(:id, id_generator()), " +
		" :code, " +
		" :description, " +
		" :is_activated, " +
		" :activated_code, " +
		" :culture_id, " +
		" :amount_decimal_places, " +
		" :amount_rounding_precision, " +
		" :unit_amount_decimal_places, " +
		" :unit_amount_rounding_precision, " +
		" :currency_lcy_id, " +
		" :vat_number, " +
		" :group_unit_code, " +
		" :vat_method_code, " +
		" :province_code, " +
		" :districts_code, " +
		" :address, " +
		" :address_transition, " +
		" :telephone, " +
		" :email, " +
		" :fax, " +
		" :website, " +
		" :representative_name, " +
		" :representative_position, " +
		" :contact_name, " +
		" :mobile, " +
		" :bank_account, " +
		" :bank_name, " +
		" :tax_authorities_id, " +
		" :version, " +
		" :rec_created_by, " +
		" :rec_modified_by, " +
		" :rec_created_at, " +
		" :rec_modified_at) " +
		" ON CONFLICT ON CONSTRAINT pk_client DO UPDATE SET " +
		" code								=	EXCLUDED.code, " +
		" description						=	EXCLUDED.description, " +
		" is_activated						=	EXCLUDED.is_activated, " +
		" activated_code					=	EXCLUDED.activated_code, " +
		" culture_id						=	EXCLUDED.culture_id, " +
		" amount_decimal_places				=	EXCLUDED.amount_decimal_places, " +
		" amount_rounding_precision			=	EXCLUDED.amount_rounding_precision, " +
		" unit_amount_decimal_places		=	EXCLUDED.unit_amount_decimal_places, " +
		" unit_amount_rounding_precision	=	EXCLUDED.unit_amount_rounding_precision, " +
		" currency_lcy_id					=	EXCLUDED.currency_lcy_id, " +
		" vat_number						=	EXCLUDED.vat_number, " +
		" group_unit_code					=	EXCLUDED.group_unit_code, " +
		" vat_method_code					=	EXCLUDED.vat_method_code, " +
		" province_code						=	EXCLUDED.province_code, " +
		" districts_code					=	EXCLUDED.districts_code, " +
		" address							=	EXCLUDED.address, " +
		" address_transition				=	EXCLUDED.address_transition, " +
		" telephone							=	EXCLUDED.telephone, " +
		" email								=	EXCLUDED.email, " +
		" fax								=	EXCLUDED.fax, " +
		" website							=	EXCLUDED.website, " +
		" representative_name				=	EXCLUDED.representative_name, " +
		" representative_position			=	EXCLUDED.representative_position, " +
		" contact_name						=	EXCLUDED.contact_name, " +
		" mobile							=	EXCLUDED.mobile, " +
		" bank_account						=	EXCLUDED.bank_account, " +
		" bank_name							=	EXCLUDED.bank_name, " +
		" tax_authorities_id				=	EXCLUDED.tax_authorities_id, " +
		" version							=	EXCLUDED.version + 1, " +
		" rec_created_by					=	EXCLUDED.rec_created_by, " +
		" rec_modified_by					=	EXCLUDED.rec_modified_by, " +
		" rec_created_at					=	EXCLUDED.rec_created_at, " +
		" rec_modified_at					=	EXCLUDED.rec_modified_at " +
		" WHERE client.version = :version " +
		" RETURNING id")

	var id int64
	err := stmt.Get(&id, c)

	if err != nil && err == sql.ErrNoRows {
		log.Error(err)
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientNotFound.Error()}}
	} else if err != nil {
		log.Error(err)
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{err.Error()}}
	}

	c.ClientID = &id
	err = c.Get(*c.ClientID)
	if err == sql.ErrNoRows {
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{ErrClientNotFound.Error()}}
	} else if err != nil {
		return TransactionalInformation{ReturnStatus: false, ReturnMessage: []string{err.Error()}}
	}
	return TransactionalInformation{ReturnStatus: true, ReturnMessage: []string{"Updated/Created successfully"}}
}

func createNewCompany(name string, success chan bool) {
	sshConfig := &ssh.ClientConfig{
		User: settings.Settings.SSHUser,
		Auth: []ssh.AuthMethod{
			utils.PublicKeyFile(settings.Settings.SSHPrivateKeyPath),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client := &utils.SSHClient{
		Config: sshConfig,
		Host:   settings.Settings.SSHHost,
	}

	cmd := &utils.SSHCommand{
		Path: fmt.Sprintf("NEW_COMPANY_NAME=%s /ehoadon/create_new_company.sh", name),
		Env:  []string{},
		//Stdin:  os.Stdin,
		//Stdout: os.Stdout,
		//Stderr: os.Stderr,
	}

	if err := client.RunCommand(cmd); err != nil {
		log.Error("command run error: ", err)
		success <- false
	}

	success <- true
}

type InitDB struct {
	UserProfile User
	Client      Client
}

func (c *Client) GetInitDB() InitDB {
	salt := utils.RandStringBytes(5)

	initDB := InitDB{
		UserProfile: User{
			Password: crypto.HashPassword("123456", salt),
			Salt:     salt,
			Email:    (*c).Email,
		},
		Client: *c,
	}
	return initDB
}
