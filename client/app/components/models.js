import {Record} from 'immutable';

import instanceTypeAccessModel from 'instance-types/access-list-model';
import alarmDefinitionModel from 'alarms-definitions/alarm-definition-model';
import alarmModel from 'alarms/alarm-model';
import appInstanceModel from 'app-catalog/app-instance-model';
import appModel from 'app-catalog/app-template-model';
import authModel from 'auth/auth-model';
import bankModel from 'protection/banks/bank-model';
import bankResourceModel from 'protection/resources/bank-resource-model';
import bankTestModel from 'protection/banks/bank-test-model';
import batchModel from 'models/batch-model';
import bondInterfaceModel from 'network/bond/bond-interface-model';
import bondSlaveModel from 'network/bond/bond-slave-model';
import bucketModel from 'object-store/bucket-model';
import cassandraModel from 'db-clusters/cassandra/cassandra-model';
import cassandraNodeModel from 'db-clusters/cassandra/nodes/cassandra-node-model';
import certificateModel from 'certificates/certificate-model';
import clusterSettingsModel from 'settings/models/cluster-settings-model';
import clusterServicesModel from 'settings/models/cluster-services-model';
import clusterStatsModel from 'cluster/cluster-stats-model';
import clusterSummaryModel from 'cluster/cluster-summary-model';
import computeRuleModel from 'compute-rules/compute-rule-model';
import conversionModel from 'tasks/conversion-model';
import dbsInstanceLogsModel from 'rds/logs/dbs-instance-logs-model';
import dbsInstanceLogModel from 'rds/logs/dbs-instance-log-model';
import defaultStoragePoolModel from 'storage/pools/models/default-pool-model';
import diskModel from 'storage/disks/disk-model';
import domainLimitsModel from 'limits/domain-limits-model';
import domainModel from 'accounts/tenants/tenant-model';
import domainLdapModel from 'accounts/tenants/tenant-ldap-model';
import ec2CredentialsModel from 'auth/credentials-model';
import engineModel from 'engine-manager/engines/engine-model';
import engineRevisionModel from 'engine-manager/engine-revisions/engine-revision-model';
import engineVersionModel from 'engine-manager/engine-versions/engine-version-model';
import ethernetInterfaceModel from 'network/ethernet/ethernet-interface-model';
import eventQueryModel from 'events/event-query-model';
import externalCredentialsModel from 'protection/external-credentials/external-credentials-model';
import eventDefinitionModel from 'events/event-definition-model';
import externalPoolModel from 'storage/pools/external-pool-model';
import externalPoolOptionsModel from 'storage/pools/external-pool-options-model';
import floatingIpModel from 'floating-ip/floating-ip-model';
import imageModel from 'storage/images/models/image-model';
import instanceTypeModel from 'instance-types/instance-type-model';
import ipv4Model from 'network/ip-interfaces/ipv4-interface-model';
import k8sClustersModel from 'kubernetes/clusters/clusters-model';
import k8sNodesModel from 'kubernetes/nodes/nodes-model';
import k8sPodsModel from 'kubernetes/pods/pods-model';
import lbaasModel from 'lbaas/lbaas-model';
import lbaasListenerModel from 'lbaas/lbaas-listener-model';
import lbaasStatusModel from 'lbaas/lbaas-init-status-model';
import lbaasTargetsModel from 'lbaas/target-groups/single-target-group/lbaas-target-model';
import lbaasTargetGroupsModel from 'lbaas/target-groups/target-group-model';
import logicalNetworkModel from 'network/logical-networks/logical-network-model';
import logicalSubnetModel from 'network/subnets/subnet-model';
import maestroModel from 'maestro/maestro-model';
import mapreduceModel from 'analytics/mapreduce/mapreduce-model';
import mapreduceInstanceModel from 'analytics/mapreduce/instances/mapreduce-instance-model';
import marketPlaceAppModel from 'app-catalog-s3/app-template-model';
import metricModel from 'metrics/metric-model';
import metricRangeModel from 'metrics/metric-range-model';
import metricTopModel from 'metrics/metric-top-model';
import networkModel from 'network/network-model';
import networkAllocationModel from 'network/network-allocation-model';
import networkRouteModel from 'network/network-routes/network-route-model';
import networkEc2ExternalModel from 'network/network-ec2-external-model';
import networkEc2ProjectModel from 'network/network-ec2-project-model';
import neutronQuotasModel from 'quotas/quota-neutron-model';
import nodeHardwareModel from 'node/node-hardware-model';
import nodeModel from 'node/node-model';
import objectModel from 'object-store/object-model';
import objectAclModel from 'object-store/object-acl-model';
import poolModel from 'storage/pools/models/pool-model';
import portModel from 'ports/models/port-model';
import projectModel from 'accounts/projects/project-model';
import protectionDefinitionModel from 'protection/definitions/protection-definition-model';
import protectionPlanModel from 'protection/plans/protection-plan-model';
import protectionTaskModel from 'protection/tasks/protection-task-model';
import protectionTriggerModel from 'protection/triggers/protection-trigger-model';
import providerModel from 'providers/provider-model';
import proxySettingsModel from 'settings/http-proxy/http-proxy-model';
import rdsInstancesModel from 'rds/instances/rds-instances-model';
import rdsParameterGroupModel from 'rds/parameter-group/rds-parameter-groups-model';
import rdsParametersModel from 'rds/parameters/rds-parameters-model';
import rdsSnapshotsModel from 'rds/snapshots/rds-snapshot-model';
import roleAssignmentsModel from 'accounts/users/user-role-assignments-model';
import routerModel from 'routers/router-model';
import securityGroupModel from 'network/security-groups/security-group-model';
import securityGroupRuleModel from 'network/security-groups/security-group-rules-model';
import serviceModel from 'services/service-model';
import serviceListModel from 'services/service-list-model';
import serviceCpuMetricTopModel from 'services/services-cpu-metric-top-model';
import serviceMemoryMetricTopModel from 'services/services-memory-metric-top-model';
import slaProfileModel from 'vms/sla-profile-model';
import snapshotModel from 'storage/snapshots/snapshot-model';
import storageQuotasModel from 'quotas/quota-storage-model';
import tagModel from 'tags/tag-model';
import targetCheckpointsModel from 'protection/checkpoints/target-checkpoint-model';
import targetModel from 'protection/targets/target-model';
import targetResourceModel from 'protection/resources/target-resource-model';
import trafficInterfaceModel from 'network/traffic/traffic-interface-model';
import trafficTypeModel from 'network/traffic/traffic-type-model';
import uiImagesModel from 'settings/models/ui-images-model';
import userModel from 'accounts/users/user-model';
import userProjectModel from 'accounts/users/user-projects-model';
import virtualIpModel from 'network/vips/vip-model';
import virtualSubnetModel from 'network/subnets/virtual-subnet-model';
import vlanModel from 'network/nodes-networks/nodes-networks-model';
import vlanInterfaceModel from 'network/vlans/vlan-interface-model';
import vmModel from 'vms/vm-model';
import vmNovaModel from 'vms/models/vm-nova-model';
import vncModel from 'vms/models/vnc-model';
import vnGroupModel from 'network/vn-groups/vn-group-model';
import vnTypeModel from 'network/vn-types/vn-type-model';
import volumeModel from 'storage/volume/volume-model';
import nfsFileSystemModel from 'storage/nfs/filesystems-model';
import nfsMountTargetModel from 'storage/nfs/mount-targets-model';
import nfsEngineModel from 'storage/nfs/engine-model';
import objectStoreModel from 'object-store/object-store-model';
import upgradesModel from 'upgrades/upgrades-model';
import upgradeTasksModel from 'upgrades/upgrade-tasks-model';
import upgradeGroupTasksModel from 'upgrades/upgrade-group-tasks-model';
import upgradesReleaseNotesModel from 'upgrades/upgrades-release-notes-model';
import userSettingsModel from 'settings/models/user-settings-model';
import vmwareAdapterModel from 'providers/adapters/vmware-adapter-model';

const models = Record({
  instanceTypeAccess:      instanceTypeAccessModel,
  alarm:                   alarmModel,
  alarmDefinition:         alarmDefinitionModel,
  app:                     appModel,
  appInstance:             appInstanceModel,
  auth:                    authModel,
  bank:                    bankModel,
  bankResource:            bankResourceModel,
  bankTest:                bankTestModel,
  batch:                   batchModel,
  bondInterface:           bondInterfaceModel,
  bondSlave:               bondSlaveModel,
  bucket:                  bucketModel,
  cassandra:               cassandraModel,
  cassandraNode:           cassandraNodeModel,
  certificate:             certificateModel,
  clusterSettings:         clusterSettingsModel,
  clusterServices:         clusterServicesModel,
  clusterStats:            clusterStatsModel,
  clusterSummary:          clusterSummaryModel,
  computeRule:             computeRuleModel,
  conversion:              conversionModel,
  dbsInstanceLogs:         dbsInstanceLogsModel,
  dbsInstanceLog:          dbsInstanceLogModel,
  defaultStoragePool:      defaultStoragePoolModel,
  disk:                    diskModel,
  domain:                  domainModel,
  domainLdap:              domainLdapModel,
  domainLimits:            domainLimitsModel,
  ec2Credentials:          ec2CredentialsModel,
  engine:                  engineModel,
  engineRevision:          engineRevisionModel,
  engineVersion:           engineVersionModel,
  ethInterface:            ethernetInterfaceModel,
  eventsDefinition:        eventDefinitionModel,
  eventQuery:              eventQueryModel,
  externalCredentials:     externalCredentialsModel,
  externalPool:            externalPoolModel,
  externalPoolOptions:     externalPoolOptionsModel,
  floatingIp:              floatingIpModel,
  image:                   imageModel,
  instanceType:            instanceTypeModel,
  ipv4:                    ipv4Model,
  k8sCluster:              k8sClustersModel,
  k8sNode:                 k8sNodesModel,
  k8sPod:                  k8sPodsModel,
  lbaas:                   lbaasModel,
  lbaasListener:           lbaasListenerModel,
  lbaasStatusModel:        lbaasStatusModel,
  lbaasTarget:             lbaasTargetsModel,
  lbaasTargetGroup:        lbaasTargetGroupsModel,
  logicalNetwork:          logicalNetworkModel,
  logicalSubnet:           logicalSubnetModel,
  mapreduce:               mapreduceModel,
  mapreduceInstance:       mapreduceInstanceModel,
  maestro:                 maestroModel,
  marketplaceApp:          marketPlaceAppModel,
  metric:                  metricModel,
  metricRange:             metricRangeModel,
  metricTop:               metricTopModel,
  network:                 networkModel,
  networkAllocation:       networkAllocationModel,
  networkRoute:            networkRouteModel,
  networkEc2External:      networkEc2ExternalModel,
  networkEc2Project:       networkEc2ProjectModel,
  neutronQuotas:           neutronQuotasModel,
  node:                    nodeModel,
  nodeHardware:            nodeHardwareModel,
  object:                  objectModel,
  objectAcl:               objectAclModel,
  objectStore:             objectStoreModel,
  pool:                    poolModel,
  port:                    portModel,
  project:                 projectModel,
  protectionDefinition:    protectionDefinitionModel,
  protectionPlan:          protectionPlanModel,
  protectionTask:          protectionTaskModel,
  protectionTrigger:       protectionTriggerModel,
  proxySettings:           proxySettingsModel,
  provider:                providerModel,
  rdsInstances:            rdsInstancesModel,
  rdsParameterGroup:       rdsParameterGroupModel,
  rdsParameter:            rdsParametersModel,
  rdsSnapshots:            rdsSnapshotsModel,
  nfsMountTarget:          nfsMountTargetModel,
  nfsFileSystem:           nfsFileSystemModel,
  nfsEngine:               nfsEngineModel,
  roleAssignment:          roleAssignmentsModel,
  router:                  routerModel,
  securityGroup:           securityGroupModel,
  securityGroupRule:       securityGroupRuleModel,
  service:                 serviceModel,
  serviceList:             serviceListModel,
  serviceCpuTop:           serviceCpuMetricTopModel,
  serviceMemoryTop:        serviceMemoryMetricTopModel,
  slaProfile:              slaProfileModel,
  snapshot:                snapshotModel,
  storageQuotas:           storageQuotasModel,
  tag:                     tagModel,
  target:                  targetModel,
  targetCheckpoint:        targetCheckpointsModel,
  targetResource:          targetResourceModel,
  trafficInterface:        trafficInterfaceModel,
  trafficType:             trafficTypeModel,
  uiImages:                uiImagesModel,
  user:                    userModel,
  userProject:             userProjectModel,
  virtualIp:               virtualIpModel,
  virtualSubnet:           virtualSubnetModel,
  vlanInterface:           vlanInterfaceModel,
  vlan:                    vlanModel,
  vm:                      vmModel,
  vmNova:                  vmNovaModel,
  vmwareAdapter:           vmwareAdapterModel,
  vnc:                     vncModel,
  vnGroup:                 vnGroupModel,
  vnType:                  vnTypeModel,
  volume:                  volumeModel,
  upgrade:                 upgradesModel,
  upgradeTask:             upgradeTasksModel,
  upgradeGroupTask:        upgradeGroupTasksModel,
  upgradeReleaseNotes:     upgradesReleaseNotesModel,
  userSettings:            userSettingsModel,
});

export default new models();
