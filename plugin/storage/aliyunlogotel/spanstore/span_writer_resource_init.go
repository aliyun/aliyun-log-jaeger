// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanstore

import (
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"go.uber.org/zap"
)

var indexJSON = `
{"keys":{"duration":{"alias":"","doc_value":true,"type":"long"},"operationName":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[","," ","'","\"",";","=","(",")","[","]","{","}","?","@","&","<",">","/",":","\n","\t","\r"],"type":"text"},"parentSpanID":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"process.serviceName":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[","," ","'","\"",";","=","(",")","[","]","{","}","?","@","&","<",">","/",":","\n","\t","\r"],"type":"text"},"process.tags.ip":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"spanID":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[","," ","'","\"",";","=","(",")","[","]","{","}","?","@","&","<",">","/",":","\n","\t","\r"],"type":"text"},"startTime":{"alias":"","doc_value":true,"type":"long"},"tags.component":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.downstream_cluster":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.guid:x-request-id":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.http.method":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.http.protocol":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.http.status_code":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.http.url":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.node_id":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.request_size":{"alias":"","doc_value":true,"type":"long"},"tags.response_flags":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.response_size":{"alias":"","doc_value":true,"type":"long"},"tags.span.kind":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.upstream_cluster":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"tags.user_agent":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"},"traceID":{"alias":"","caseSensitive":false,"chn":false,"doc_value":true,"token":[],"type":"text"}},"line":{"caseSensitive":false,"chn":false,"token":[","," ","'","\"",";","=","(",")","[","]","{","}","?","@","&","<",">","/",":","\n","\t","\r"]}}
`

var dashboardJSON = `
{"charts":[{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"%","description":"当前小时/同比昨日","descriptionSize":12,"displayName":"外部请求平均延迟","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":3,"maxval":100,"numberType":"compare","showTitle":false,"threshold":0,"unit":"ms","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":2.50,"yAxis":["_col1"],"yPos":0},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000, 2) as pv from log ))","start":"-3600s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000, 2) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697384000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"%","description":"当前小时/同比昨天","descriptionSize":12,"displayName":"外部请求P90延迟","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":3,"maxval":100,"numberType":"compare","showTitle":false,"threshold":0,"unit":"ms","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":5,"yAxis":["_col1"],"yPos":0},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(approx_percentile(duration, 0.90) / 1000000.0, 2) as pv from log ))","start":"-3600s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(approx_percentile(duration, 0.90) / 1000000.0, 2) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697450000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"%","description":"当前小时/同比昨天","descriptionSize":12,"displayName":"外部请求P99延迟","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":3,"maxval":100,"numberType":"compare","showTitle":false,"threshold":0,"unit":"ms","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":7.50,"yAxis":["_col1"],"yPos":0},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(approx_percentile(duration, 0.99) / 1000000.0, 2) as pv from log ))","start":"-3600s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(approx_percentile(duration, 0.99) / 1000000.0, 2) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697473000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"%","description":"PV/环比昨天","descriptionSize":12,"displayName":"日PV","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":3,"maxval":100,"numberType":"compare","showTitle":false,"threshold":0,"unit":"","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":0,"yAxis":["_col1"],"yPos":0},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  count(1) as pv from log ))","start":"-86400s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT  process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  count(1) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697539000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"","description":"请求延迟（15分钟）","descriptionSize":12,"displayName":"内部组件平均延迟","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":5,"maxval":100,"numberType":"dashboard","showTitle":true,"threshold":0,"unit":"","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":2.50,"yAxis":["_col1"],"yPos":3},"search":{"end":"now","logstore":"tracing","query":"not parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000.0, 2) as pv from log ))","start":"-900s","tokenQuery":"not parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000.0, 2) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697848000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"","description":"请求延迟（15分钟）","descriptionSize":12,"displayName":"外部请求平均延迟","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":5,"maxval":200,"numberType":"dashboard","showTitle":true,"threshold":0,"unit":"ms","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":0,"yAxis":["_col1"],"yPos":3},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000.0, 2) as pv from log ))","start":"-900s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(avg(duration) / 1000000.0, 2) as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542697919000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[3,10,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"","description":"4XX请求占比","descriptionSize":12,"displayName":"4XX请求比例","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":5,"maxval":100,"numberType":"dashboard","showTitle":false,"threshold":0,"unit":"%","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":5,"yAxis":["_col1"],"yPos":3},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(sum(CASE when  cast(\"tags.http.status_code\" as bigint) >= 400 and cast(\"tags.http.status_code\" as bigint) < 500  then 1 else 0 END ) *  100.0 / count(1), 2)  as pv from log ))","start":"-900s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(sum(CASE when  cast(\"tags.http.status_code\" as bigint) >= 400 and cast(\"tags.http.status_code\" as bigint) < 500  then 1 else 0 END ) *  100.0 / count(1), 2)  as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542698066000","type":"number"},{"action":{"type":"disabled"},"display":{"bgColor":{"a":"1","b":"255","g":"255","r":"255"},"colorAreaNum":3,"colors":["rgba(35,143,255,1)","rgba(255,191,0,1)","rgba(245,34,45,1)"],"colorval":[33,66,100],"compareColors":[{"color":"rgba(245,34,45,1)"},{"color":"rgba(255,255,255,1)"},{"color":"rgba(55,214,122,1)"},{"color":"rgba(255,255,255,1)"}],"compareSize":16,"compareUnit":"","description":"5XX请求占比","descriptionSize":12,"displayName":"5XX请求占比","fontColor":{"a":"1","b":"85","g":"85","r":"85"},"fontSize":32,"height":5,"maxval":100,"numberType":"dashboard","showTitle":true,"threshold":0,"unit":"万分比","unitSize":14,"width":2.50,"xAxis":["_col0"],"xPos":7.50,"yAxis":["_col1"],"yPos":3},"search":{"end":"now","logstore":"tracing","query":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(sum(CASE when  cast(\"tags.http.status_code\" as bigint) < 500  then 0 else 1 END ) *  10000.0 / count(1), 2)  as pv from log ))","start":"-900s","tokenQuery":"parentSpanID:  0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | select diff[1], round(100.0 * (diff[1] - diff[2] ) / diff[2], 2) from( SELECT  compare( pv , 86400)  as diff from ( SELECT  round(sum(CASE when  cast(\"tags.http.status_code\" as bigint) < 500  then 0 else 1 END ) *  10000.0 / count(1), 2)  as pv from log ))","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542698774000","type":"number"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"内部组件平均延迟","height":5,"legendPosition":"top","margin":[94,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":5,"yAxis":["平均延迟(ms)"],"yPos":8},"search":{"end":"absolute","logstore":"tracing","query":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699473000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"内部组件P99延迟","height":5,"legendPosition":"top","margin":[94,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":5,"yAxis":["99延迟(ms)"],"yPos":13},"search":{"end":"absolute","logstore":"tracing","query":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699498000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"内部组件P90延迟","height":5,"legendPosition":"top","margin":[94,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":0,"yAxis":["90延迟(ms)"],"yPos":13},"search":{"end":"absolute","logstore":"tracing","query":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"not traceID:0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699522000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"服务组件P99延迟","height":5,"legendPosition":"top","margin":[100,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":5,"yAxis":["99延迟(ms)"],"yPos":23},"search":{"end":"absolute","logstore":"tracing","query":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699687000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"服务组件P90延迟","height":5,"legendPosition":"top","margin":[94,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":0,"yAxis":["90延迟(ms)"],"yPos":23},"search":{"end":"absolute","logstore":"tracing","query":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699720000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"服务组件平均延迟","height":5,"legendPosition":"top","margin":[99,"auto","auto","auto"],"width":5,"xAxis":["time"],"xPos":5,"yAxis":["平均延迟(ms)"],"yPos":18},"search":{"end":"absolute","logstore":"tracing","query":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699791000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"服务组件请求数","height":5,"legendPosition":"top","margin":[99,0,"auto","auto"],"width":5,"xAxis":["time"],"xPos":0,"yAxis":["请求数"],"yPos":18},"search":{"end":"absolute","logstore":"tracing","query":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699823000","type":"agg"},{"action":{},"display":{"aggAxis":["process.serviceName"],"chartType":"line","displayName":"内部组件请求数","height":5,"legendPosition":"top","margin":[99,3,"auto","auto"],"width":5,"xAxis":["time"],"xPos":0,"yAxis":["请求数"],"yPos":8},"search":{"end":"absolute","logstore":"tracing","query":"not parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","start":"-86400s","tokenQuery":"not parentSpanID : 0 NOT  process.serviceName: istio-mixer NOT    process.serviceName: istio-telemetry NOT    process.serviceName: istio-ingressgateway | SELECT avg(duration) / 1000000.0 as \"平均延迟(ms)\", approx_percentile(duration, 0.99) / 1000000.0 as \"99延迟(ms)\",  approx_percentile(duration, 0.9) / 1000000.0 as \"90延迟(ms)\", count(1) as \"请求数\", time_series(__time__, '1h', '%d-%H', '0') as time, \"process.serviceName\" group by time, \"process.serviceName\" order by time","tokens":[],"topic":""},"title":"sls-zc-test-hz-pub1542699883000","type":"agg"}],"displayName":"istio-tracing","dashboardName":"istio-tracing-dashboard"}
`

func retryCreateDashboard(logger *zap.Logger, client sls.ClientInterface, project, logstore, dashboardStr string) (err error) {
	time.Sleep(time.Second * 10)
	// vendor dashboard
	dashboardStr = strings.Replace(dashboardStr, `"logstore":"tracing"`, "\"logstore\":\""+logstore+"\"", -1)
	dashboardStr = strings.Replace(dashboardStr, `"dashboardName":"istio-tracing-dashboard"`, "\"dashboardName\":\""+project+logstore+"\"", -1)
	// create index, create index do not return error
	for i := 0; i < 10; i++ {
		err = client.CreateDashboardString(project, dashboardStr)
		if err != nil {
			// if IndexAlreadyExist, just return
			if clientError, ok := err.(*sls.Error); ok && strings.Contains(clientError.Message, "already exist") {
				logger.With(zap.String("project", project)).
					With(zap.String("logstore", logstore)).
					Info("dashboard already exist")
				return nil
			}
			time.Sleep(time.Second)
		} else {
			logger.With(zap.String("project", project)).With(zap.String("logstore", logstore)).Info("create dashboard success")
			break
		}
	}
	return err
}

func retryCreateIndex(logger *zap.Logger, client sls.ClientInterface, project, logstore, indexStr string) (err error) {
	// create index, create index do not return error
	createFlag := true
	for i := 0; i < 10; i++ {
		if createFlag {
			err = client.CreateIndexString(project, logstore, indexStr)
		} else {
			err = client.UpdateIndexString(project, logstore, indexStr)
		}
		if err != nil {
			// if IndexAlreadyExist, just return
			if clientError, ok := err.(*sls.Error); ok && clientError.Code == "IndexAlreadyExist" {
				logger.With(zap.String("project", project)).
					With(zap.String("logstore", logstore)).
					Info("index already exist, try update index")
				createFlag = false
				continue
			}
			time.Sleep(time.Second)
		} else {
			logger.With(zap.String("project", project)).With(zap.String("logstore", logstore)).Info("create or update index success")
			break
		}
	}
	return err
}

func makesureLogstoreExist(logger *zap.Logger, client sls.ClientInterface, project, logstore string, shardCount, lifeCycle int) (new bool, err error) {
	for i := 0; i < 5; i++ {
		if ok, err := client.CheckLogstoreExist(project, logstore); err != nil {
			time.Sleep(time.Millisecond * 100)
		} else {
			if ok {
				return false, nil
			}
			break
		}
	}
	ttl := 180
	if shardCount <= 0 {
		shardCount = 2
	}
	// @note max init shard count limit : 10
	if shardCount > 10 {
		shardCount = 10
	}
	for i := 0; i < 5; i++ {
		err = client.CreateLogStore(project, logstore, ttl, shardCount, true, 32)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
		} else {
			logger.With(zap.String("project", project)).With(zap.String("logstore", logstore)).Info("create logstore success")
			break
		}
	}
	if err != nil {
		return true, err
	}
	// after create logstore success, wait 1 sec
	time.Sleep(time.Second)
	return true, nil
}

func makesureProjectExist(logger *zap.Logger, client sls.ClientInterface, project string) error {
	ok := false
	var err error

	for i := 0; i < 5; i++ {
		if ok, err = client.CheckProjectExist(project); err != nil {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	if ok {
		return nil
	}
	for i := 0; i < 5; i++ {
		_, err = client.CreateProject(project, "istio log project, created by alibaba cloud jeager collector")
		if err != nil {
			time.Sleep(time.Millisecond * 100)
		} else {
			logger.With(zap.String("project", project)).Info("create project success")
			break
		}
	}
	return err
}

// InitSpanWriterLogstoreResource create project, logstore, index, dashboard for jeager collector
func InitSpanWriterLogstoreResource(client sls.ClientInterface, project string, logstore string, logger *zap.Logger) error {
	if err := makesureProjectExist(logger, client, project); err != nil {
		return err
	}
	if _, err := makesureLogstoreExist(logger, client, project, logstore, 2, 90); err != nil {
		return err
	}
	if err := retryCreateIndex(logger, client, project, logstore, indexJSON); err != nil {
		return err
	}
	if err := retryCreateDashboard(logger, client, project, logstore, dashboardJSON); err != nil {
		return err
	}
	return nil
}
