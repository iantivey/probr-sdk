package connection

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/citihub/probr-sdk/utils"
)

type AzureManagedCluster struct {
	ctx                     context.Context
	credentials             AzureCredentials
	azManagedClustersClient containerservice.ManagedClustersClient
}

//var azConnection connection.Azure // Provides functionality to interact with Azure

// NewContainerService provides a new instance of AzureContainerService
func NewContainerService(c context.Context, creds AzureCredentials) (cs *AzureManagedCluster, err error) {

	// Guard clause - context
	if c == nil {
		err = utils.ReformatError("Context instance cannot be nil")
		return
	}

	// Guard clause - authorizer
	if creds.Authorizer == nil {
		err = utils.ReformatError("Authorizer instance cannot be nil")
		return
	}

	cs = &AzureManagedCluster{
		ctx:         c,
		credentials: creds,
	}

	// Create an azure storage account client object via the connection config vars
	var csErr error
	cs.azManagedClustersClient, csErr = cs.getManagedClusterClient(creds)
	if csErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Kubernetes Service client: %v", csErr)
		return
	}

	return
}

// GetJSONRepresentation returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (amc *AzureManagedCluster) GetJSONRepresentation(resourceGroupName string, clusterName string) (aksJSON []byte, err error) {
	var cs containerservice.ManagedCluster
	cs, err = amc.azManagedClustersClient.Get(amc.ctx, resourceGroupName, clusterName)
	if err != nil {
		log.Printf("Error getting ContainerServiceClient: %v", err)
		return
	}
	aksJSON, err = cs.MarshalJSON()
	return
}

func (amc *AzureManagedCluster) getManagedClusterClient(creds AzureCredentials) (csClient containerservice.ManagedClustersClient, err error) {

	log.Printf("Credentials: Subscription: %s", creds.SubscriptionID)

	csClient = containerservice.NewManagedClustersClient(creds.SubscriptionID)
	// Create an azure container services client object via the connection config vars

	csClient.Authorizer = creds.Authorizer

	return
}
