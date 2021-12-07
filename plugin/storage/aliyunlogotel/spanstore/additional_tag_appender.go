package spanstore

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var (
	logger, _ = zap.NewDevelopment()
)

const (
	OperationWithPrefix = "OperationWithPrefix"
	Tags                = "Tags"
)

type TagMappingValue struct {
	TagKey   string
	TagValue string
}

type TagAppendRules interface {
	SpanTagRules() map[string]*TagMappingValue
	OperationPrefixRules() map[string]*TagMappingValue
}

type tagRules struct {
	SpanTagRule        map[string]*TagMappingValue
	OperationNameRules map[string]*TagMappingValue
}

func (t tagRules) SpanTagRules() map[string]*TagMappingValue {
	return t.SpanTagRule
}

func (t tagRules) OperationPrefixRules() map[string]*TagMappingValue {
	return t.OperationNameRules
}

func initTagAppendRules(ruleFile string) TagAppendRules {
	var spanTagsAppendRules = map[string]*TagMappingValue{
		"db.instance": {TagKey: "db.system", TagValue: "Database"},
		"redis.key":   {TagKey: "db.system", TagValue: "Redis"},
	}

	var operationPrefixAppendRules = map[string]*TagMappingValue{
		"elastic-POST": {TagKey: "db.system", TagValue: "ElasticSearch"},
	}

	data := initTagAppenderRule(ruleFile)

	if d, ok := data[OperationWithPrefix]; ok {
		for k, v := range d {
			operationPrefixAppendRules[k] = v
		}
	}

	if d, ok := data[Tags]; ok {
		for k, v := range d {
			spanTagsAppendRules[k] = v
		}
	}

	d, _ := json.Marshal(spanTagsAppendRules)
	d2, _ := json.Marshal(operationPrefixAppendRules)
	logger.Info("The tag append rules.", zap.String("TagAppendRules", string(d)), zap.String("OperationNamePrefixAppendRules", string(d2)))

	return &tagRules{
		SpanTagRule:        spanTagsAppendRules,
		OperationNameRules: operationPrefixAppendRules,
	}
}

func initTagAppenderRule(tagAppenderRuleFile string) map[string]map[string]*TagMappingValue {
	if tagAppenderRuleFile != "" {
		if file, err := os.Open(tagAppenderRuleFile); err == nil {
			if data, e := ioutil.ReadAll(file); e == nil {
				logger.Info("The context of tag append rule file.", zap.String("content", string(data)))
				tagMapping := make(map[string]map[string]*TagMappingValue)
				if e1 := json.Unmarshal(data, &tagMapping); e1 == nil {
					return tagMapping
				} else {
					logger.Warn("Failed to pared the tag append rule.", zap.Error(e1))
				}
			} else {
				logger.Warn("Failed to read the tag append rule.", zap.Error(e))
			}
		} else {
			logger.Warn("the tag append rule file is not exist.", zap.Error(err))
		}
	}
	return nil
}
