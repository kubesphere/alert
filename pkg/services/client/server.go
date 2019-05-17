package client

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

func WebService() *restful.WebService {
	restful.RegisterEntityAccessor(constants.MIME_MERGEPATCH, restful.NewEntityAccessorJSON(restful.MIME_JSON))

	ws := new(restful.WebService)
	ws.Path("/api/v1").Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).Produces(restful.MIME_JSON)

	tags := []string{"ResourceType apis"}

	ws.Route(ws.GET("/resource_type").To(DescribeResourceTypes).
		Doc("Describe Resource Types").
		Param(ws.QueryParameter("rs_type_ids", "resource type id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_type_names", "resource type name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"Metric apis"}

	ws.Route(ws.GET("/metric").To(DescribeMetrics).
		Doc("Describe Metrics").
		Param(ws.QueryParameter("metric_ids", "metric id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("metric_names", "metric name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("status", "metric status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_type_ids", "metric resource type id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"Policy apis"}

	ws.Route(ws.PATCH("/clusters/policy").To(ModifyPolicyByAlertCluster).
		Doc("Modify Policy By Alert Cluster level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/policy").To(ModifyPolicyByAlertNode).
		Doc("Modify Policy By Alert Node Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/workspaces/policy").To(ModifyPolicyByAlertWorkspace).
		Doc("Modify Policy By Alert Workspace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/workspaces/{ws_name}/policy").To(ModifyPolicyByAlertWorkspace).
		Doc("Modify Policy By Alert Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/policy").To(ModifyPolicyByAlertNamespace).
		Doc("Modify Policy By Alert Namespace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/policy").To(ModifyPolicyByAlertNamespace).
		Doc("Modify Policy By Alert Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/workloads/policy").To(ModifyPolicyByAlertWorkload).
		Doc("Modify Policy By Alert Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/pods/policy").To(ModifyPolicyByAlertPod).
		Doc("Modify Policy By Alert Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/{node_id}/pods/policy").To(ModifyPolicyByAlertPod).
		Doc("Modify Policy By Alert Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/pods/{pod_name}/containers/policy").To(ModifyPolicyByAlertContainer).
		Doc("Modify Policy By Alert Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/{node_id}/pods/{pod_name}/containers/policy").To(ModifyPolicyByAlertContainer).
		Doc("Modify Policy By Alert Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"Alert apis"}

	ws.Route(ws.POST("/clusters/alert").To(CreateAlertCluster).
		Doc("Create Alert Cluster level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/nodes/alert").To(CreateAlertNode).
		Doc("Create Alert Node Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/workspaces/alert").To(CreateAlertWorkspace).
		Doc("Create Alert Workspace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/workspaces/{ws_name}/alert").To(CreateAlertWorkspace).
		Doc("Create Alert Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/namespaces/alert").To(CreateAlertNamespace).
		Doc("Create Alert Namespace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/namespaces/{ns_name}/alert").To(CreateAlertNamespace).
		Doc("Create Alert Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/namespaces/{ns_name}/workloads/alert").To(CreateAlertWorkload).
		Doc("Create Alert Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/namespaces/{ns_name}/pods/alert").To(CreateAlertPod).
		Doc("Create Alert Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/nodes/{node_id}/pods/alert").To(CreateAlertPod).
		Doc("Create Alert Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/namespaces/{ns_name}/pods/{pod_name}/containers/alert").To(CreateAlertContainer).
		Doc("Create Alert Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/nodes/{node_id}/pods/{pod_name}/containers/alert").To(CreateAlertContainer).
		Doc("Create Alert Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/clusters/alert").To(ModifyAlertByNameCluster).
		Doc("Modify Alert By Name Cluster level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/alert").To(ModifyAlertByNameNode).
		Doc("Modify Alert By Name Node Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/workspaces/alert").To(ModifyAlertByNameWorkspace).
		Doc("Modify Alert By Name Workspace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/workspaces/{ws_name}/alert").To(ModifyAlertByNameWorkspace).
		Doc("Modify Alert By Name Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/alert").To(ModifyAlertByNameNamespace).
		Doc("Modify Alert By Name Namespace Level").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/alert").To(ModifyAlertByNameNamespace).
		Doc("Modify Alert By Name Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/workloads/alert").To(ModifyAlertByNameWorkload).
		Doc("Modify Alert By Name Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/pods/alert").To(ModifyAlertByNamePod).
		Doc("Modify Alert By Name Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/{node_id}/pods/alert").To(ModifyAlertByNamePod).
		Doc("Modify Alert By Name Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/namespaces/{ns_name}/pods/{pod_name}/containers/alert").To(ModifyAlertByNameContainer).
		Doc("Modify Alert By Name Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PATCH("/nodes/{node_id}/pods/{pod_name}/containers/alert").To(ModifyAlertByNameContainer).
		Doc("Modify Alert By Name Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/clusters/alert").To(DeleteAlertsByNameCluster).
		Doc("Delete Alerts By Name Cluster level").
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/nodes/alert").To(DeleteAlertsByNameNode).
		Doc("Delete Alerts By Name Node Level").
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/workspaces/alert").To(DeleteAlertsByNameWorkspace).
		Doc("Delete Alerts By Name Workspace Level").
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/workspaces/{ws_name}/alert").To(DeleteAlertsByNameWorkspace).
		Doc("Delete Alerts By Name Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/namespaces/alert").To(DeleteAlertsByNameNamespace).
		Doc("Delete Alerts By Name Namespace Level").
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/namespaces/{ns_name}/alert").To(DeleteAlertsByNameNamespace).
		Doc("Delete Alerts By Name Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/namespaces/{ns_name}/workloads/alert").To(DeleteAlertsByNameWorkload).
		Doc("Delete Alerts By Name Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/namespaces/{ns_name}/pods/alert").To(DeleteAlertsByNamePod).
		Doc("Delete Alerts By Name Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/nodes/{node_id}/pods/alert").To(DeleteAlertsByNamePod).
		Doc("Delete Alerts By Name Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/namespaces/{ns_name}/pods/{pod_name}/containers/alert").To(DeleteAlertsByNameContainer).
		Doc("Delete Alerts By Name Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.DELETE("/nodes/{node_id}/pods/{pod_name}/containers/alert").To(DeleteAlertsByNameContainer).
		Doc("Delete Alerts By Name Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/clusters/alert").To(DescribeAlertDetailsCluster).
		Doc("Describe Alert Details Cluster level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/alert").To(DescribeAlertDetailsNode).
		Doc("Describe Alert Details Node Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/alert").To(DescribeAlertDetailsWorkspace).
		Doc("Describe Alert Details Workspace Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{ws_name}/alert").To(DescribeAlertDetailsWorkspace).
		Doc("Describe Alert Details Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/alert").To(DescribeAlertDetailsNamespace).
		Doc("Describe Alert Details Namespace Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/alert").To(DescribeAlertDetailsNamespace).
		Doc("Describe Alert Details Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/workloads/alert").To(DescribeAlertDetailsWorkload).
		Doc("Describe Alert Details Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/alert").To(DescribeAlertDetailsPod).
		Doc("Describe Alert Details Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/alert").To(DescribeAlertDetailsPod).
		Doc("Describe Alert Details Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/{pod_name}/containers/alert").To(DescribeAlertDetailsContainer).
		Doc("Describe Alert Details Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/{pod_name}/containers/alert").To(DescribeAlertDetailsContainer).
		Doc("Describe Alert Details Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/clusters/alert_status").To(DescribeAlertStatusCluster).
		Doc("Describe Alert Status Cluster level").
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/alert_status").To(DescribeAlertStatusNode).
		Doc("Describe Alert Status Node Level").
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/alert_status").To(DescribeAlertStatusWorkspace).
		Doc("Describe Alert Status Workspace Level").
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{ws_name}/alert_status").To(DescribeAlertStatusWorkspace).
		Doc("Describe Alert Status Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/alert_status").To(DescribeAlertStatusNamespace).
		Doc("Describe Alert Status Namespace Level").
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/alert_status").To(DescribeAlertStatusNamespace).
		Doc("Describe Alert Status Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/workloads/alert_status").To(DescribeAlertStatusWorkload).
		Doc("Describe Alert Status Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/alert_status").To(DescribeAlertStatusPod).
		Doc("Describe Alert Status Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/alert_status").To(DescribeAlertStatusPod).
		Doc("Describe Alert Status Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/{pod_name}/containers/alert_status").To(DescribeAlertStatusContainer).
		Doc("Describe Alert Status Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/{pod_name}/containers/alert_status").To(DescribeAlertStatusContainer).
		Doc("Describe Alert Status Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("alert_ids", "alert id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("disables", "alert disabled list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("running_status", "alert running status list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("policy_ids", "alert policy id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("creators", "policy creator list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rs_filter_ids", "alert resource filter id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("executor_ids", "alert executor id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"History apis"}

	ws.Route(ws.GET("/clusters/history").To(DescribeHistoryDetailCluster).
		Doc("Describe History Detail Cluster level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/history").To(DescribeHistoryDetailNode).
		Doc("Describe History Detail Node Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/history").To(DescribeHistoryDetailWorkspace).
		Doc("Describe History Detail Workspace Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{ws_name}/history").To(DescribeHistoryDetailWorkspace).
		Doc("Describe History Detail Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/history").To(DescribeHistoryDetailNamespace).
		Doc("Describe History Detail Namespace Level").
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/history").To(DescribeHistoryDetailNamespace).
		Doc("Describe History Detail Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/workloads/history").To(DescribeHistoryDetailWorkload).
		Doc("Describe History Detail Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/history").To(DescribeHistoryDetailPod).
		Doc("Describe History Detail Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/history").To(DescribeHistoryDetailPod).
		Doc("Describe History Detail Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/{pod_name}/containers/history").To(DescribeHistoryDetailContainer).
		Doc("Describe History Detail Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/{pod_name}/containers/history").To(DescribeHistoryDetailContainer).
		Doc("Describe History Detail Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("search_word", "search word specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_names", "history name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("alert_names", "alert name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_names", "rule name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("events", "history event list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("rule_ids", "rule id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("resource_names", "resource name list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("recent", "recent history specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"Comment apis"}

	ws.Route(ws.POST("/comment").To(CreateComment).
		Doc("Create Comment").
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/comment").To(DescribeComments).
		Doc("Describe Comments").
		Param(ws.QueryParameter("comment_ids", "comment id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("addressers", "addresser list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("contents", "content list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("history_ids", "history id list specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_key", "sort key specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("reverse", "order specify").DataType("bool").DefaultValue("false").Required(false)).
		Param(ws.QueryParameter("offset", "record set offset specify").DataType("uint32").Required(false)).
		Param(ws.QueryParameter("limit", "record set limit specify").DataType("uint32").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	tags = []string{"Resource apis"}

	ws.Route(ws.GET("/clusters/resource").To(DescribeResourcesCluster).
		Doc("Describe Resources Cluster level").
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/resource").To(DescribeResourcesNode).
		Doc("Describe Resources Node Level").
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/resource").To(DescribeResourcesWorkspace).
		Doc("Describe Resources Workspace Level").
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{ws_name}/resource").To(DescribeResourcesWorkspace).
		Doc("Describe Resources Workspace Level").
		Param(ws.PathParameter("ws_name", "specific workspace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/resource").To(DescribeResourcesNamespace).
		Doc("Describe Resources Namespace Level").
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/resource").To(DescribeResourcesNamespace).
		Doc("Describe Resources Namespace Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/workloads/resource").To(DescribeResourcesWorkload).
		Doc("Describe Resources Workload Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("workload_kind", "workload kind specify").DataType("string").Required(false)).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/resource").To(DescribeResourcesPod).
		Doc("Describe Resources Pod Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/resource").To(DescribeResourcesPod).
		Doc("Describe Resources Pod Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{ns_name}/pods/{pod_name}/containers/resource").To(DescribeResourcesContainer).
		Doc("Describe Resources Container Level").
		Param(ws.PathParameter("ns_name", "specific namespace").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/nodes/{node_id}/pods/{pod_name}/containers/resource").To(DescribeResourcesContainer).
		Doc("Describe Resources Container Level").
		Param(ws.PathParameter("node_id", "specific node id").DataType("string").Required(true).DefaultValue("")).
		Param(ws.PathParameter("pod_name", "specific pod").DataType("string").Required(true).DefaultValue("")).
		Param(ws.QueryParameter("selector", "selector specify").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, tags)).
		Consumes(restful.MIME_JSON, constants.MIME_MERGEPATCH).
		Produces(restful.MIME_JSON)

	return ws
}

func Run() {
	restful.DefaultContainer.Add(WebService())
	enableCORS()

	global.GetInstance()

	cfg := config.GetInstance()
	apiPort, _ := strconv.Atoi(cfg.App.ApiPort)
	listen := fmt.Sprintf(":%d", apiPort)

	logger.Info(nil, "%+v", http.ListenAndServe(listen, nil))
}

func enableCORS() {
	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		CookiesAllowed: false,
		AllowedDomains: []string{"*"},
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)
}
