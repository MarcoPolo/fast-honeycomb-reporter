{
  description = "my project description";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let pkgs = nixpkgs.legacyPackages.${system}; in
        {
          defaultPackage = import ./fast-honeycomb-reporter.nix {inherit pkgs;};
        }
      );
}