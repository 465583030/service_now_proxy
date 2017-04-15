package main

import (
	"encoding/json"
	"log"
)


type IncidentResult struct {
	Incidents []Incident `json:"result"`
	Count     int `json:"incident_count"`
}

type SNLink struct {
	Link string `json:"link"`
	Value string `json:"value"`
}

type Incident struct {
	Number string `json:"number"`
	SysCreatedBy string `json:"sys_created_by"`
	UIncidentType string `json:"u_incident_type"`
	IncidentState string `json:"incident_state"`
	Impact string `json:"impact"`
	Active string `json:"active"`
	Priority string `json:"priority"`
	ShortDescription string `json:"short_description"`
	TicketID string `json:"sys_id"`
	ClosedBy string `json:"closed_by"`
	AssignedToRaw json.RawMessage `json:"assigned_to,omitempty"` //todo,demote this to unexported value
	ULsmCustomerImpacting string `json:"u_lsm_customer_impacting"`
	UResolvedOn string `json:"u_resolved_on"`
	UCategoryTier1 string `json:"u_category_tier_1"`
	SysUpdatedBy string `json:"sys_updated_by"`
	UCategoryTier3 string `json:"u_category_tier_3"`
	UCategoryTier2 string `json:"u_category_tier_2"`
	SysCreatedOn string `json:"sys_created_on"`
	USLA string `json:"u_sla"`
	AssignmentGroup SNLink `json:"assignment_group,omitempty"`
	Urgency string `json:"urgency"`
	Severity string `json:"severity"`
	LSMAssigned interface{} `json:"lsm_assigned"`
}

type IncidentParams struct {
	limit string
	active bool
	teamID string
	incidentID string
	query string
}


func (c client) Incidents(p IncidentParams) (IncidentResult){
	gp := make(map[string]string)
	if p.teamID != "" {
		gp["assignment_group"] = p.teamID
	}
	gp["sysparm_limit"] = p.limit
	if p.limit != "" {
		gp["sysparm_limit"] = "100"
	}
	if p.active {
		gp["sysparm_query"] = "active=true"
	}
	if p.incidentID != "" {
		gp["number"] = p.incidentID
	}

	if gp["assignment_group"] == "" && gp["number"] == "" {
		log.Fatal("either teamID or incidentID must be provided")
	}

	IncidentRequest := getParams{}
	IncidentRequest.params = gp
	IncidentRequest.path = INCIDENTLISTPATH
	IncidentRequest.Client = c
	return IncidentRequest.Get().IncidentsData()
}

func (i Incident) AssignedUser() User {
	if string(i.AssignedToRaw) != "" {
		user := SNLink{}
		err := json.Unmarshal(i.AssignedToRaw,&user)
		if err != nil {
			log.Printf("Could not parse Assigned to details, %+v",i.AssignedToRaw)
			var u = User{"N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A"}
			return u
		}
		userInfo := serviceNow.User(user.Value)
		if len(userInfo) > 0 {
			return userInfo[0]
		}
	}
	var u = User{"N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A","N/A"}
	return u
}

func (ir IncidentResult) DataPresent() bool {
	if ir.Count > 0 {
		return true
	}
	return false
}

func (rd returnData) IncidentsData() (res IncidentResult){
	err := json.Unmarshal(rd, &res)
	if err != nil {
		log.Printf("Could not unmarshall Incident response to struct - %+v\n",err)
		return
	}
	res.Count = len(res.Incidents)
	for index,incident := range res.Incidents {
		res.Incidents[index].LSMAssigned = incident.AssignedUser()
	}
        return
}