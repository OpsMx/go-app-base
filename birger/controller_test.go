/*
 * Copyright 2022 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package birger

import (
	"reflect"
	"testing"
)

func Test_parseAgentStatistics(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		filter  []string
		want    map[string]controllerService
		wantErr bool
	}{
		{
			"empty errors",
			[]byte(""),
			[]string{"whoami"},
			map[string]controllerService{},
			true,
		}, {
			"One agent, one service, one annotation",
			[]byte(`
			{
				"serverTime": 1662067531436,
				"version": "v3.4.6-6-g4eee038",
				"connectedAgents": [
				  {
					"name": "smith",
					"session": "0001HH270W7TD8DZZ6STNY2ASX",
					"connectionType": "direct",
					"endpoints": [
					  {
						"name": "whoami",
						"type": "whoami",
						"configured": true,
						"annotations": {
						  "description": "demo service"
						}
					  }
					],
					"version": "v3.4.6-6-g4eee038",
					"hostname": "studio.local",
					"connectedAt": 1662065692965,
					"lastPing": 1662067522916,
					"agentInfo": {
					  "annotations": {
						"description": "demo agent"
					  }
					}
				  }
				]
			  }
			`),
			[]string{"whoami"},
			map[string]controllerService{
				"smith:whoami:whoami": {
					Name:      "whoami",
					Type:      "whoami",
					AgentName: "smith",
					Annotations: map[string]string{
						"description": "demo service",
					},
				},
			},
			false,
		}, {
			"most recently connected agent endpoints are used",
			[]byte(`
				{
					"serverTime": 1662067531436,
					"version": "v3.4.6-6-g4eee038",
					"connectedAgents": [
					  {
						"name": "smith",
						"session": "session-one",
						"connectionType": "direct",
						"endpoints": [
						  {
							"name": "whoami",
							"type": "whoami",
							"configured": true,
							"annotations": {
							  "description": "demo service",
							  "otherAnnotation": "newer annotation"
							}
						  }
						],
						"version": "v3.4.6-6-g4eee038",
						"hostname": "studio.local",
						"connectedAt": 999,
						"lastPing": 1662067522916,
						"agentInfo": {
						  "annotations": {
							"description": "demo agent"
						  }
						}
					  },
					  {
						"name": "smith",
						"session": "session-two",
						"connectionType": "direct",
						"endpoints": [
						  {
							"name": "whoami",
							"type": "whoami",
							"configured": true,
							"annotations": {
							  "description": "demo service",
							  "otherAnnotation": "very old annotation"
							}
						  }
						],
						"version": "v3.4.6-6-g4eee038",
						"hostname": "studio.local",
						"connectedAt": 111,
						"lastPing": 1662067522916,
						"agentInfo": {
						  "annotations": {
							"description": "demo agent"
						  }
						}
					  },
					  {
						"name": "smith",
						"session": "session-two",
						"connectionType": "direct",
						"endpoints": [
						  {
							"name": "whoami",
							"type": "whoami",
							"configured": true,
							"annotations": {
							  "description": "demo service",
							  "otherAnnotation": "old annotation"
							}
						  }
						],
						"version": "v3.4.6-6-g4eee038",
						"hostname": "studio.local",
						"connectedAt": 222,
						"lastPing": 1662067522916,
						"agentInfo": {
						  "annotations": {
							"description": "demo agent"
						  }
						}
					  }
					]
				  }
				`),
			[]string{"whoami"},
			map[string]controllerService{
				"smith:whoami:whoami": {
					Name:      "whoami",
					Type:      "whoami",
					AgentName: "smith",
					Annotations: map[string]string{
						"description":     "demo service",
						"otherAnnotation": "newer annotation",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ControllerManager{serviceTypes: tt.filter}
			got, err := m.parseAgentStatistics(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAgentStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAgentStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}
