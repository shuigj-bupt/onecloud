package aws

import (
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/util/secrules"
	"github.com/aws/aws-sdk-go/service/ec2"
	"yunion.io/x/log"
)

type Tags struct {
	Tag []Tag
}

type Tag struct {
	TagKey   string
	TagValue string
}

type SSecurityGroup struct {
	vpc               *SVpc

	RegionId          string
	VpcId             string
	SecurityGroupId   string
	Description       string
	SecurityGroupName string
	Permissions       []secrules.SecurityRule
	Tags              Tags

	// CreationTime      time.Time
	// InnerAccessPolicy string
}

func (self *SSecurityGroup) GetId() string {
	return self.SecurityGroupId
}

func (self *SSecurityGroup) GetName() string {
	if len(self.SecurityGroupName) > 0 {
		return self.SecurityGroupName
	}
	return self.SecurityGroupId
}

func (self *SSecurityGroup) GetGlobalId() string {
	return self.SecurityGroupId
}

func (self *SSecurityGroup) GetStatus() string {
	return ""
}

func (self *SSecurityGroup) Refresh() error {
	if new, err := self.vpc.region.GetSecurityGroupDetails(self.SecurityGroupId); err != nil {
		return err
	} else {
		return jsonutils.Update(self, new)
	}
}

func (self *SSecurityGroup) IsEmulated() bool {
	return false
}

func (self *SSecurityGroup) GetMetadata() *jsonutils.JSONDict {
	if len(self.Tags.Tag) == 0 {
		return nil
	}
	data := jsonutils.NewDict()
	for _, value := range self.Tags.Tag {
		data.Add(jsonutils.NewString(value.TagValue), value.TagKey)
	}
	return data
}

func (self *SSecurityGroup) GetDescription() string {
	return self.Description
}

func (self *SSecurityGroup) GetRules() ([]secrules.SecurityRule, error) {
	rules := make([]secrules.SecurityRule, 0)
	if secgrp, err := self.vpc.region.GetSecurityGroupDetails(self.SecurityGroupId); err != nil {
		return rules, err
	} else {
		rules = secgrp.Permissions
	}

	return rules, nil
}

func (self *SRegion) addSecurityGroupRules(secGrpId string, rule *secrules.SecurityRule) error {
	// todo: add sercurity rules
	return nil
}

func (self *SRegion) addSecurityGroupRule(secGrpId string, rule *secrules.SecurityRule) error {
	// todo: add sercurity rules
	return nil
}

func (self *SRegion) createSecurityGroup(vpcId string, name string, desc string) (string, error) {
	params := &ec2.CreateSecurityGroupInput{}
	params.SetVpcId(vpcId)
	params.SetDescription(desc)
	params.SetGroupName(name)

	group, err := self.ec2Client.CreateSecurityGroup(params)
	if err != nil {
		return "", err
	}

	return *group.GroupId, nil
}

func (self *SRegion) createDefaultSecurityGroup(vpcId string) (string, error) {
	secId, err := self.createSecurityGroup(vpcId, "vpc default", "vpc default group")
	if err != nil {
		return "", err
	}
	// todo : add sercurity rules
	return secId, nil
}

func (self *SRegion) GetSecurityGroupDetails(secGroupId string) (*SSecurityGroup, error) {
	return nil, nil
}

func (self *SRegion) getSecurityGroupByTag(vpcId, secgroupId string) (*SSecurityGroup, error) {
	return nil, nil
}

func (self *SRegion) addTagToSecurityGroup(secgroupId, key, value string, index int) error {
	return nil
}

func (self *SRegion) modifySecurityGroup(secGrpId string, name string, desc string) error {
	return nil
}

func (self *SRegion) syncSecgroupRules(secgroupId string, rules []secrules.SecurityRule) error {
	return nil
}

func (self *SRegion) getSecRules(ingress []*ec2.IpPermission, egress []*ec2.IpPermission) ([]secrules.SecurityRule) {
	rules := []secrules.SecurityRule{}
	for _, p := range ingress {
		ret, err := AwsIpPermissionToYunion(secrules.SecurityRuleIngress, p)
		if err != nil {
			log.Debugf(err.Error())
		}

		for _, rule := range ret {
			rules = append(rules, rule)
		}
	}

	for _, p := range egress {
		ret, err := AwsIpPermissionToYunion(secrules.SecurityRuleEgress, p)
		if err != nil {
			log.Debugf(err.Error())
		}

		for _, rule := range ret {
			rules = append(rules, rule)
		}
	}

	return rules
}

func (self *SRegion) GetSecurityGroups(vpcId string, offset int, limit int) ([]SSecurityGroup, int, error) {
	params := &ec2.DescribeSecurityGroupsInput{}
	filters := make([]*ec2.Filter, 0)
	if len(vpcId) > 0 {
		filters = AppendSingleValueFilter(filters, "vpc-id", vpcId)
	}

	if len(filters) > 0 {
		params.SetFilters(filters)
	}

	ret, err := self.ec2Client.DescribeSecurityGroups(params)
	if err != nil {
		return nil, 0 , err
	}

	securityGroups := []SSecurityGroup{}
	for _,item := range ret.SecurityGroups {
		vpc, err := self.getVpc(*item.VpcId)
		if err != nil {
			log.Errorf("vpc %s not found", *item.VpcId)
			continue
		}

		permissions := self.getSecRules(item.IpPermissions, item.IpPermissionsEgress)

		group := SSecurityGroup{
			vpc:               vpc,
			Description:       *item.Description,
			SecurityGroupId:   *item.GroupId,
			SecurityGroupName: *item.GroupName,
			VpcId:             *item.VpcId,
			Permissions:       permissions,
			RegionId:          self.RegionId,
			// Tags:              *item.Tags,
		}

		securityGroups = append(securityGroups, group)
	}

	return securityGroups, len(securityGroups), nil
}