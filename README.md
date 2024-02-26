# bsputil

just a toy project to parse lumps from an RBSP file and emit linewise-JSON of lump objects.

![Go](https://github.com/Razish/bsputil/actions/workflows/build.yml/badge.svg)

## Supported games/formats

- RBSP version 1 (Jedi Academy, Jedi Outcast)

## Supported lumps

- `ents`, `entities` for the entitystring
- `shaders` for the shaders

## usage

```console
$ bsputil /path/to/my/map.bsp entities | jq -r '.classname + " @ " + .origin // "N/A"'
target_position @ -1984 1280 1152
ammo_powercell @ 1536 64 -688
ammo_blaster @ 1536 0 -688
ammo_thermal @ 1536 -64 -688
weapon_trip_mine @ 2208 -1088 16
# [...]
```
