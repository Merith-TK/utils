To use the `mc2se` command, follow the syntax below:

```shell
mc2se -i schematic.litematica -o outputfile -d definitions.json
```

Here's a breakdown of the command:

- `-i schematic.litematica`: Specifies the input file, in this case, `schematic.litematica`.
- `-o output.sbc`: (optional) Specifies the output file name. If not provided, the default name will be `output.sbc`.
- `-d definitions.json` (optional): Specifies the output file for the block definitions. You can replace `definitions.json` with your desired file name.

The definitions file will always generate at the location if not found

Make sure to replace `schematic.litematica` with the actual name of your input file and `output.sbc` with your desired output file name.
