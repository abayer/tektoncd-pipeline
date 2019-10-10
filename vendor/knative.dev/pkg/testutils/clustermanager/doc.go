/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package clustermanager provides support for managing clusters for e2e tests,
responsible for creating/deleting cluster, and cluster life cycle management if
running in Prow
usage example:
func acquireCluster() {
    clusterOps := GKEClient{}.Setup(2, "n1-standard-8", "us-east1", "a", "myproject")
    // Cast to GKEOperation
    GKEOps := clusterOps.(GKECluster)
    if err = GKEOps.Acquire(); err != nil {
        log.Fatalf("Failed acquire cluster: '%v'", err)
    }
    log.Printf("GKE project is: %s", GKEOps.Project)
    log.Printf("GKE cluster is: %v", GKEOps.Cluster)
}
*/
package clustermanager
