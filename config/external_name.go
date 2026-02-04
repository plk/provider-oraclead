package config

import (
	"errors"
	"fmt"
	"strings"
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Import requires full Azure ID
	"azurerm_oracle_autonomous_database": idFromName(),
}

func idFromName() config.ExternalName {
	e := config.NameAsIdentifier
	e.GetExternalNameFn = getNameFromFullyQualifiedID
	e.GetIDFn = getFullyQualifiedIDfunc
	return e
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}

func getNameFromFullyQualifiedID(tfstate map[string]any) (string, error) {
	id, ok := tfstate["id"]
	if !ok {
		return "", errors.New(fmt.Sprint("No attribute 'id'"))
	}
	idStr, ok := id.(string)
	if !ok {
		return "", errors.New(fmt.Sprint("Unexpected type 'id'"))
	}
	words := strings.Split(idStr, "/")
	return words[len(words)-1], nil
}

func getFullyQualifiedIDfunc(ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
	subID, ok := providerConfig["subscription_id"]
	if !ok {
		return "", errors.New(fmt.Sprint("No attribute 'subscription_id'"))
    }
    subIDStr, ok := subID.(string)
    if !ok {
		return "", errors.New(fmt.Sprint("Unexpected type for 'subscription_id'"))
    }
    rg, ok := parameters["resource_group_name"]
    if !ok {
		return "", errors.New(fmt.Sprint("No attribute 'resource_group_name'"))
    }
    rgStr, ok := rg.(string)
    if !ok {
		return "", errors.New(fmt.Sprint("Unexpected type for 'resource_group_name'"))
    }
	name, ok := parameters["name"]
    if !ok {
		return "", errors.New(fmt.Sprint("No attribute 'name'"))
    }
    nameStr, ok := name.(string)
    if !ok {
		return "", errors.New(fmt.Sprint("Unexpected type for 'name'"))
    }

    return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/autonomousDatabases/%s", subIDStr, rgStr, nameStr), nil
}
