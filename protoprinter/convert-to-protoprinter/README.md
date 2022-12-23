## `convert-to-protoprinter`

Convert stls to be printed on the protoprinter.

Uses a hack to insert extra instructions, processed after gcode generation.

The instructions deal with starting with a high initial bed temperature so that the base can stick, but lowering it afterwards,
because the material burns otherwise.

### Using

```bash
nix develop
convert-to-protoprinter <stl files>
```
