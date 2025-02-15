package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "a110f4f4-bb11-4532-b1f2-cd0e17e3dec6"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "babe0013",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Test 2: Confirm NIC exists and is connected to VM
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID), "NIC does not exist")
	nic := azure.GetNetworkInterface(t, nicName, resourceGroupName, subscriptionID)
	assert.Equal(t, vmName, *nic.VirtualMachine.ID, "NIC is not attached to the correct VM")

	// Test 3: Confirm the VM is running the correct Ubuntu version
	expectedUbuntuVersion := "22.04" // Change this based on the Terraform configuration
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, *vm.StorageProfile.ImageReference.Offer, "Ubuntu", "VM is not running Ubuntu")
	assert.Contains(t, *vm.StorageProfile.ImageReference.Sku, expectedUbuntuVersion, "VM is not running the expected Ubuntu version")
}