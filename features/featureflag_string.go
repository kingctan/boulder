// Code generated by "stringer -type=FeatureFlag"; DO NOT EDIT.

package features

import "strconv"

const _FeatureFlag_name = "unusedAllowAccountDeactivationAllowKeyRolloverResubmitMissingSCTsOnlyUseAIAIssuerURLAllowTLS02ChallengesGenerateOCSPEarlyReusePendingAuthzCountCertificatesExactRandomDirectoryEntryIPv6FirstDirectoryMetaAllowRenewalFirstRLRecheckCAAUDPDNSROCACheckWildcardDomainsEnforceChallengeDisable"

var _FeatureFlag_index = [...]uint16{0, 6, 30, 46, 69, 84, 104, 121, 138, 160, 180, 189, 202, 221, 231, 237, 246, 261, 284}

func (i FeatureFlag) String() string {
	if i < 0 || i >= FeatureFlag(len(_FeatureFlag_index)-1) {
		return "FeatureFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FeatureFlag_name[_FeatureFlag_index[i]:_FeatureFlag_index[i+1]]
}
