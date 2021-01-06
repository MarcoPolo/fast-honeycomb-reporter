{pkgs ? import <nixpkgs> {}}:
pkgs.buildGoModule rec {
  pname = "fast-honeycomb-reporter";
  version = "0.0.1";

  src = ./.;
  vendorSha256 = "sha256:P3Iyb0dXBIXTHT50I1FdQx5V8Zg/jznyX7q4SPoMCH8=";

  subPackages = [ "." ];

  meta = with pkgs.lib; {
    description = "Simple command-line snippet manager, written in Go";
    homepage = https://github.com/marcopolo/fast-honeycomb-reporter;
    license = licenses.mit;
    platforms = platforms.linux ++ platforms.darwin;
  };
}