{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
}:

let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  packages = [
    goEnv
    pkgs.gomod2nix
    pkgs.niv
  ];
}

## run in current shell :  `nix-shell -E '{ pkgs ?  ( let sources = import ./nix/sources.nix; in import sources.nixpkgs { overlays = [ (import "${sources.gomod2nix}/overlay.nix") ]; }) }: pkgs.mkShell { nativeBuildInputs =  [ (pkgs.callPackage ./default.nix {} )  ]; } '`
