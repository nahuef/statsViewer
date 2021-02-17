# Find your `stats` folder
(To complete the second step and second option of "How to use")
1. Open your Steam application and go to Library -> Collections, find `Kovaak 2.0` in the list, right-click -> Manage -> Browse local files.
    <p align="center">
        <img alt="Screenshot" src="browseLocalFiles.png">
    </p>

2. Now we are in Kovaak's installation directory, we need to open the folder `FPSAimTrainer`.
    <p align="center">
        <img alt="Screenshot" src="installationFolder.png">
    </p>

3. Here we can see the `stats` folder, we are almost there! Open it.
    <p align="center">
        <img alt="Screenshot" src="FPSAimTrainerDir.png">
    </p>

4. From inside the `stats` folder, click on the blank space to the right of the selected text shown in the screenshot below, once you do that you should also have the path to the `stats` folder selected and ready to be copied and pasted into the `config.json` file.
    <p align="center">
        <img alt="Screenshot" src="statsPath.png">
    </p>

Paste the path in the `config.json` file, once pasted make sure you duplicate each `\`.

For the screenshots shown above the `config.json` file should look like this:
```
{
	"stats_path": "E:\\GamesSSD\\SteamLibrary\\steamapps\\common\\FPSAimTrainer\\FPSAimTrainer\\stats"
}
```
