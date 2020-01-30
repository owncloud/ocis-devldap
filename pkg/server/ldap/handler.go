package ldap

import (
	"encoding/asn1"
	"encoding/hex"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/butonic/ldapserver/pkg/constants"
	"github.com/butonic/ldapserver/pkg/ldap"
	"github.com/lor00x/goldap/message"
	"github.com/micro/go-micro/util/log"
)

// Handler implements the handlers for LDAP Requests
type Handler struct {
	data *gabs.Container
}

// NotFound returns success for a bind request, unwilling to perform otherwise
func (h *Handler) NotFound(w ldap.ResponseWriter, r *ldap.Request) {
	switch r.ProtocolOpType() {
	case constants.ApplicationBindRequest:
		res := ldap.NewBindResponse(constants.LDAPResultSuccess)
		res.SetDiagnosticMessage("Default binding behavior set to return Success")

		w.Write(res)

	default:
		res := ldap.NewResponse(constants.LDAPResultUnwillingToPerform)
		res.SetDiagnosticMessage("Operation not implemented by server")
		w.Write(res)
	}
}

// Abandon tries to Abandon the given request
func (h *Handler) Abandon(w ldap.ResponseWriter, m *ldap.Request) {
	var req = m.GetAbandonRequest()
	// retrieve the request to abandon, and send a abort signal to it
	if requestToAbandon, ok := m.Conn.GetMessageByID(int(req)); ok {
		requestToAbandon.Abandon()
		log.Infof("Abandon signal sent to request processor [messageID=%d]", int(req))
	}
}

// Bind authenticates using user and password
// the password is taken from a `userpassword` property, which is non standard
func (h *Handler) Bind(w ldap.ResponseWriter, m *ldap.Request) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(constants.LDAPResultSuccess)
	if r.AuthenticationChoice() == "simple" {
		password := h.data.Search(string(r.Name()), "userpassword").Data()
		log.Debugf("User=%s, password=%v", string(r.Name()), password)

		if password == nil {
			log.Debugf("User=%s, has no userpassword", string(r.Name()))
		} else if string(r.AuthenticationSimple()) == password.(string) {
			w.Write(res)
			return
		}
		log.Debugf("Bind failed User=%s, Pass=%#v", string(r.Name()), r.Authentication())
		res.SetResultCode(constants.LDAPResultInvalidCredentials)
		res.SetDiagnosticMessage("invalid credentials")
	} else {
		res.SetResultCode(constants.LDAPResultUnwillingToPerform)
		res.SetDiagnosticMessage("Authentication choice not supported")
	}

	w.Write(res)
}

// Extended handles extended requests
func (h *Handler) Extended(w ldap.ResponseWriter, m *ldap.Request) {
	r := m.GetExtendedRequest()
	log.Debugf("Extended request received, name=%s", r.RequestName())
	log.Debugf("Extended request received, value=%x", r.RequestValue())
	res := ldap.NewExtendedResponse(constants.LDAPResultSuccess)
	w.Write(res)
}

// WhoAmI TODO return the currently logged in user
func (h *Handler) WhoAmI(w ldap.ResponseWriter, m *ldap.Request) {
	res := ldap.NewExtendedResponse(constants.LDAPResultSuccess)
	w.Write(res)
}

type searchControlValue struct {
	Size   int
	Cookie string
}

func addAttributeValue(e *message.SearchResultEntry, attribute message.LDAPString, values []string) {
	//log.Printf("Adding Attribute %s with values %s", attribute, values)
	attributeValues := make([]message.AttributeValue, len(values))
	for i, value := range values {
		if strings.HasPrefix(value, "{hex}") {
			bytes, err := hex.DecodeString(value[5:])
			if err != nil {
				log.Debugf("could not decode hex string %s to bytes", value)
			}
			log.Debugf("Adding Attribute %s with hex value %s", attribute, value)
			attributeValues[i] = message.AttributeValue(bytes)
		} else {
			log.Debugf("Adding Attribute %s with value %s", attribute, value)
			attributeValues[i] = message.AttributeValue(value)
		}
	}
	e.AddAttribute(message.AttributeDescription(attribute), attributeValues...)
}

// Search looks up nodes and attributes in the tree
func (h *Handler) Search(w ldap.ResponseWriter, m *ldap.Request) {
	r := m.GetSearchRequest()
	log.Debugf("Request BaseDn=%s", r.BaseObject())
	log.Debugf("Request Filter=%s", r.Filter())
	log.Debugf("Request FilterString=%s", r.FilterString())
	log.Debugf("Request Attributes=%s", r.Attributes())
	log.Debugf("Request TimeLimit=%d", r.TimeLimit().Int())
	log.Debugf("Request Controls=%+v", m.Controls())

	// Handle Stop Signal (server stop / client disconnected / Abandoned request....)
	select {
	case <-m.Done:
		log.Debugf("Leaving handleSearch...")
		return
	default:
	}

	if m.Controls() != nil {
		for _, control := range *m.Controls() {
			if control.ControlType() == "1.2.840.113556.1.4.319" {
				var controlValue searchControlValue
				/*rest, err := */ asn1.Unmarshal(control.ControlValue().Bytes(), &controlValue)
				log.Debugf("Paged search request %+v", controlValue)
				// TODO implement paged search
			}
		}
	}

	children, _ := h.data.ChildrenMap()
	for key, child := range children {
		if strings.HasSuffix(key, string(r.BaseObject())) {
			log.Debugf("checking node: %v\n", key)
			if matches(child, r.Filter()) {
				log.Debugf("found match %v\n", child)
				e := ldap.NewSearchResultEntry(key)
				for _, ldapAttribute := range r.Attributes() {
					attribute := strings.ToLower(string(ldapAttribute))
					if attribute == "dn" {
						continue
					}
					value := child.Search(attribute)
					log.Debugf("checking attribute: %+v for value: %+v\n", attribute, value)
					if value != nil {

						children, err := value.Children()
						var values []string
						if err != nil {
							values = []string{value.Data().(string)}
						} else {
							values = make([]string, len(children))
							for i, child := range children {
								values[i] = child.Data().(string)
							}
						}
						addAttributeValue(&e, ldapAttribute, values)

					}
				}
				w.Write(e)
			}
		} else {
			log.Debugf("node: %v not in basedn %v\n", key, r.BaseObject())
		}
	}

	res := ldap.NewSearchResultDoneResponse(constants.LDAPResultSuccess)
	w.Write(res)

}
