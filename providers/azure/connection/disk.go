package connection

import (
	"context"
	"log"
	"strings"
	//"sync"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	//"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	//"github.com/Azure/go-autorest/autorest"
	//"github.com/Azure/go-autorest/autorest/azure/auth"
	//azconnection "github.com/citihub/probr-sdk/providers/azure/connection"
	"github.com/citihub/probr-sdk/utils"
)

type AzureDisk struct {
	ctx          context.Context
	credentials  AzureCredentials
	azDiskClient compute.DisksClient
}

func NewDisk(c context.Context, creds AzureCredentials) (dsk *AzureDisk, err error) {
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

	dsk = &AzureDisk{
		ctx:         c,
		credentials: creds,
	}

	// Create an azure storage account client object via the connection config vars
	var dskErr error
	dsk.azDiskClient, dskErr = dsk.getDisksClient(creds)
	if dskErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Disks client: %v", dskErr)
		return
	}

	return
}

func (dsk *AzureDisk) getDisksClient(creds AzureCredentials) (dskClient compute.DisksClient, err error) {

	log.Printf("Credentials: Subscription: %s", creds.SubscriptionID)

	dskClient = compute.NewDisksClient(creds.SubscriptionID)
	// Create an azure container services client object via the connection config vars

	dskClient.Authorizer = creds.Authorizer

	return
}

func (dsk *AzureDisk) GetDisk(resourceGroupName string, diskName string) (d compute.Disk, err error) {
	d, err = dsk.azDiskClient.Get(dsk.ctx, resourceGroupName, diskName)
	log.Printf("[DEBUG] GetDisk.d: %v", d)
	return
}

func (dsk *AzureDisk) ParseDiskDetails(diskURI string) (resourceGroupName, diskName string) {
	///subscriptions/b82de549-1f90-4ed0-9bbe-7632e52b7b3b/resourceGroups/mc_probr-demo-rg_probr-demo-cluster_eastus2/providers/Microsoft.Compute/disks/kubernetes-dynamic-pvc-17a4da65-a206-41af-a3e0-9661dfee9195
	s := strings.Split(diskURI, "/")
	resourceGroupName = s[4]
	diskName = s[8]
	return
}

// GetJSONRepresentation returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (dsk *AzureDisk) GetJSONRepresentation(resourceGroupName string, diskName string) (dskJSON []byte, err error) {
	var d compute.Disk
	d, err = dsk.azDiskClient.Get(dsk.ctx, resourceGroupName, diskName)
	if err != nil {
		log.Printf("Error getting ContainerServiceClient: %v", err)
		return
	}
	dskJSON, err = d.MarshalJSON()
	return
}
