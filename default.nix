{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  ),
  pname? "bard-cli",
  pversion ? "0.1"
}:

pkgs.buildGoApplication {
  pname = pname;
  version = pversion;
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
}
