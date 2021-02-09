# Runtime Env
`HONEYCOMB_API_KEY`: Needs to be set to the honeycomb api key

# NixOS Flake Module

Inside `flake.nix`
```nix

  # Add this as a dependency
  inputs.fast-honeycomb-reporter.url = "github:marcopolo/fast-honeycomb-reporter";
  # Optional: If you want this to follow your existing nixpkgs pin
  # inputs.fast-honeycomb-reporter.inputs.nixpkgs.follows = "nixpkgs";
  # inputs.fast-honeycomb-reporter.inputs.flake-utils.follows = "flake-utils";

  # Inside your nixos system config:
  # ...
  {
    imports = [ fast-honeycomb-reporter.nixosModule ];
    services.fast-honeycomb-reporter.apiKey = "HONEYCOMB_API_KEY";
  }
  # ...

```


Complete `flake.nix` example:
```nix
{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-20.09";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  inputs.fast-honeycomb-reporter.url = "github:marcopolo/fast-honeycomb-reporter";
  inputs.fast-honeycomb-reporter.inputs.nixpkgs.follows = "nixpkgs";
  inputs.fast-honeycomb-reporter.inputs.flake-utils.follows = "flake-utils";

  outputs =
    { self
    , nixpkgs
    , fast-honeycomb-reporter
    }: {
      nixosConfigurations = {
        pi4 = nixpkgs.lib.nixosSystem
          {
            system = "aarch64-linux";
            modules = [
              ./pi4_config.nix
              {
                imports = [ fast-honeycomb-reporter.nixosModule ];
                services.fast-honeycomb-reporter.apiKey =
                  (import ./secrets.nix).honeycomb_api_key;
              }
            ];
          };
      };
    };
}

```



# Building for ARM (raspberry pi 4)
```
GOOS=linux GOARCH=arm GOARM=7 go build
```