# KSP GameData Manager
The KSP GameData Manager is a simple graphical program that facilitates easy swapping of GameData folders in KSP.

## Building
Dependencies:
- Go 1.14.5 or later
- C Compiler
- Go Package `fyne.io/fyne`
- Go Package `github.com/plus3it/gorecurcopy`
- Go Package `github.com/gen2brain/dlgs`

`go build gamedatamanager.go`

## Installing
Put the `gamedatamanager` excecutable in the KSP root directory. Create a folder called `GameArchive`.

I highly recommend that you either start with the stock game, or else make a folder that contains the Squad (and, if you wish, SquadExpansion) folders.
That way, you can clone the stock GameData easily for fresh installs, etc.

## Using
Run the `gamedatamanager` executable. Usage should be fairly intuitive.
- Clone GameData: allows you to make a copy of a GameData folder you have archived.
- Edit Archived GameData Info: allows you to edit the data associated with a GameData folder stored in the archive.
- Use Archived GameData: Replaces the KSP GameData folder with one from the Archive. PLEASE NOTE: THIS WILL ERASE THE CONTENTS OF YOUR GAMEDATA. BE SURE TO ARCHIVE YOUR GAMEDATA IF YOU WANT TO KEEP IT.
- Delete Archived Gamedata: Allows you to delete an archived GameData folder.

Please remember that the directory structure inside your KSP Directory is essential: There must be both a GameData and GameArchive folder for this program to work correctly. If you tamper with the JSON files in the GameArchive, it might cause issues.
