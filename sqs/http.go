package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go-aws-sqs/sqs/option"
)

func optionsHttp(opt *option.OptionsHttp) func(options *sqs.Options) {
	return func(options *sqs.Options) {
		if opt == nil {
			return
		}
		options = &sqs.Options{
			APIOptions:                       opt.APIOptions,
			AppID:                            opt.AppID,
			BaseEndpoint:                     opt.BaseEndpoint,
			ClientLogMode:                    opt.ClientLogMode,
			Credentials:                      opt.Credentials,
			DefaultsMode:                     opt.DefaultsMode,
			DisableMessageChecksumValidation: opt.DisableMessageChecksumValidation,
			EndpointOptions:                  opt.EndpointOptions,
			EndpointResolverV2:               opt.EndpointResolverV2,
			HTTPSignerV4:                     opt.HTTPSignerV4,
			Logger:                           opt.Logger,
			Region:                           opt.Region,
			RetryMaxAttempts:                 opt.RetryMaxAttempts,
			RetryMode:                        opt.RetryMode,
			Retryer:                          opt.Retryer,
			RuntimeEnvironment:               opt.RuntimeEnvironment,
			HTTPClient:                       opt.HTTPClient,
			AuthSchemeResolver:               opt.AuthSchemeResolver,
			AuthSchemes:                      opt.AuthSchemes,
		}
	}
}
