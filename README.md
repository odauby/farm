# farm
FARM Ain't a Rom Manager

The purpose of this project is to create a small command-line based tool to manage roms in the simplest way possible.

Basically it needs 3 parameters:
* a **DAT file** -- common XML file format used in the emulation and in the digital asset collectors worlds.
* a **source directory** -- FARM will recursively parse the content of that directory in search for (rom) files described in the DAT file.
* a **destination directory** -- all the sets that seems to be complete (required assets in the source directory are matching the content of the DAT file) are copied to the destination directory.
* a zip creationg **mode** -- Assets are zipped in the destination directory, but these zip files will be created according to one of these strategies: non-merge, merge or split. More information about this on this [MameDEV page](https://docs.mamedev.org/usingmame/aboutromsets.html)

This is very early stage of development, but I'm already getting the first results:
```
$ farm -datfile mame209.dat -source ./chaosdir -dest ./roms -mode nonmerge
2019/05/09 22:16:59 f.a.r.m. is starting
2019/05/09 22:16:59 start reading file mame209.dat
2019/05/09 22:16:59 done reading file 
2019/05/09 22:16:59 start marshalling XML
2019/05/09 22:17:02 done marshalling XML
2019/05/09 22:17:02 [name]        MAME - ROMs
2019/05/09 22:17:02 [description] MAME - ROMs (v0.209)
2019/05/09 22:17:02 [category]    Arcade
2019/05/09 22:17:02 [version]     0.209
2019/05/09 22:17:02 [author]      ARMAX
2019/05/09 22:17:02 [comment]     Created by Data File Manager v3.0
2019/05/09 22:17:02 start parsing sets and roms
2019/05/09 22:17:02 done parsing sets and roms
2019/05/09 22:17:02 found 36330 sets and 135199 roms
2019/05/09 22:17:02 parsing source directory ./chaosdir
2019/05/09 22:17:02 identified 22 unique roms in directory ./chaosdir
2019/05/09 22:17:02 Set pacman is complete
2019/05/09 22:17:02 creating zip: ./roms/pacman.zip
2019/05/09 22:17:02 Set puckman is complete
2019/05/09 22:17:02 creating zip: ./roms/puckman.zip
```
