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

// ServiceUpdate contains an update message sent when a new service type is
// discovered or is no longer present in the controller.
//
// Operation is either 'update' or 'delete'.  For both,
// Name, Type, and AgentName will be set.  For update only, the
// URL and Token will also be included.
type ServiceUpdate struct {
	Operation   string // delete, update (implies add)
	Name        string
	Type        string
	AgentName   string
	Annotations map[string]string // Only set for update
	Token       string            // Only set for update
	URL         string            // Only set for update
}
