# lols
slack integration that automatically uploads images to digital ocean spaces and can recall them based on file name descriptors.

## Building from source

`make plugin` to create a plugin (designed to work with github.com/thorfour/sillyputty)

## Running

upload the plugin found at `/bin/plugin/` to the plugins directory of a [sillyputty](https://github.com/thorfour/sillyputty) server

### Using from slack integration

lols supports the following commands:
  - new <image_url_to_upload> <new_image_name>
  - \<sub strings found in image name>
