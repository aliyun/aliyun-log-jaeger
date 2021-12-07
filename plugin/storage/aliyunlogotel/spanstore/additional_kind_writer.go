package spanstore

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type KindRewriteRules interface {
	SpanKindRules() map[string]string
	OperationPrefixRules() map[string]string
}

type kindRules struct {
	SpanTagRule        map[string]string
	OperationNameRules map[string]string
}

func (t kindRules) SpanKindRules() map[string]string {
	return t.SpanTagRule
}

func (t kindRules) OperationPrefixRules() map[string]string {
	return t.OperationNameRules
}

func initKindRewriteRules(ruleFile string) KindRewriteRules {
	var spanTagsRewriteRules = map[string]string{
		"db.instance": "client",
		"redis.key":   "client",
	}

	var operationPrefixRewriteRules = map[string]string{
		"elastic-POST": "client",
	}

	data := initKindRewriteRule(ruleFile)

	if d, ok := data[OperationWithPrefix]; ok {
		for k, v := range d {
			operationPrefixRewriteRules[k] = v
		}
	}

	if d, ok := data[Tags]; ok {
		for k, v := range d {
			spanTagsRewriteRules[k] = v
		}
	}

	d, _ := json.Marshal(spanTagsRewriteRules)
	d2, _ := json.Marshal(operationPrefixRewriteRules)
	logger.Info("The tag rewrite rules.", zap.String("KindRewriteRules", string(d)), zap.String("OperationNamePrefixAppendRules", string(d2)))

	return &kindRules{
		SpanTagRule:        spanTagsRewriteRules,
		OperationNameRules: operationPrefixRewriteRules,
	}
}

func initKindRewriteRule(kindAppenderRuleFile string) map[string]map[string]string {
	if kindAppenderRuleFile != "" {
		if file, err := os.Open(kindAppenderRuleFile); err == nil {
			if data, e := ioutil.ReadAll(file); e == nil {
				logger.Info("The context of kind rewrite rule file.", zap.String("content", string(data)))
				tagMapping := make(map[string]map[string]string)
				if e1 := json.Unmarshal(data, &tagMapping); e1 == nil {
					return tagMapping
				} else {
					logger.Warn("Failed to pared the kind rewrite rule.", zap.Error(e1))
				}
			} else {
				logger.Warn("Failed to read the kind rewrite rule.", zap.Error(e))
			}
		} else {
			logger.Warn("the kind rewrite rule file is not exist.", zap.Error(err))
		}
	}
	return nil
}
