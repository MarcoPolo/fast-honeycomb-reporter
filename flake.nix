{
  description = "Internet speed and latency reporter using fast.com and 1.1.1.1";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-20.09";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.flake-utils.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          defaultPackage = import ./fast-honeycomb-reporter.nix { inherit pkgs; };
        }) // {
      nixosModule = { config, pkgs, ... }:
        let cfg = config.services.fast-honeycomb-reporter;
        in
        {
          options.services.fast-honeycomb-reporter = {
            apiKey = nixpkgs.lib.mkOption {
              description = "Honeycomb API Key";
              type = nixpkgs.lib.types.str;
            };
          };
          config = {
            users.users.fast-honeycomb-reporter.isSystemUser = true;
            systemd.services.fast-honeycomb-reporter = {
              description = "fast-honeycomb-reporter";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];
              serviceConfig = {
                Environment = "HONEYCOMB_API_KEY=${cfg.apiKey}";
                ExecStart = "${
                      self.defaultPackage.${pkgs.system}
                    }/bin/fast-honeycomb-reporter";
                Restart = "on-failure";
                User = "fast-honeycomb-reporter";
              };
            };
          };
        };

    };
}
