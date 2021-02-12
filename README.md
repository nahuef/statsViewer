<h1 align="center">Stats Viewer</h1>

<p align="center">
    <img width="100%" alt="Screenshot" src="docs/statsViewer.png">
</p>

## The tool & why
*Stats Viewer* is a tool to help you track score progress and improvement over time by displaying a list **sorted by your most played scenarios** (so you can focus on the bigger picture) with *interactive charts* that shows your max and average scores for each, grouped by day.

Originally, I just wanted to know which were the scenarios I had played the most so I could look for those charts directly, instead of clicking one by one on tons of scenarios I had only played once or twice. And then, the idea to have the charts on the page itself came to me and... voilá!

- Using the toolbox in the top right of each chart you can **zoom in** and **download them as an image** to share on social media.

- Hover over data points (the dots on the lines) to check the **FOV and sensibility you used to get high scores** and figure out what worked best for you.

- Toggle `max scores` and `average scores` lines by clicking on the legend at the top of the chart.

- Quickly go to a specific scenario by using the browser's search function (`ctrl + f` by default in most browsers).

## Online sample
Not hyped yet?!

Check out this **[online sample](https://nahuef.github.io/statsViewer/)** before you download and run the tool to see if it would add value to your training and progress tracking.

## How does it work?
The application will process the files in the `stats` folder to generate a list with a chart per scenario, and output `StatsViewer.html` file with the data in the same directory as the tool's executable `statsViewer.exe`.

That also means that if you lost some of your stats files or your progress by switching PC's, uninstalling Kovaak's, reinstalling your operating system or by any other means, that data won't be available in the result.

## How to use
*It's super simple, I just like to write a lot.*

Basicaly you download the app, set the path to your `stats` folder (which is even optional, see third option on step 2) and run it.
1. Download and extract the latest release of the tool [here](https://github.com/nahuef/statsViewer/releases).
2. Set the path to your Kovaak's `stats` folder. There are three **different options** for this step.
    1. First option. If you are already using the [Progress Sheet Updater](https://github.com/VoltaicHQ/Progress-Sheet-Updater), copy the `config.json` file from that tool and paste it inside `StatsViewer` folder (the one you extracted in step 1). Done.

    2. Second option. Inside `StatsViewer` folder (the one you extracted in step 1) create a file named `config.json` and paste the following snippet:
        ```
        {
            "stats_path": "C:\\Program Files (x86)\\Steam\\steamapps\\common\\FPSAimTrainer\\FPSAimTrainer\\stats"
        }
        ```
        If your steam library is not installed in the default location you will need change the path manually to point it to the right `stats` folder (see below "How to find `stats` folder"). Make sure you either use double `\\` or single `/`. Done.

    3. Third option **for the lazy ones**. Copy the contents from `StatsViewer` folder (the one you extracted in step 1) and paste them in the same directory of your `stats` folder.

        The executable `statsViewer.exe` has to be in the same directory as the `stats` folder (not inside). Done.

        For example, if your `stats` path is
        ```
        C:\Steam\steamapps\common\FPSAimTrainer\FPSAimTrainer\stats
        ```
        You want to paste them in
        ```
        C:\Steam\steamapps\common\FPSAimTrainer\FPSAimTrainer\
        ```
3. Run the tool by executing or double clicking `statsViewer.exe`.

4. A terminal/console will open and if everything goes right you'll see a progress bar. Wait for a few seconds as it processes yor data.

5. Once finished, that window will close itself and you should have a new file `StatsViewer.html` in the same directory as the executable `statsViewer.exe`. Double click the `.html` file to open it in your default browser. **That's it, done!**

> *TIP: Run the tool everytime you want to update the charts*

### How to find your `stats` folder
(To complete the second step and second option of "How to use")
1. Open your Steam application and go to Library -> Collections, find `Kovaak 2.0` in the list, right-click -> Manage -> Browse local files.
    <p align="center">
        <img alt="Screenshot" src="docs/browseLocalFiles.png">
    </p>

2. Now we are in Kovaak's installation directory, we need to open the folder `FPSAimTrainer`.
    <p align="center">
        <img alt="Screenshot" src="docs/installationFolder.png">
    </p>

3. Here we can see the `stats` folder, we are almost there! Open it.
    <p align="center">
        <img alt="Screenshot" src="docs/FPSAimTrainerDir.png">
    </p>

4. From inside the `stats` folder, click on the blank space to the right of the selected text shown in the screenshot below, once you do that you should also have the path to the `stats` folder selected and ready to be copy pasted into the `config.json` file.
    <p align="center">
        <img alt="Screenshot" src="docs/statsPath.png">
    </p>

Once you paste that path in the `config.json` file, make sure you duplicate each `\`.
For the screenshots shown above the `config.json` file would look like this:
```
{
	"stats_path": "E:\\GamesSSD\\SteamLibrary\\steamapps\\common\\FPSAimTrainer\\FPSAimTrainer\\stats"
}
```


## Troubleshooting & support
Tested on Windows 10.

If you need help or encounter any bug, feel free to [open an issue](https://github.com/nahuef/statsViewer/issues/new) or contact me via Discord at Malhumoradour#5542 and send me a screenshot of the error, if any.

Suggestions and PR's welcome!

## Disclaimer
This is my first free open source side project and also my first time programming in Go.

Shipping was prioritized over code quality so if you are a dev and want to check the code, I offer my sincere apology for what your eyes are about to see.

## Build it from source
Go 1.15 required.

```bash
$ git clone https://github.com/nahuef/statsViewer
$ cd statsViewer
$ go build
```
