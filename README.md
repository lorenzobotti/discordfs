# discordfs

## Installation
To compile the CLI executable:
```
go build -o discfs ./cmd
```

The executable looks for a config file in your home folder. On UNIX-like systems (Linux, the BSDs, MacOS) it's `$home` (`/home/username`), on Windows it's `%USERPROFILE%` (`C:\\Users\\username`)

Create a file called `disc.yaml` in your home folder and write:
```
token: "<your-bot-token>"
channel: "<channel-id>"
```

Alternatively you can set a different location for the config file with `discfs --config="<config-file>"`


## Development
Create a file named `secrets.go` in the root folder of the project and write:
```
package discordfs

const authToken = `<your-bot-token>`
const channelId = `<channel-id>`
```

`TestNewDownload` looks for test cases in the `test_files` folder of the project. You can:
 * Populate this folder with test files that also exist on the cloud channel. They must have the same name as the ones on the cloud
 * Comment out the test. It's not necessary, `TestUploadAndDelete` fulfills the same purpose