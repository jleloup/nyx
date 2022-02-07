nyx {
  changelog {
    path = 'CHANGELOG.md'
    sections = [
      'Added' : '^feat$',
      'Fixed' : '^fix$',
    ]
    substitutions = [
      '(?m)#([0-9]+)(?s)': '[#%s](https://example.com/issues/%s)'
    ]
  }
  commitMessageConventions {
    enabled = [ 'conventionalCommits' ]
    items {
      conventionalCommits {
        expression = '(?m)^(?<type>[a-zA-Z0-9_]+)(!)?(\\\\((?<scope>[a-z ]+)\\\\))?:( (?<title>.+))$(?s).*'
        bumpExpressions {
          major = '(?s)(?m)^[a-zA-Z0-9_]+(!|.*^(BREAKING( |-)CHANGE: )).*'
          minor = '(?s)(?m)^feat(?!!|.*^(BREAKING( |-)CHANGE: )).*'
          patch = '(?s)(?m)^fix(?!!|.*^(BREAKING( |-)CHANGE: )).*'
        }
      }
    }
  }
  configurationFile = '.nyx.json'
  dryRun = false
  git {
    remotes {
      origin {
        user = 'jdoe'
        password = 'somepassword'
      }
      replica {
        user = 'stiger'
        password = 'somesecret'
      }
    }
  }
  initialVersion = '1.0.0'
  preset = 'extended'
  releaseLenient = true
  releasePrefix = 'v'
  releaseTypes {
    enabled = [ 'mainline', 'maturity', 'integration', 'hotfix', 'feature', 'release', 'maintenance', 'internal' ]
    publicationServices = [ 'github', 'gitlab' ]
    remoteRepositories = [ 'origin', 'replica' ]
    items {
      mainline {
        collapseVersions = false
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)$'
        gitCommit = 'false'
        gitCommitMessage = 'Release version {{version}}'
        gitPush = 'true'
        gitTag = 'true'
        gitTagMessage = 'Tag version {{version}}'
        matchBranches = '^(master|main)$'
        matchEnvironmentVariables {
          CI = '^true$'
        }
        matchWorkspaceStatus = 'CLEAN'
        publish = 'true'
        versionRangeFromBranchName = false
      }
      maturity {
        collapseVersions = true
        collapsedVersionQualifier = '{{#sanitizeLower}}{{branch}}{{/sanitizeLower}}'
        description = 'Maturity release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)(-(alpha|beta|gamma|delta|epsilon|zeta|eta|theta|iota|kappa|lambda|mu|nu|xi|omicron|pi|rho|sigma|tau|upsilon|phi|chi|psi|omega)(\\.([0-9]\\d*))?)?$'
        gitCommit = 'false'
        gitPush = 'true'
        gitTag = 'true'
        matchBranches = '^(alpha|beta|gamma|delta|epsilon|zeta|eta|theta|iota|kappa|lambda|mu|nu|xi|omicron|pi|rho|sigma|tau|upsilon|phi|chi|psi|omega)$'
        matchWorkspaceStatus = 'CLEAN'
        publish = 'true'
        versionRangeFromBranchName = false
      }
      integration {
        collapseVersions = true
        collapsedVersionQualifier = '{{#sanitizeLower}}{{branch}}{{/sanitizeLower}}'
        description = 'Integration release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)(-(develop|development|integration|latest)(\\.([0-9]\\d*))?)$'
        gitCommit = 'false'
        gitPush = 'true'
        gitTag = 'true'
        matchBranches = '^(develop|development|integration|latest)$'
        matchWorkspaceStatus = 'CLEAN'
        publish = 'true'
        versionRangeFromBranchName = false
      }
      hotfix {
        collapseVersions = true
        collapsedVersionQualifier = '{{#sanitizeLower}}{{branch}}{{/sanitizeLower}}'
        description = 'Hotfix release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)(-(fix|hotfix)(([0-9a-zA-Z]*)(\\.([0-9]\\d*))?)?)$'
        gitCommit = 'false'
        gitPush = 'true'
        gitTag = 'true'
        matchBranches = '^(fix|hotfix)((-|\\/)[0-9a-zA-Z-_]+)?$'
        matchWorkspaceStatus = 'CLEAN'
        publish = 'true'
        versionRangeFromBranchName = false
      }
      feature {
        collapseVersions = true
        collapsedVersionQualifier = '{{#sanitizeLower}}{{branch}}{{/sanitizeLower}}'
        description = 'Feature release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)(-(feat|feature)(([0-9a-zA-Z]*)(\\.([0-9]\\d*))?)?)$'
        gitCommit = 'false'
        gitPush = 'false'
        gitTag = 'false'
        matchBranches = '^(feat|feature)((-|\\/)[0-9a-zA-Z-_]+)?$'
        publish = 'false'
        versionRangeFromBranchName = false
      }
      release {
        collapseVersions = true
        collapsedVersionQualifier = '{{#sanitizeLower}}{{branch}}{{/sanitizeLower}}'
        description = 'Release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)(-(rel|release)((\\.([0-9]\\d*))?)?)$'
        gitCommit = 'false'
        gitPush = 'true'
        gitTag = 'true'
        matchBranches = '^(rel|release)(-|\\/)({{configuration.releasePrefix}})?([0-9|x]\\d*)(\\.([0-9|x]\\d*)(\\.([0-9|x]\\d*))?)?$'
        matchWorkspaceStatus = 'CLEAN'
        publish = 'false'
        versionRangeFromBranchName = true
      }
      maintenance {
        collapseVersions = false
        description = 'Maintenance release {{version}}'
        filterTags = '^({{configuration.releasePrefix}})?([0-9]\\d*)\\.([0-9]\\d*)\\.([0-9]\\d*)$'
        gitCommit = 'false'
        gitPush = 'true'
        gitTag = 'true'
        matchBranches = '^[a-zA-Z]*([0-9|x]\\d*)(\\.([0-9|x]\\d*)(\\.([0-9|x]\\d*))?)?$'
        matchWorkspaceStatus = 'CLEAN'
        publish = 'true'
        versionRangeFromBranchName = true
      }
      internal {
        collapseVersions = true
        collapsedVersionQualifier = 'internal'
        description = 'Internal release {{version}}'
        gitCommit = 'false'
        gitPush = 'false'
        gitTag = 'false'
        identifiers {
          '0' {
            position = 'BUILD'
            qualifier = 'branch'
            value = '{{#sanitize}}{{branch}}{{/sanitize}}'
          }
          '1' {
            position = 'BUILD'
            qualifier = 'commit'
            value = '{{#short7}}{{releaseScope.finalCommit}}{{/short7}}'
          }
          '2' {
            position = 'BUILD'
            qualifier = 'timestamp'
            value = '{{#timestampYYYYMMDDHHMMSS}}{{timestamp}}{{/timestampYYYYMMDDHHMMSS}}'
          }
          '3' {
            position = 'BUILD'
            qualifier = 'user'
            value = '{{#sanitizeLower}}{{environment.user}}{{/sanitizeLower}}'
          }
        }
        publish = 'false'
        versionRangeFromBranchName = false
      }
    }
  }
  resume = true
  scheme = 'SEMVER'
  services {
    github {
      type = 'GITHUB'
      options {
        AUTHENTICATION_TOKEN = '{{#environment.variable}}GITHUB_TOKEN{{/environment.variable}}'
        REPOSITORY_NAME = 'myrepo'
        REPOSITORY_OWNER = 'acme'
      }
    }
    gitlab {
      type = 'GITLAB'
      options {
        AUTHENTICATION_TOKEN = 'abcde-a1b2c3d4e5f6g7h8i9j0'
        REPOSITORY_NAME = 'myrepo'
        REPOSITORY_OWNER = 'acme'
      }
    }
  }
  sharedConfigurationFile = 'example-shared.config.json'
  stateFile = '.nyx-state.yml'
  verbosity = 'INFO'
  version = '1.8.12'
}