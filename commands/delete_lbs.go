package commands

import (
	"github.com/pivotal-cf-experimental/bosh-bootloader/bosh"
	"github.com/pivotal-cf-experimental/bosh-bootloader/flags"
	"github.com/pivotal-cf-experimental/bosh-bootloader/storage"
)

type DeleteLBs struct {
	awsCredentialValidator    awsCredentialValidator
	availabilityZoneRetriever availabilityZoneRetriever
	certificateManager        certificateManager
	infrastructureManager     infrastructureManager
	logger                    logger
	boshCloudConfigurator     boshCloudConfigurator
	cloudConfigManager        cloudConfigManager
	boshClientProvider        boshClientProvider
	stateStore                stateStore
}

type cloudConfigManager interface {
	Update(cloudConfigInput bosh.CloudConfigInput, boshClient bosh.Client) error
}

type deleteLBsConfig struct {
	skipIfMissing bool
}

func NewDeleteLBs(awsCredentialValidator awsCredentialValidator, availabilityZoneRetriever availabilityZoneRetriever,
	certificateManager certificateManager, infrastructureManager infrastructureManager, logger logger,
	boshCloudConfigurator boshCloudConfigurator, cloudConfigManager cloudConfigManager,
	boshClientProvider boshClientProvider, stateStore stateStore,
) DeleteLBs {
	return DeleteLBs{
		awsCredentialValidator:    awsCredentialValidator,
		availabilityZoneRetriever: availabilityZoneRetriever,
		certificateManager:        certificateManager,
		infrastructureManager:     infrastructureManager,
		logger:                    logger,
		boshCloudConfigurator:     boshCloudConfigurator,
		cloudConfigManager:        cloudConfigManager,
		boshClientProvider:        boshClientProvider,
		stateStore:                stateStore,
	}
}

func (c DeleteLBs) Execute(subcommandFlags []string, state storage.State) error {
	config, err := c.parseFlags(subcommandFlags)
	if err != nil {
		return err
	}

	err = c.awsCredentialValidator.Validate()
	if err != nil {
		return err
	}

	if config.skipIfMissing && !lbExists(state.Stack.LBType) {
		c.logger.Println("no lb type exists, skipping...")
		return nil
	}

	if err := checkBBLAndLB(state, c.boshClientProvider, c.infrastructureManager); err != nil {
		return err
	}

	azs, err := c.availabilityZoneRetriever.Retrieve(state.AWS.Region)
	if err != nil {
		return err
	}

	stack, err := c.infrastructureManager.Describe(state.Stack.Name)
	if err != nil {
		return err
	}

	cloudConfigInput := c.boshCloudConfigurator.Configure(stack, azs)
	cloudConfigInput.LBs = nil

	boshClient := c.boshClientProvider.Client(state.BOSH.DirectorAddress, state.BOSH.DirectorUsername, state.BOSH.DirectorPassword)

	err = c.cloudConfigManager.Update(cloudConfigInput, boshClient)
	if err != nil {
		return err
	}

	_, err = c.infrastructureManager.Update(state.KeyPair.Name, len(azs), state.Stack.Name, "", "", state.EnvID)
	if err != nil {
		return err
	}

	c.logger.Step("deleting certificate")
	err = c.certificateManager.Delete(state.Stack.CertificateName)
	if err != nil {
		return err
	}

	state.Stack.LBType = "none"
	state.Stack.CertificateName = ""

	err = c.stateStore.Set(state)
	if err != nil {
		return err
	}

	return nil
}

func (DeleteLBs) parseFlags(subcommandFlags []string) (deleteLBsConfig, error) {
	lbFlags := flags.New("delete-lbs")

	config := deleteLBsConfig{}
	lbFlags.Bool(&config.skipIfMissing, "skip-if-missing", "", false)

	err := lbFlags.Parse(subcommandFlags)
	if err != nil {
		return config, err
	}

	return config, nil
}
