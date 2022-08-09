package shared

// MessageKey type
type MessageKey int

const (

	// MemberRegistration message key
	MemberRegistration MessageKey = iota

	// MemberUpdate message key
	MemberUpdate

	// MemberActivation message key
	MemberActivation

	// MerchantRegistration message key
	MerchantRegistration

	//AddressCreate message key
	AddressCreate

	//AddressModify message key
	AddressModify

	//AddressPrimary message key
	AddressPrimary

	//AddressDelete message key
	AddressDelete

	textMemberRegistration = "member-registration"
)

// String method
func (k MessageKey) String() string {
	switch k {
	case MemberRegistration:
		return textMemberRegistration
	case MemberUpdate:
		return "member-update"
	case MemberActivation:
		return "member-activation"
	case MerchantRegistration:
		return "merchant-registration"
	case AddressCreate:
		return "address-create"
	case AddressModify:
		return "address-modify"
	case AddressPrimary:
		return "address-primary"
	case AddressDelete:
		return "address-delete"
	default:
		return textMemberRegistration
	}
}

//MessageKeyFromString convert string to MessageKey
func MessageKeyFromString(s string) MessageKey {
	switch s {
	case textMemberRegistration, "register":
		return MemberRegistration
	case "member-update", "update":
		return MemberUpdate
	case "member-activation", "activation":
		return MemberActivation
	case "address-create", "create":
		return AddressCreate
	case "address-modify", "modify":
		return AddressModify
	case "address-primary", "primary":
		return AddressPrimary
	case "address-delete", "delete":
		return AddressDelete
	case "merchant-registration", "merchant-register":
		return MerchantRegistration
	default:
		return MemberRegistration
	}
}
