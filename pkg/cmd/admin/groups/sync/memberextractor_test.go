package syncgroups

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/ldap.v2"

	"github.com/openshift/origin/pkg/cmd/admin/groups/sync/interfaces"
)

func TestExtractMembers(t *testing.T) {
	var testCases = []struct {
		name            string
		baseExtractor   interfaces.LDAPMemberExtractor
		nameMappper     interfaces.LDAPUserNameMapper
		ldapGroupUID    string
		blacklist       []string
		expectedMembers []*ldap.Entry
	}{
		{
			name:            "blacklist is subset of whitelist",
			baseExtractor:   &parrotLDAPMemberExtractor{members: []*ldap.Entry{newTestUser("bob"), newTestUser("alice"), newTestUser("george")}},
			nameMappper:     NewUserNameMapper([]string{"cn"}),
			ldapGroupUID:    "some-group",
			blacklist:       []string{"bob"},
			expectedMembers: []*ldap.Entry{newTestUser("alice"), newTestUser("george")},
		},
		{
			name:            "blacklist is a superset of whitelist",
			baseExtractor:   &parrotLDAPMemberExtractor{members: []*ldap.Entry{newTestUser("bob"), newTestUser("alice"), newTestUser("george")}},
			nameMappper:     NewUserNameMapper([]string{"cn"}),
			ldapGroupUID:    "some-group",
			blacklist:       []string{"bob", "alice", "george", "someone else"},
			expectedMembers: []*ldap.Entry{},
		},
		{
			name:            "blacklist is empty", // should not happen if the create flow is working correctly
			baseExtractor:   &parrotLDAPMemberExtractor{members: []*ldap.Entry{newTestUser("bob"), newTestUser("alice"), newTestUser("george")}},
			nameMappper:     NewUserNameMapper([]string{"cn"}),
			ldapGroupUID:    "some-group",
			blacklist:       []string{},
			expectedMembers: []*ldap.Entry{newTestUser("bob"), newTestUser("alice"), newTestUser("george")},
		},
	}

	for _, testCase := range testCases {
		memberExtractor := NewBlacklistLDAPMemberExtractor(testCase.blacklist, testCase.baseExtractor, testCase.nameMappper)

		members, err := memberExtractor.ExtractMembers(testCase.ldapGroupUID)
		if err != nil {
			t.Errorf("%s: unexpected error extracting members: %v", err)
			continue
		}

		if actual, expected := members, testCase.expectedMembers; !reflect.DeepEqual(actual, expected) {
			t.Errorf("%s: member extractor did not generate correct list of names:\n\texpected:\n\t%v\n\tgot:\n\t%v", testCase.name, expected, actual)
		}
	}
}

// parrotLDAPMemberExtractor parrots the list of members given to it for use in testing
type parrotLDAPMemberExtractor struct {
	members []*ldap.Entry
}

func (e *parrotLDAPMemberExtractor) ExtractMembers(ldapGroupUID string) ([]*ldap.Entry, error) {
	return e.members, nil
}

// newTestUser returns a new LDAP entry with the CN
func newTestUser(CN string) *ldap.Entry {
	return ldap.NewEntry(fmt.Sprintf("cn=%s,ou=users,dc=example,dc=com", CN), map[string][]string{"cn": {CN}})
}
