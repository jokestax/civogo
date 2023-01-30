package civogo

import (
	"reflect"
	"testing"
)

func TestListApplications(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `{"page":1,"per_page":20,"pages":1,"items":[{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-app-name",
		  "status": "ACTIVE",
		  "account_id": "12345",
		  "network_id": "34567",
		  "process_info": [
			  	{
					"processType": "web",
					"processCount": 1
				}],
			"domains": [
				"your-app-name.example.com"
			],
			}]}`,
	})

	defer server.Close()
	got, err := client.ListApplications()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedApplications{
		Page:    1,
		PerPage: 20,
		Pages:   1,
		Items: []Application{
			{
				ID:        "69a23478-a89e-41d2-97b1-6f4c341cee70",
				Name:      "your-app-name",
				Status:    "ACTIVE",
				NetworkID: "34567",
				ProcessInfo: []ProcInfo{
					{
						ProcessType:  "web",
						ProcessCount: 1,
					},
				},
				Domains: []string{"your-app-name.example.com"},
			},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestCreateApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `{"name":"test-app"}`,
	})
	defer server.Close()

	cfg := &ApplicationConfig{
		Name: "test-app",
		Size: "small",
	}

	got, err := client.CreateApplication(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "test-app" {
		t.Errorf("Expected %s, got %s", "test-app", got.Name)
	}
}

func TestDeleteApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications/12345": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.DeleteApplication("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
