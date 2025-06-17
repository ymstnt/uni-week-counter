{
  lib,
  version,
  buildGoModule,
}:

buildGoModule {
  pname = "uni-week-counter";
  inherit version;

  src = ./.;

  vendorHash = null;

  meta = {
    description = "University week counter API.";
    homepage = "https://github.com/ymstnt/uni-week-counter/";
    license = lib.licenses.gpl3;
    maintainers = with lib.maintainers; [ ymstnt ];
    platforms = lib.platforms.all;
    mainProgram = "uni-week-counter";
  };
}
