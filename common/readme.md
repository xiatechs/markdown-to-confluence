# markdown-to-confluence/common readme

## the common package is for storing common constants used in app

### The package contains one exported variable
```
// ConstantsHardCoded - set to true if you are testing this locally & want to edit constants to actual parameters
ConstantsHardCoded = false
```

### The package contains multiple exported constants
```
// ConfluenceBaseURL is the base URL for the confluence page you want to work with
ConfluenceBaseURL     = "https://xiatech-markup.atlassian.net"
	
// ConfluenceUsernameEnv is to collect external env var INPUT_CONFLUENCE_USERNAME
ConfluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	
// ConfluenceAPIKeyEnv is to collect external env var INPUT_CONFLUENCE_API_KEY
ConfluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	
// ConfluenceSpaceEnv is to collect external env var INPUT_CONFLUENCE_SPACE
ConfluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
	
// EnvsNotSetError is template env var not set error
EnvsNotSetError       = "environment variable not set, please assign values for: "
```
