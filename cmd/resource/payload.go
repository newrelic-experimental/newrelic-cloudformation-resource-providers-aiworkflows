package resource

import (
   "fmt"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

//
// Generic, should be able to leave these as-is
//

type Payload struct {
   model  *Model
   models []interface{}
}

func NewPayload(m *Model) *Payload {
   return &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
}

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}

func (p *Payload) GetTags() map[string]string {
   return nil
}

func (p *Payload) HasTags() bool {
   return false
}

// TODO Move this to common

//
// These are API specific, must be configured per API
//

var typeName = "NewRelic::Observability::AIWorkflows"

var blank = " "

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Id: &s})
}

func (p *Payload) GetCaptureKeys(a model.Action) []string {
   return []string{"id"}
}

var emptyString = ""

func (p *Payload) GetGraphQLFragment() *string {
   return &emptyString
}

func (p *Payload) SetIdentifier(g *string) {
   p.model.Id = g
   log.Debugf("SetIdentifier: %s", *p.model.Id)
}

func (p *Payload) GetTagIdentifier() *string {
   return nil
}

func (p *Payload) GetIdentifier() *string {
   return p.model.Id
}

func (p *Payload) GetVariables() map[string]string {
   log.Debugf("GetVariables: enter: p.model.Variables: %+v", p.model.Variables)
   vars := make(map[string]string)
   if p.model.Variables != nil {
      for k, v := range p.model.Variables {
         vars[k] = v
      }
   }

   if p.model.Id != nil {
      vars["ID"] = *p.model.Id
   }

   if p.model.WorkflowData != nil {
      vars["FRAGMENT"] = *p.model.WorkflowData
   }

   lqf := ""
   if p.model.ListQueryFilter != nil {
      lqf = *p.model.ListQueryFilter
   }
   vars["LISTQUERYFILTER"] = lqf

   log.Debugf("GetVariables: exit: vars: %+v", vars)
   return vars
}

func (p *Payload) GetErrorKey() string {
   return ""
}

func (p *Payload) GetIdentifierKey(a model.Action) string {
   return "id"
}

func (p *Payload) NeedsPropagationDelay(a model.Action) bool {
   return true
}

func (p *Payload) GetCreateMutation() string {
   return `
mutation {
  aiWorkflowsCreateWorkflow(accountId: {{{ACCOUNTID}}},
    createWorkflowData: 
      {{{FRAGMENT}}}
    ) {
    errors {
      description
      type
    }
    workflow {
      id
    }
  }
}
`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  aiWorkflowsDeleteWorkflow(accountId: {{{ACCOUNTID}}}, deleteChannels: false, id: "{{{ID}}}") {
    errors {
      description
      type
    }
    id
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
  aiWorkflowsUpdateWorkflow(accountId: {{{ACCOUNTID}}}, deleteUnusedChannels: false, 
    updateWorkflowData: { id: "{{{ID}}}",
      {{{FRAGMENT}}}
    }) {
    errors {
      description
      type
    }
    workflow {
      id
    }
  }
}
`
}

func (p *Payload) GetReadQuery() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      aiWorkflows {
        workflows(filters: {id: "{{{ID}}}"}) {
          entities {
            id
          }
          totalCount
        }
      }
    }
  }
}
`
}

func (p *Payload) GetListQuery() string {
   return `{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      aiWorkflows {
        workflows(cursor: "{{{NEXTCURSOR}}}") {
          entities {
            id
          }
          totalCount
          nextCursor
        }
      }
    }
  }
}
`
}

func (p *Payload) GetListQueryNextCursor() string {
   return p.GetListQuery()
}
