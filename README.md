## Overview
face-api detects and outlines faces in an image using the [Azure Cognitive Service for Vision](https://azure.microsoft.com/en-au/products/cognitive-services/vision-services/) Face API.

## Prerequisites
- Azure subscription: Get started with Azure [here](https://azure.microsoft.com/en-us/)
- Face API resource: [Create a Face API resource](https://portal.azure.com/#create/Microsoft.CognitiveServicesFace) in the Azure portal after you have set up your Azure subscription
- Go version 1.20 or higher: [Install Go](https://go.dev/doc/install) for your operating system

## Usage
face-api accepts one image as an argument. Supported image types are JPEG, PNG, GIF (first frame), and BMP. A maximum of 100 faces can be detected in a single image, and faces must be at least 36x36 pixels for images up to 1920x1080 pixels. For larger images, faces will need to be larger in size.

### Running face-api
Navigate to the repository root for the following steps:

Create a `.env` file and add values for the following variables (found in the Azure portal):
```
API_KEY=
ENDPOINT=
```
> The endpoint will be in the format of https://{face-api-resource-name}.cognitiveservices.azure.com/face/v1.0/detect

To build and run the program, run:
```
> go build
> ./face-api <path/to/myImage.png>
```

The resulting image will be saved in the same directory and under the same name as the original image with an `_output` suffix (ex: `"myImage.png"` will generate a `"myImage_output.png"` file).