package syncgroups

import (
	"gopkg.in/ldap.v2"

	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/openshift/origin/pkg/cmd/admin/groups/sync/interfaces"
)

// NewBlacklistLDAPMemberExtractor layers a LDAP user UID blacklist on top of an existing member extractor
func NewBlacklistLDAPMemberExtractor(blacklist []string, baseExtractor interfaces.LDAPMemberExtractor, nameMapper interfaces.LDAPUserNameMapper) interfaces.LDAPMemberExtractor {
	return &blacklistLDAPMemberExtractor{
		blacklist:     sets.NewString(blacklist...),
		baseExtractor: baseExtractor,
		nameMapper:    nameMapper,
	}
}

type blacklistLDAPMemberExtractor struct {
	blacklist     sets.String
	baseExtractor interfaces.LDAPMemberExtractor
	nameMapper    interfaces.LDAPUserNameMapper
}

// ExtractMembers returns the LDAP users whose LDAP user UIDs are not in the user blacklist
func (e *blacklistLDAPMemberExtractor) ExtractMembers(ldapGroupUID string) ([]*ldap.Entry, error) {
	allMembers, err := e.baseExtractor.ExtractMembers(ldapGroupUID)
	if err != nil {
		return nil, err
	}

	wantedMembers := []*ldap.Entry{}
	for _, member := range allMembers {
		ldapUserUID, err := e.nameMapper.UserNameFor(member)
		if err != nil {
			return nil, err
		}

		if !e.blacklist.Has(ldapUserUID) {
			wantedMembers = append(wantedMembers, member)
		}
	}

	return wantedMembers, nil
}
