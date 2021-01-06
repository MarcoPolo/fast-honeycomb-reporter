{
  description = "my project description";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system}; in
        {
          defaultPackage = import ./fast-honeycomb-reporter.nix { inherit pkgs; };
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
                users.users.honeycomb-reporter.isSystemUser = true;
                systemd.services.honeycomb-reporter = {
                  description = "fast-honeycomb-reporter";
                  wantedBy = [ "multi-user.target" ];
                  after = [ "network.target" ];
                  serviceConfig = {
                    ExecStart = "HONEYCOMB_API_KEY=${cfg.apiKey} ${self.defaultPackage.${system}}/bin/fast-honeycomb-reporter";
                    Restart = "on-failure";
                    User = "honeycomb-reporter";
                  };
                };
              };
            };
        }
      );
}
