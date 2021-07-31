# ElectronVisualizer

https://angelmariages.github.io/ElectronVisualizer/

![image](https://user-images.githubusercontent.com/425506/127748775-43e3c7c9-40d3-49d8-864f-c1b8037d8b54.png)

In order to understand better the habits of electricity consumption of a house it's easier to have a broad perspective to find patterns.

This is why, with the inspiration of the [video](https://www.youtube.com/watch?v=jYPFiMaeOv8), I have decided to start creating such helpful graphs.

## How to use with your data

### 1. Get the CSV from your distributor

Example from Endesa [Link](https://zonaprivada.edistribucion.com/areaprivada/s/wp-massivemeasuredownload-v3):

![image](https://user-images.githubusercontent.com/425506/127748852-6c387d65-9c93-4907-b98d-1871215c20e2.png)


### 2. Generate the JSON to be used in the app

Put your csv in the `csv-parser` folder with the name `raw_data.csv` and execute `main.go`. A new data.json will be generated with all your values.
