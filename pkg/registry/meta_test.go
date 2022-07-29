/*
Copyright 2022 Upbound Inc.
*/

package registry

import (
	"io/ioutil"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"

	xptest "github.com/crossplane/crossplane-runtime/pkg/test"
)

func TestScrapeRepo(t *testing.T) {
	type args struct {
		config *ScrapeConfiguration
	}
	type want struct {
		err    error
		pmPath string
	}
	tests := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ScrapeAWSResources": {
			reason: "Should successfully scrape AWS resource metadata",
			args: args{
				config: &ScrapeConfiguration{
					RepoPath:      "testdata/aws/r",
					CodeXPath:     `//code[@class="language-terraform" or @class="language-hcl"]/text()`,
					PreludeXPath:  `//text()[contains(., "description") and contains(., "subcategory")]`,
					FieldDocXPath: `//ul/li//code[1]/text()`,
					ImportXPath:   `//code[@class="language-shell"]/text()`,
				},
			},
			want: want{
				pmPath: "testdata/aws/pm.yaml",
			},
		},
		"ScrapeAzureResources": {
			reason: "Should successfully scrape Azure resource metadata",
			args: args{
				config: &ScrapeConfiguration{
					RepoPath:      "testdata/azure/r",
					CodeXPath:     `//code[@class="language-terraform" or @class="language-hcl"]/text()`,
					PreludeXPath:  `//text()[contains(., "description") and contains(., "subcategory")]`,
					FieldDocXPath: `//ul/li//code[1]/text()`,
					ImportXPath:   `//code[@class="language-shell"]/text()`,
				},
			},
			want: want{
				pmPath: "testdata/azure/pm.yaml",
			},
		},
		"ScrapeGCPResources": {
			reason: "Should successfully scrape GCP resource metadata",
			args: args{
				config: &ScrapeConfiguration{
					RepoPath:      "testdata/gcp/r",
					CodeXPath:     `//code[@class="language-terraform" or @class="language-hcl"]/text()`,
					PreludeXPath:  `//text()[contains(., "description") and contains(., "subcategory")]`,
					FieldDocXPath: `//ul/li//code[1]/text()`,
					ImportXPath:   `//code[@class="language-shell"]/text()`,
				},
			},
			want: want{
				pmPath: "testdata/gcp/pm.yaml",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pm := NewProviderMetadata("test-provider")
			err := pm.ScrapeRepo(tc.args.config)
			if diff := cmp.Diff(tc.want.err, err, xptest.EquateErrors()); diff != "" {
				t.Errorf("\n%s\nScrapeRepo(error): -want, +got:\n%s", tc.reason, diff)
			}
			if err != nil {
				return
			}
			pmExpected := ProviderMetadata{}
			buff, err := ioutil.ReadFile(tc.want.pmPath)
			if err != nil {
				t.Errorf("Failed to load expected ProviderMetadata from file: %s", tc.want.pmPath)
			}
			if err := yaml.Unmarshal(buff, &pmExpected); err != nil {
				t.Errorf("Failed to unmarshal expected ProviderMetadata from file: %s", tc.want.pmPath)
			}
			if diff := cmp.Diff(&pmExpected, pm, cmpopts.IgnoreUnexported(fieldpath.Paved{})); diff != "" {
				t.Errorf("\n%s\nScrapeRepo(ProviderConfig): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}